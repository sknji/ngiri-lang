package vm

import (
	"github.com/marmotini/monkey-lang/code"
	"github.com/marmotini/monkey-lang/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
	basePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
