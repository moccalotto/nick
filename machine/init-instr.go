package machine

import (
	"github.com/moccaloto/nick/field"
	"regexp"
)

func init() {
	InstructionHandlers["init"] = Init
}

func Init(m *Machine) {
	m.Assert(m.Field == nil, "You cannot call 'init' more than once!")

	arg0 := m.ArgAsString(0)
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
