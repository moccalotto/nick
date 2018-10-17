package machine

import (
	"strconv"
)

func init() {
	InstructionHandlers["seed"] = Seed
}

func Seed(m *Machine) {
	seed, err := strconv.ParseInt(m.ArgAsString(0), 0, 64)

	m.Assert(err == nil, "First arg to seed must be a 64-bit integer")

	m.Rng.Seed(seed)
	m.Seed = seed
}
