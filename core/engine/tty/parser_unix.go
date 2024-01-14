//go:build unix
// +build unix

package tty

const (
	ResetEscape = "\x1b[0m"
)

type unixInputtParser struct{}

func newInputParser() InputParser {
	return &unixInputtParser{}
}

func (p *unixInputtParser) Parse(b []byte) []Event {
	return []Event{p.parseKey(b)}
}

func (p *unixInputtParser) parseKey(b []byte) Event {
	if len(b) == 1 || b[1] == 0 {
		if b[0] == 13 {
			return newKeyboardEvent(NewByteKey(b[0]))
		}
		if b[0] < 36 {
			return newKeyboardEvent(NewByteKey(b[0]+96), Modifier_Ctrl)
		}
		return newKeyboardEvent(NewByteKey(b[0]), 0)
	}

	if len(b) == 2 || b[2] == 0 && b[0] == 27 {
		return newKeyboardEvent(NewByteKey(b[1]), Modifier_Alt)
	}

	if len(b) == 6 || b[6] == 0 && b[0] == 27 && b[1] == 91 && b[2] == 49 && b[3] == 59 {
		return newKeyboardEvent(NewSpecialKey(string([]byte{27, 91, b[5]})), byteToModifier(b[4]))
	}

	if len(b) == 7 || b[7] == 0 && b[0] == 27 && b[1] == 91 && b[2] == 50 && b[3] == 51 && b[4] == 59 {
		return newKeyboardEvent(NewSpecialKey(string([]byte{27, 91, 50, 52, b[6]})), byteToModifier(b[5]))
	}

	return newKeyboardEvent(NewSpecialKey(string(cleanByteArray(b))))
}

func (p *unixInputtParser) modifiers(b []byte) Event {
	var m Modifiers

	if len(b) > 1 && b[0] == 0x1b {
		m |= Modifier_Alt
	}

	if len(b) == 1 && b[0] < 32 {
		m |= Modifier_Ctrl
	}

	return nil
}

func cleanByteArray(b []byte) []byte {
	for i, v := range b {
		if v == 0 {
			return b[:i]
		}
	}
	return b
}

func byteToModifier(n byte) Modifiers {
	var md Modifiers
	switch n {
	case 50:
		md = Modifier_Shift
	case 51:
		md = Modifier_Alt
	case 52:
		md = Modifier_Alt | Modifier_Shift
	case 53:
		md = Modifier_Ctrl
	case 54:
		md = Modifier_Ctrl | Modifier_Shift
	case 55:
		md = Modifier_Ctrl | Modifier_Alt
	case 56:
		md = Modifier_Ctrl | Modifier_Shift | Modifier_Alt
	}
	return md
}
