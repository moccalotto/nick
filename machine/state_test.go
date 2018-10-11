package machine

import "testing"

var script string = `

commandWithoutArgs
evolve 				ABC
kim    FOO 		BAR		# comment that will be stripped but stored
		mov    eax ebx

`

func TestScriptToInstructions(t *testing.T) {
	m := MachineFromScript(script)
	s := m.tape

	x := []Instruction{
		Instruction{
			"commandWithoutArgs",
			[]string{},
			"",
		},
		Instruction{
			"evolve",
			[]string{"ABC"},
			"",
		},
		Instruction{
			"kim",
			[]string{"FOO", "BAR"},
			"comment that will be stripped but stored",
		},
		Instruction{
			"mov",
			[]string{"eax", "ebx"},
			"",
		},
	}

	if len(s) != len(x) {
		t.Errorf("Expected %d instructions, but got %d instead (%+v)", len(x), len(s), s)
		return
	}

	for i, xInstr := range x {
		sInstr := s[i]

		if sInstr.cmd != xInstr.cmd {
			t.Errorf("Expected instruction %d to be »%s«, but it was »%s«", i, xInstr.cmd, sInstr.cmd)
		}

		if len(xInstr.args) != len(sInstr.args) {
			t.Errorf(
				"Expected instruction %d (%s) to have %d args, but it had %d",
				i, 
				xInstr.cmd,
				len(xInstr.args),
				len(sInstr.args),
			)
			continue
		}

		for j, xArg := range sInstr.args {

			sArg := sInstr.args[j]

			if sArg != xArg {
				t.Errorf(
					"Expected arg #%d for instruction %d to be »%s«, but it was »%s«",
					j,
					i,
					xArg,
					sArg,
				)
			}

		}
	}
}
