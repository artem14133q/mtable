package table

import "main/src/formater"

type RawItem struct {
	parent *Raw
	value *formater.TextFormatter
	position int
}

func (item *RawItem) Value() *formater.TextFormatter {
	return item.value
}

func NewRawItem(parent *Raw, value string) *RawItem {
	formatter := formater.New(value)
	formatter.SetTextColor(parent.parent.columns[parent.size].defaultRawColor)

	size := formatter.GetSize()

	width := &parent.parent.columns[parent.size].width

	if size > *width {
		*width = size
	}

	return &RawItem{parent: parent, value: formatter, position: parent.size}
}
