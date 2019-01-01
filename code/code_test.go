package code

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       OpCode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetGlobal), 255}},
	}

	for _, tt := range tests {
		Instruction := Make(tt.op, tt.operands...)

		assert.Equal(t, len(Instruction), len(tt.expected), "instruction has wrong length")

		for i, _ := range tt.expected {
			assert.Equal(t, Instruction[i], tt.expected[i], fmt.Sprintf("wrong byte at pos %d", i))
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpGetLocal 1
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	assert.Equal(t, expected, concatted.String(), "instructions wrongly formatted")
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        OpCode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpGetLocal, []int{255}, 1},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		assert.NoError(t, err, "definition not found")

		operandsRead, n := ReadOperands(def, instruction[1:])
		assert.Equal(t, n, tt.bytesRead, "n wrong")

		for i, want := range tt.operands {
			assert.Equal(t, want, operandsRead[i], "operands wrong")
		}
	}
}
