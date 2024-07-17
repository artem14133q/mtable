package table

import (
	"main/src/color"
	"main/src/formater"
)

const LeftAlignment = 1
const CenterAlignment = 2
const RightAlignment = 3

type Column struct {
	parent *Table
	name *formater.TextFormatter
	defaultRawColor *color.Color
	width int
	size int
	alignment int
}

func NewColumn(parent *Table, name string) *Column {
	formatter := formater.New(name)

	size := formatter.GetSize()

	return &Column{parent: parent, name: formatter, size: size, width: size, alignment: LeftAlignment}
}

func (column *Column) GetName() *formater.TextFormatter {
	return column.name
}

func (column *Column) SetDefaultRawColor(color *color.Color)  {
	column.defaultRawColor = color
}

func (column *Column) SetAlignment(alignment int) {
	column.alignment = alignment
}