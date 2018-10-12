package machine

import (
	"github.com/moccaloto/nick/field"
	"regexp"
)

func HandleInitInstruction(m *Machine, args []string) {
	m.Assert(m.Field == nil, "You cannot call 'init' more than once!")
	m.Assert(len(args) == 1, "'init' instruction must have exactly 1 argument, but it was given %d", len(args))

	arg0 := m.MustGetString(args[0])
	nums := regexp.MustCompile(`^(\d+)x(\d+)$`).FindStringSubmatch(arg0)

	m.Assert(
		len(nums) == 3,
		"Invalid argument for 'init'. Must be [number]x[number], but it was '%s' (%+v)",
		arg0,
		nums,
	)

	m.Field = field.NewField(
		m.StrToInt(nums[1]),
		m.StrToInt(nums[2]),
	)
}
