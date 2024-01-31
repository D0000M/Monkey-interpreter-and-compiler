package vm

import (
	"monkey/code"
	"monkey/object"
)

// 栈帧
type Frame struct {
	cl          *object.Closure // 让Frame保持对Closure的支持
	ip          int             // 指向该帧的命令指针
	basePointer int             // 该帧执行完后恢复的命令指针值
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}

// 每当pushFrame，就运行新进的指令
func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

// 进入函数体时运行，将vm运行的指令推到新来的frame中
func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}
