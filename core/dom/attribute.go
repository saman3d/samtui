package dom

import (
	"errors"
)

var (
	ErrAttrMustBeInt = errors.New("attr must be int")
)

// --------------------
//       RawAttr
// --------------------

type RawAttribute [2]string

type RawAttributeList []RawAttribute

// --------------------
// 	      Attrs
// --------------------

type Attributes struct {
	Display         Display
	Position        Position
	Flex            int
	Focusable       bool
	FlexDirection   FlexDirection
	Color           int
	BackGroundColor int
	Width           int
	MaxWidth        int
	MinWidth        int
	Height          int
	MaxHeight       int
	MinHeight       int
	Border          bool
	VCenter         bool
	HCenter         bool
	Top             int
	Left            int
	ID              string
	TextAlign       TextAlign
	Writable        bool
	TextType        InputType
	ZIndex          uint8
	Overflow        Overflow
}

func NewAttributes(opts ...AttributesOpt) *Attributes {
	a := Attributes{}
	for _, opt := range opts {
		a = opt(a)
	}
	return &a
}

// --------------------
//   Attributes Opts
// --------------------

type AttributesOpt func(Attributes) Attributes

func WithDefaultAttributes() AttributesOpt {
	return func(a Attributes) Attributes {
		return Attributes{
			Display:         Display_Block,
			Position:        Position_Relative,
			FlexDirection:   FlexDirection_Row,
			Focusable:       false,
			Color:           0,
			BackGroundColor: 0,
			Width:           0,
			MaxWidth:        0,
			MinWidth:        0,
			Height:          0,
			MaxHeight:       0,
			MinHeight:       0,
			Flex:            1,
			Border:          false,
			VCenter:         false,
			HCenter:         false,
			Top:             0,
			Left:            0,
			ID:              "",
			TextAlign:       TextAlign_Left,
			Writable:        false,
			TextType:        InputType_Text,
			ZIndex:          0,
		}
	}
}

// --------------------
//  Attributes Methods
// --------------------

func (a *Attributes) Parse(rawAttrs RawAttributeList) {
	for _, rawAttr := range rawAttrs {
		a.AddRaw(rawAttr[0], rawAttr[1])
	}
}

func (a *Attributes) AddRaw(attr string, value string) {
	switch AttrName(attr) {
	case AttrName_Display:
		a.Display = stringToDisplay(value)
	case AttrName_Position:
		a.Position = stringToPosition(value)
	case AttrName_Flex:
		a.Flex = stringToInt(value)
	case AttrName_FlexDirection:
		a.FlexDirection = stringToFlexDirection(value)
	case AttrName_Focusable:
		a.Focusable = stringToBool(value)
	case AttrName_Color:
		a.Color = stringToInt(value)
	case AttrName_BackGroundColor:
		a.BackGroundColor = stringToInt(value)
	case AttrName_Width:
		a.Width = stringToInt(value)
	case AttrName_MaxWidth:
		a.MaxWidth = stringToInt(value)
	case AttrName_MinWidth:
		a.MinWidth = stringToInt(value)
	case AttrName_Height:
		a.Height = stringToInt(value)
	case AttrName_MaxHeight:
		a.MaxHeight = stringToInt(value)
	case AttrName_MinHeight:
		a.MinHeight = stringToInt(value)
	case AttrName_Border:
		a.Border = stringToBool(value)
	case AttrName_VCenter:
		a.VCenter = stringToBool(value)
	case AttrName_HCenter:
		a.HCenter = stringToBool(value)
	case AttrName_Top:
		a.Top = stringToInt(value)
	case AttrName_Left:
		a.Left = stringToInt(value)
	case AttrName_ID:
		a.ID = value
	case AttrName_TextAlign:
		a.TextAlign = stringToTextAlign(value)
	case AttrName_Writable:
		a.Writable = stringToBool(value)
	case AttrName_TextType:
		a.TextType = stringToInputType(value)
	case AttrName_ZIndex:
		a.ZIndex = stringToUint8(value)
	case AttrName_Overflow:
		a.Overflow = stringToOverflow(value)
	}
}

func (a *Attributes) InheritFrom(parent *Attributes) {
	// a.Display = parent.Display
	// a.Position = parent.Position
	// a.Flex = parent.Flex
	// a.FlexDirection = parent.FlexDirection
	a.Color = parent.Color
	a.BackGroundColor = parent.BackGroundColor
	// a.Width = parent.Width
	// a.MaxWidth = parent.MaxWidth
	// a.MinWidth = parent.MinWidth
	// a.Height = parent.Height
	// a.MaxHeight = parent.MaxHeight
	// a.MinHeight = parent.MinHeight
	// a.Border = parent.Border
	// a.VCenter = parent.VCenter
	// a.HCenter = parent.HCenter
	// a.Top = parent.Top
	// a.Left = parent.Left
	// a.ID = parent.ID
	// a.TextAlign = parent.TextAlign
	// a.Writable = parent.Writable
	// a.TextType = parent.TextType
}

// --------------------
//   Attribute Types
// --------------------

type AttrName string

const (
	AttrName_Display         AttrName = "display"
	AttrName_Position        AttrName = "position"
	AttrName_Flex            AttrName = "flex"
	AttrName_FlexDirection   AttrName = "flex-direction"
	AttrName_Focusable       AttrName = "focusable"
	AttrName_Color           AttrName = "color"
	AttrName_BackGroundColor AttrName = "background-color"
	AttrName_Width           AttrName = "width"
	AttrName_MaxWidth        AttrName = "max-width"
	AttrName_MinWidth        AttrName = "min-width"
	AttrName_Height          AttrName = "height"
	AttrName_MaxHeight       AttrName = "max-height"
	AttrName_MinHeight       AttrName = "min-height"
	AttrName_Border          AttrName = "border"
	AttrName_VCenter         AttrName = "vcenter"
	AttrName_HCenter         AttrName = "hcenter"
	AttrName_Top             AttrName = "top"
	AttrName_Left            AttrName = "left"
	AttrName_ID              AttrName = "id"
	AttrName_TextAlign       AttrName = "text-align"
	AttrName_Writable        AttrName = "writable"
	AttrName_TextType        AttrName = "text-type"
	AttrName_ZIndex          AttrName = "z-index"
	AttrName_Overflow        AttrName = "overflow"
)

type Display uint8

const (
	Display_Block Display = iota
	Display_Flex
	Display_Absolute
)

type Position uint8

const (
	Position_Relative Position = iota
	Position_Absolute
)

type InputType uint8

const (
	InputType_Text InputType = iota
	InputType_Password
	InputType_Number
	InputType_Email
	InputType_Tel
	InputType_URL
)

type TextAlign uint8

const (
	TextAlign_Left TextAlign = iota
	TextAlign_Center
	TextAlign_Right
)

type FlexDirection uint8

const (
	FlexDirection_Row FlexDirection = iota
	FlexDirection_RowReverse
	FlexDirection_Column
	FlexDirection_ColumnReverse
)

type TextDecoration uint8

const (
	TextDecoration_None TextDecoration = iota
	TextDecoration_Underline
	TextDecoration_LineThrough
)

type Overflow uint8

const (
	Overflow_Hidden Overflow = iota
	Overflow_Scroll
)
