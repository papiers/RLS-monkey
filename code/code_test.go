package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length: got %d, want %d", len(instruction), len(tt.expected))
		}
		for i, b := range tt.expected {
			if instruction[i] != b {
				t.Errorf("instruction[%d]: got %d, want %d", i, instruction[i], b)
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}
	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`
	concat := Instructions{}
	for _, ins := range instructions {
		concat = append(concat, ins...)
	}
	if concat.String() != expected {
		t.Errorf("instructions string: got %s, want %s", concat.String(), expected)
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %v", err)
		}
		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("number of bytes read: got %d, want %d", n, tt.bytesRead)
		}
		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Fatalf("operand %d: got %d, want %d", i, operandsRead[i], want)
			}
		}
	}
}
