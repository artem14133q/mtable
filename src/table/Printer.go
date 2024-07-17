package table

import (
	"fmt"
	"strings"
)

type Printer struct {
	table *Table
	result string
	pure bool
}

func NewTablePrinter(table *Table, pure bool) Printer {
	return Printer{table: table, pure: pure}
}

func (printer *Printer) PrintResult() {
	fmt.Print(printer.result)
}

func (printer *Printer) CreateTable(short bool)  {
	printer.result += printer.CreateTopLine()
	printer.result += printer.CreateHeader()
	printer.result += printer.CreateDelimiterLine()

	size := len(printer.table.raws)

	for i := 0; i < size; i++ {
		printer.result += printer.CreateDataLine(i)

		if !short && i + 1 != size {
			printer.result += printer.CreateDelimiterLine()
		}
	}

	printer.result += printer.CreateBottomLine()
}

func (printer *Printer) CreateDelimiterLine() string {
	boxSettings := printer.table.boxSettings

	return printer.CreateHorizontalLine(
		boxSettings.LeftOuterIntersection,
		boxSettings.InnerHorizontalLine,
		boxSettings.InnerIntersection,
		boxSettings.RightOuterIntersection)
}

func (printer *Printer) CreateTopLine() string {
	boxSettings := printer.table.boxSettings

	return printer.CreateHorizontalLine(
		boxSettings.LeftTopCorner,
		boxSettings.OuterHorizontalLine,
		boxSettings.TopOuterIntersection,
		boxSettings.RightTopCorner)
}

func (printer *Printer) CreateBottomLine() string {
	boxSettings := printer.table.boxSettings

	return printer.CreateHorizontalLine(
		boxSettings.LeftBottomCorner,
		boxSettings.OuterHorizontalLine,
		boxSettings.BottomOuterIntersection,
		boxSettings.RightBottomCorner)
}

func (printer *Printer) CreateHorizontalLine(
	leftIntersectionChar string, lineChar string, delimiterChar string, rightIntersectionChar string) string {
	result := leftIntersectionChar
	sizes := printer.table.GetColumnSizes()

	var list []string

	for i := 0; i < len(sizes); i++ {
		list = append(list, mulLine(sizes[i] + 2, lineChar))
	}

	return result + strings.Join(list, delimiterChar) + rightIntersectionChar + "\n"
}

func (printer *Printer) CreateHeader() string {
	boxSettings := printer.table.boxSettings
	result := boxSettings.OuterVerticalLine

	var lines []string

	for _, column := range printer.table.columns {
		lines = append(lines, createCellContent(
			column.name.Format(printer.pure), column.width - column.size, column.alignment))
	}

	return result + strings.Join(lines, boxSettings.InnerVerticalLine) + boxSettings.OuterVerticalLine + "\n"
}

func (printer *Printer) CreateDataLine(rawIndex int) string {
	boxSettings := printer.table.boxSettings
	result := printer.table.boxSettings.OuterVerticalLine

	var lines []string

	raw := printer.table.raws[rawIndex]

	for i, column := range printer.table.columns {
		value := raw.items[i].value

		lines = append(
			lines, createCellContent(value.Format(printer.pure), column.width - value.GetSize(), column.alignment))
	}

	return result + strings.Join(lines, boxSettings.InnerVerticalLine) + boxSettings.OuterVerticalLine + "\n"
}

func createSpace(len int) string {
	return mulLine(len, " ")
}

func mulLine(len int, chars string) string {
	result := ""

	for i := 0; i < len; i++ {
		result += chars
	}

	return result
}

func createCellContent(content string, width int, alignment int) string {
	space := 0

	if alignment == CenterAlignment {
		space = width / 2
	} else if alignment == RightAlignment {
		space = width
	}

	return strings.Join([]string{" ", createSpace(space), content, createSpace(width - space), " "}, "")
}
