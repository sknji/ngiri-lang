package vm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/marmotini/monkey-lang/code"
	"github.com/marmotini/monkey-lang/compiler"
	"github.com/marmotini/monkey-lang/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

var ErrStackOverflow = errors.New("stack overflow")
var ErrStackUnderflow = errors.New("stack underflow")

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int

	globals []object.Object

	frames     []*Frame
	frameIndex int
}

func NewVM(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalsSize),
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		frames:       frames,
		frameIndex:   1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := NewVM(bytecode)
	vm.globals = s

	return vm
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

// Run contains the fetch-decode-execute cycle
func (vm *VM) Run() error {
	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ins := vm.currentFrame().Instructions()
		op := code.OpCode(ins[vm.currentFrame().ip])

		switch op {
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[vm.currentFrame().ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case code.OpJump:
			pos := int(code.ReadUint16(ins[vm.currentFrame().ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[vm.currentFrame().ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpDiv, code.OpMul, code.OpSub:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[vm.currentFrame().ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[vm.currentFrame().ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[vm.currentFrame().ip+1:])
			vm.currentFrame().ip += 1

			vm.stack[vm.currentFrame().basePointer+int(localIndex)] = vm.pop()
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[vm.currentFrame().ip+1:])
			vm.currentFrame().ip += 1

			err := vm.push(vm.stack[vm.currentFrame().basePointer+int(localIndex)])
			if err != nil {
				return err
			}
		case code.OpCall:
			vm.currentFrame().ip += 1
			fn, ok := vm.stack[vm.sp-1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non-function")
			}

			frame := NewFrame(fn, vm.sp)
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + fn.NumLocals
		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeComparison(op code.OpCode) error {
	right := vm.pop()
	left := vm.pop()

	if right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func nativeToBooleanObject(result bool) object.Object {
	if result {
		return True
	}

	return False
}

func (vm *VM) executeIntegerComparison(op code.OpCode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(nativeToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeBinaryOperation(op code.OpCode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryStringOperation(op code.OpCode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operation: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

func (vm *VM) executeBinaryIntegerOperation(op code.OpCode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return ErrStackOverflow
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	if vm.sp < 1 {
		return nil
	}

	o := vm.stack[vm.sp-1]
	vm.sp--

	return o
}

func (vm *VM) debugStack() string {
	var buf bytes.Buffer

	buf.WriteString("Stack:\n")
	for i := 0; i < vm.sp; i++ {
		if vm.stack[i] != nil {
			buf.WriteString(fmt.Sprintf("%d. %s\n", i, vm.stack[i].Inspect()))
		}
	}

	return buf.String()
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}

func (vm *VM) pushFrame(frame *Frame) {
	vm.frames[vm.frameIndex] = frame
	vm.frameIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.frameIndex--
	return vm.frames[vm.frameIndex]
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}

}
