package machine

import "testing"

var script string = `

commandWithoutArgs
evolve 				ABC
kim    FOO 		BAR		# comment that will be stripped but stored
		mov    eax ebx


# ÆØÅUUU
`

func TestScriptToInstructions(t *testing.T) {
	m := MachineFromScript(script)
	s := m.Tape

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

		if sInstr.Cmd != xInstr.Cmd {
			t.Errorf("Expected instruction %d to be »%s«, but it was »%s«", i, xInstr.Cmd, sInstr.Cmd)
		}

		if len(xInstr.Args) != len(sInstr.Args) {
			t.Errorf(
				"Expected instruction %d (%s) to have %d Args, but it had %d",
				i,
				xInstr.Cmd,
				len(xInstr.Args),
				len(sInstr.Args),
			)
			continue
		}

		for j, xArg := range sInstr.Args {

			sArg := sInstr.Args[j]

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
