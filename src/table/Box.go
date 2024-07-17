package table

type BoxSettings struct {
	LeftTopCorner string
	RightTopCorner string

	LeftBottomCorner string
	RightBottomCorner string

	OuterHorizontalLine string
	InnerHorizontalLine string

	OuterVerticalLine string
	InnerVerticalLine string

	LeftOuterIntersection string
	RightOuterIntersection string
	TopOuterIntersection string
	BottomOuterIntersection string

	LeftInnerIntersection string
	RightInnerIntersection string
	TopInnerIntersection string
	BottomInnerIntersection string
	InnerIntersection string
}

func NewBoxSettings() *BoxSettings {
	return &BoxSettings{
		"┏",
		"┓",
		"┗",
		"┛",
		"━",
		"─",
		"┃",
		"│",
		"┠",
		"┨",
		"┯",
		"┷",
		"├",
		"┤",
		"┬",
		"┴",
		"┼",
	}
}