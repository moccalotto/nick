package machine

import (
	"log"
	"strings"
)

func init() {
	InstructionHandlers["log"] = Log
}

// TODO: this should log via the machine. The machine should have a mechanism for outputting log commands.
func Log(m *Machine) {
	buf := make([]string, m.ArgCount())
	for n := 0; n < m.ArgCount(); n++ {
		buf[n] = m.ArgAsString(n)
	}
	log.Printf("LOG: %s", strings.Join(buf, " "))
}
