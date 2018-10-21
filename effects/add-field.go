package effects

import (
	"github.com/moccalotto/nick/field"
)

type AddField struct {
	source field.Field

	/*
		Modes/features/Flags

		Overwrite:
			Cells in source simply overwrite cells in target
			this is only useful of an offset and/or size restriction is applied.

		LivingOnly
			cells that are alive in the source field will be copied to the
			target field
		DeadOnly
			cells that are dead in the source field will be copied to the
			target field
		Add
			THe value of the cells in the source branch will be added to the
			value cells of the target branch
			in the target branch
			Since dead cells are normally zero. Dead cells from the source field
			would not affect the cells in the target field.
		Subtract
			Values of the cells in the source branch will be subtracted from the
			value of the cells in the target branch
	*/
}
