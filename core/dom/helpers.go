package dom

import (
	"strconv"
)

// ------------------------
//     Primal Type Converters
// ------------------------

func stringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func stringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func stringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

// ------------------------
//     Attr Converters
// ------------------------

func rawAttrsToAttibuteList(attrs [][2]string) RawAttributeList {
	var list RawAttributeList
	for _, attr := range attrs {
		list = append(list, RawAttribute{attr[0], attr[1]})
	}
	return list
}

func stringToDisplay(s string) Display {
	switch s {
	case "block":
		return Display_Block
	case "flex":
		return Display_Flex
	case "absolute":
		return Display_Absolute
	default:
		return Display_Block
	}
}

func stringToFlexDirection(s string) FlexDirection {
	switch s {
	case "row":
		return FlexDirection_Row
	case "row-reverse":
		return FlexDirection_RowReverse
	case "column":
		return FlexDirection_Column
	case "column-reverse":
		return FlexDirection_ColumnReverse
	default:
		return FlexDirection_Row
	}
}

func stringToTextAlign(s string) TextAlign {
	switch s {
	case "left":
		return TextAlign_Left
	case "center":
		return TextAlign_Center
	case "right":
		return TextAlign_Right
	default:
		return TextAlign_Left
	}
}

func stringToTextDecoration(s string) TextDecoration {
	switch s {
	case "none":
		return TextDecoration_None
	case "underline":
		return TextDecoration_Underline
	case "line-through":
		return TextDecoration_LineThrough
	default:
		return TextDecoration_None
	}
}

func stringToPosition(s string) Position {
	switch s {
	case "relative":
		return Position_Relative
	case "absolute":
		return Position_Absolute
	default:
		return Position_Relative
	}
}

func stringToInputType(s string) InputType {
	switch s {
	case "text":
		return InputType_Text
	case "password":
		return InputType_Password
	case "email":
		return InputType_Email
	case "number":
		return InputType_Number
	case "tel":
		return InputType_Tel
	case "url":
		return InputType_URL
	default:
		return InputType_Text
	}
}

func stringToUint8(s string) uint8 {
	i, _ := strconv.Atoi(s)
	return uint8(i)
}

// ------------------------
//      Math Helpers
// ------------------------

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
