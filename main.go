package main

import (
	"bufio"
	"log"
	"main/src/color"
	"main/src/table"
	"os"
	"strconv"
	"strings"
)

const DefaultSeparator = ","

const (
	ArgTypeFlag = iota + 1
	ArgTypeParameter
)

const (
	ArgCompactTable = iota + 1
	ArgDisableAsciiCharsInBox
	ArgDisableAnsi
	ArgSetTableParameters
	ArgDefineCsvSeparator
)

var params = map[string]string{
	"headersColor":  "coral",
	"rawsColor":     "sea green",
	"firstColColor": "sky blue",
	"numerate":      "",
}

type args struct {
	name     string
	alias    string
	argType  int
	callback func(arg string) string
}

var argsTypes = map[int]args{
	// Flags
	ArgCompactTable:           {"--compact", "", ArgTypeFlag, nil},
	ArgDisableAsciiCharsInBox: {"--non-ascii-box", "", ArgTypeFlag, nil},
	ArgDisableAnsi:            {"--pure", "", ArgTypeFlag, nil},

	// Params with values
	ArgSetTableParameters: {"--table-params", "-t", ArgTypeParameter, func(arg string) string {
		if arg == "" {
			return ""
		}

		paramItems := strings.Split(arg, ",")

		for _, item := range paramItems {
			keyValue := strings.SplitN(item, "=", 2)

			if len(keyValue) < 2 {
				log.Fatalf("Cannot find value for key %s", keyValue[0])
			}

			params[keyValue[0]] = keyValue[1]
		}

		return ""
	}},
	ArgDefineCsvSeparator: {"--sep", "-s", ArgTypeParameter, func(arg string) string {
		if arg == "" {
			return DefaultSeparator
		}

		return arg
	}},
}

var argFlags = map[int]bool{}
var argParams = map[int]string{}

func getStdin() string {
	stdin := os.Stdin

	content := ""

	stat, err := stdin.Stat()

	if err != nil {
		log.Fatal("Cannot read from stdin")
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return content
	}

	reader := bufio.NewReader(stdin)

	for {
		text, err := reader.ReadString('\n')

		if text == "\n" {
			break
		}

		if err == nil {
			content += text
			continue
		}

		break
	}

	return content
}

func createTable(headers []string, raws [][]string) {
	var boxSettings *table.BoxSettings

	if argFlags[ArgDisableAsciiCharsInBox] {
		boxSettings = &table.BoxSettings{
			LeftTopCorner:           "+",
			RightTopCorner:          "+",
			LeftBottomCorner:        "+",
			RightBottomCorner:       "+",
			OuterHorizontalLine:     "-",
			InnerHorizontalLine:     "-",
			OuterVerticalLine:       "|",
			InnerVerticalLine:       "|",
			LeftOuterIntersection:   "+",
			RightOuterIntersection:  "+",
			TopOuterIntersection:    "+",
			BottomOuterIntersection: "+",
			LeftInnerIntersection:   "+",
			RightInnerIntersection:  "+",
			TopInnerIntersection:    "+",
			BottomInnerIntersection: "+",
			InnerIntersection:       "+",
		}
	} else {
		boxSettings = table.NewBoxSettings()
	}

	_table := table.New()
	_table.SetBoxSettings(boxSettings)

	headerColor := color.New(params["headersColor"])
	rawsColor := color.New(params["rawsColor"])
	firstColColor := color.New(params["firstColColor"])

	if params["numerate"] != "" {
		headers = append([]string{params["numerate"]}, headers...)
	}

	for i, header := range headers {
		column := _table.CreateColumn(header)
		column.GetName().SetTextColor(headerColor)

		if i == 0 {
			column.SetDefaultRawColor(firstColColor)
			continue
		}

		column.SetDefaultRawColor(rawsColor)
	}

	for i, line := range raws {
		if params["numerate"] != "" {
			line = append([]string{strconv.Itoa(i + 1)}, line...)
		}

		_table.CreateRawWithValues(line)
	}

	tablePrinter := table.NewTablePrinter(_table, argFlags[ArgDisableAnsi])
	tablePrinter.CreateTable(argFlags[ArgCompactTable])
	tablePrinter.PrintResult()
}

func workCsv(content string) {
	lines := strings.Split(content, "\n")
	headers := strings.Split(strings.Trim(lines[0], "\r"), argParams[ArgDefineCsvSeparator])

	headersLen := len(headers)

	var raws [][]string

	if lines[len(lines)-1] == "" {
		lines = lines[1 : len(lines)-1]
	}

	for _, raw := range lines {
		raws = append(raws, strings.SplitN(strings.Trim(raw, "\r"), argParams[ArgDefineCsvSeparator], headersLen))
	}

	createTable(headers, raws)
}

func getParam(argv []string, param string, short string) (bool, string) {
	if param[:2] != "--" || param == "" {
		log.Fatalf("Cannot resolve parameter: '%s'\n", param)
	}

	ok, index := indexOf(argv, param)
	if !ok {
		if short == "" {
			return false, ""
		}

		if short[:1] != "-" {
			log.Fatalf("Cannot resolve short parameter: '%s'\n", short)
		}

		ok, index = indexOf(argv, short)
		if !ok {
			return false, ""
		}
	}

	return true, argv[index+1]
}

func indexOf(argv []string, param string) (bool, int) {
	for i, v := range argv {
		if v == param {
			return true, i
		}
	}

	return false, 0
}

func hasFlag(argv []string, param string) bool {
	ok, _ := indexOf(argv, param)
	if ok {
		return true
	}

	return false
}

func main() {
	argv := os.Args[1:]

	for argId, argOptions := range argsTypes {
		if argOptions.argType == ArgTypeFlag {
			argFlags[argId] = hasFlag(argv, argOptions.name)
		} else if argOptions.argType == ArgTypeParameter {
			_, value := getParam(argv, argOptions.name, argOptions.alias)

			if argOptions.callback != nil {
				value = argOptions.callback(value)
			}

			argParams[argId] = value
		}
	}

	stdin := getStdin()

	if stdin != "" {
		workCsv(stdin)
	}
}
