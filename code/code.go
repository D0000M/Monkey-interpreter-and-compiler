package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte //字节集合
type Opcode byte         // 操作码

const (
	OpConstant Opcode = iota //添加常数
	OpAdd                    // +
	OpSub                    // -
	OpMul                    // *
	OpDiv                    // /
	OpPop                    // 每个表达式语句执行后，弹出栈顶元素
	OpTrue                   // 将true压入栈
	OpFalse                  // 将false压入栈
	OpEqual
	OpNotEqual
	OpGreaterThan // 有且只有'>',而'<'通过编译器代码重排序实现
	OpMinus
	OpBang
	OpJumpNotTruthy // 有条件跳转
	OpJump          // 直接跳转
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex       // 索引运算符
	OpCall        // 调用函数
	OpReturnValue // 函数的隐式返回和显式返回值
	OpReturn      // 函数没有返回值的返回，用于回到调用之前的状态
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpClosure
	OpGetFree        // 用于在closure里存储自由变量，存储在其Free变量中
	OpCurrentClosure // 引用自身closure
)

type Definition struct {
	Name          string
	OperandWidths []int // 有几个操作数，以及每个操作数占用的字节数
}

var definitions = map[Opcode]*Definition{
	OpConstant:       {"OpConstant", []int{2}}, //只有一个占两个字节宽的操作数
	OpAdd:            {"OpAdd", []int{}},       // OpAdd没有操作数，只是顶部两个弹栈相加后，运算结果压栈
	OpSub:            {"OpSub", []int{}},
	OpMul:            {"OpMul", []int{}},
	OpDiv:            {"OpDiv", []int{}},
	OpPop:            {"OpPop", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpEqual:          {"OpEqual", []int{}},
	OpNotEqual:       {"OpNotEqual", []int{}},
	OpGreaterThan:    {"OpGreaterThan", []int{}},
	OpMinus:          {"OpMinus", []int{}},
	OpBang:           {"OpBang", []int{}},
	OpJumpNotTruthy:  {"OpJumpNotTruthy", []int{2}},
	OpJump:           {"OpJump", []int{2}},
	OpNull:           {"OpNull", []int{}},
	OpGetGlobal:      {"OpGetGlobal", []int{2}},
	OpSetGlobal:      {"OpSetGlobal", []int{2}},
	OpArray:          {"OpArray", []int{2}}, // 操作数为数组元素的数量
	OpHash:           {"OpHash", []int{2}},
	OpIndex:          {"OpIndex", []int{}}, // 默认栈顶有两个元素，一个被索引的对象和作为索引的对象
	OpCall:           {"OpCall", []int{1}}, // 调用函数，有一个占据一字节的操作数表示调用函数的参数数目
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpGetLocal:       {"OpGetLocal", []int{1}},
	OpSetLocal:       {"OpSetLocal", []int{1}},
	OpGetBuiltin:     {"OpGetBuiltin", []int{1}},
	OpClosure:        {"OpClosure", []int{2, 1}}, // 第一个是常量索引用于常量池指定哪个函数转化成闭包，第二个是有多少个自由变量
	OpGetFree:        {"OpGetFree", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])

	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen) // 若没有传入操作数，后面初始化为0
	instruction[0] = byte(op)

	offset := 1

	// 遍历定义好的OperandWidths，从操作数operands一个个取出匹配元素放入指令中
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width { // 取决于操作数的宽度
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

// 用于解码，返回操作数列表，以及已读的字节数
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }
