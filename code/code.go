package code

import (
	"encoding/binary"
	"fmt"
)

type Definition struct {
	Name           string
	OperandsWidths []int
}

var definitions = map[OpCode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("opCode %d undefinied", op)
	}

	return def, nil
}

type Instructions []byte

type OpCode byte

const (
	OpConstant OpCode = iota
)

func Make(op OpCode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandsWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandsWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}

		offset += width
	}
	return instruction
}

func (ins Instructions) String() string {
	return ""
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandsWidths))
	offset := 0

	for i, width := range def.OperandsWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
