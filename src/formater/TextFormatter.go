package formater

import (
	"log"
	"main/src/color"
	"strconv"
	"strings"
)

const EscapeStart = "\x1B["
const EscapeEnd = EscapeStart + "0m"

const EscapeTypeTextColor = 1
const EscapeTypeBackgroundColor = 2
const EscapeTypeBolt = 3
const EscapeTypeUnderline = 4
const EscapeTypeItalic = 5
const EscapeTypeDim = 6
const EscapeTypeBlinking = 7
const EscapeTypeInverse = 8
const EscapeTypeHidden = 9
const EscapeTypeStrikethrough = 10

var EscapeParams = map[int][]string{
	EscapeTypeTextColor:       {"38", "2"},
	EscapeTypeBackgroundColor: {"48", "2"},
	EscapeTypeBolt:            {"1"},
	EscapeTypeDim:             {"2"},
	EscapeTypeItalic:          {"3"},
	EscapeTypeUnderline:       {"4"},
	EscapeTypeBlinking:        {"5"},
	EscapeTypeInverse:         {"7"},
	EscapeTypeHidden:          {"8"},
	EscapeTypeStrikethrough:   {"9"},
}

type EscapeParam struct {
	escapeParamType int
	color           *color.Color
}

type TextFormatter struct {
	text            string
	textColor       *color.Color
	backgroundColor *color.Color
	bolt            bool
	underline       bool
	italic          bool
	dim             bool
	blinking        bool
	inverse         bool
	hidden          bool
	strikethrough   bool
	size            int
}

func New(text string) *TextFormatter {
	formatter := &TextFormatter{}

	for {
		if !IfStringEscaped(text) {
			break
		}

		s, param := GetEscapedStringData(text)

		_type := param.escapeParamType

		if _type == EscapeTypeTextColor {
			formatter.textColor = param.color
		} else if _type == EscapeTypeBackgroundColor {
			formatter.backgroundColor = param.color
		}

		formatter.bolt = _type == EscapeTypeBolt
		formatter.underline = _type == EscapeTypeUnderline
		formatter.italic = _type == EscapeTypeItalic
		formatter.dim = _type == EscapeTypeDim
		formatter.blinking = _type == EscapeTypeBlinking
		formatter.inverse = _type == EscapeTypeInverse
		formatter.hidden = _type == EscapeTypeHidden
		formatter.strikethrough = _type == EscapeTypeStrikethrough

		text = s
	}

	formatter.text = text
	formatter.size = len([]rune(text))

	return formatter
}

func (formatter *TextFormatter) GetSize() int {
	return formatter.size
}

func (formatter *TextFormatter) SetTextColor(color *color.Color) {
	formatter.textColor = color
}

func (formatter *TextFormatter) SetBackgroundColor(color *color.Color) {
	formatter.backgroundColor = color
}

func (formatter *TextFormatter) SetBolt(bolt bool) {
	formatter.bolt = bolt
}

func (formatter *TextFormatter) SetUnderline(underline bool) {
	formatter.underline = underline
}

func (formatter *TextFormatter) Format(pure bool) string {
	if pure {
		return formatter.text
	}

	result := formatter.text

	if formatter.italic {
		result = ItalicString(result)
	}

	if formatter.textColor != nil {
		result = ColorString(result, formatter.textColor)
	}

	if formatter.bolt {
		result = BoltString(result)
	}

	if formatter.underline {
		result = UnderlineString(result)
	}

	if formatter.backgroundColor != nil {
		result = ColorBackgroundString(result, formatter.backgroundColor)
	}

	if formatter.dim {
		result = DimString(result)
	}

	if formatter.blinking {
		result = BlinkingString(result)
	}

	if formatter.inverse {
		result = InverseString(result)
	}

	if formatter.hidden {
		result = HiddenString(result)
	}

	if formatter.strikethrough {
		result = StrikethroughString(result)
	}

	return result
}

func EscapeString(line string, params []string) string {
	return strings.Join([]string{EscapeStart, strings.Join(params, ";"), "m", line, EscapeEnd}, "")
}

func ColorString(line string, color *color.Color) string {
	return ColorStringTyped(line, color, EscapeTypeTextColor)
}

func ColorBackgroundString(line string, color *color.Color) string {
	return ColorStringTyped(line, color, EscapeTypeBackgroundColor)
}

func ColorStringTyped(line string, color *color.Color, _type int) string {
	r, g, b := color.ToStringList()

	return EscapeString(line, append(EscapeParams[_type], r, g, b))
}

func BoltString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeBolt])
}

func UnderlineString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeUnderline])
}

func ItalicString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeItalic])
}

func DimString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeDim])
}

func BlinkingString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeBlinking])
}

func InverseString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeInverse])
}

func HiddenString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeHidden])
}

func StrikethroughString(line string) string {
	return EscapeString(line, EscapeParams[EscapeTypeStrikethrough])
}

func IfStringEscaped(line string) bool {
	return strings.Contains(line, EscapeStart)
}

func EscapeStringParameter(escapeType int) string {
	return strings.Join(EscapeParams[escapeType], ";")
}

func DeterminateEscapeType(line string) int {
	for key := range EscapeParams {
		if strings.Contains(line, EscapeStringParameter(key)) {
			return key
		}
	}

	return 0
}

func GetColorParameter(line string) *color.Color {
	params := strings.Split(line, ";")

	r, rErr := strconv.ParseInt(params[2], 10, 64)
	g, gErr := strconv.ParseInt(params[3], 10, 64)
	b, bErr := strconv.ParseInt(params[4], 10, 64)

	if !(rErr == nil || gErr == nil || bErr == nil) {
		log.Fatal("Error get color parameters: '" + line + "'")
	}

	return &color.Color{R: r, G: g, B: b}
}

func GetEscapedStringData(line string) (string, *EscapeParam) {
	lines := strings.Split(line, EscapeStart)

	str := strings.Join(lines[1:len(lines)-1], EscapeStart)

	lines = strings.Split(str, "m")

	stringParams := lines[0]

	escapeType := DeterminateEscapeType(stringParams)

	escapeParam := &EscapeParam{escapeParamType: escapeType}

	if EscapeTypeTextColor == escapeType || EscapeTypeBackgroundColor == escapeType {
		escapeParam.color = GetColorParameter(stringParams)
	}

	str = strings.Join(lines[1:], "m")

	return str, escapeParam
}
