package tty

// -----------------
//      Event
// -----------------

type Event interface {
	// Type returns the type of the event.
	Type() EventType
}

type EventType int

const (
	// EventTypeResize is the type of a resize event.
	EventType_Resize EventType = iota
	// EventTypeKeyboard is the type of a keyboard event.
	EventType_Keyboard
	// EventTypeMouse is the type of a mouse event.
	EventType_Mouse
)

// -----------------
//   ResizeEvent
// -----------------

type ResizeEvent struct {
	Width  int
	Height int
}

func (e ResizeEvent) Type() EventType {
	return EventType_Resize
}

// -----------------
//  KeyaboardEvent
// -----------------

type KeyboardEvent struct {
	Key       Key
	Modifiers Modifiers
}

func newKeyboardEvent(key Key, mod ...Modifiers) KeyboardEvent {
	k := KeyboardEvent{
		Key: key,
	}
	for _, m := range mod {
		k.Modifiers |= m
	}
	return k
}

func (e KeyboardEvent) SetModifiers(modifiers Modifiers) KeyboardEvent {
	e.Modifiers = modifiers
	return e
}

func (e KeyboardEvent) Type() EventType {
	return EventType_Keyboard
}

func (e KeyboardEvent) Is(key Key, mdfr ...Modifiers) bool {
	for _, m := range mdfr {
		if m&e.Modifiers > 0 || m == e.Modifiers {
			continue
		}
		return false
	}

	return e.Key.Is(key)
}

type Key interface {
	Type() KeyType
	Is(Key) bool
}

type KeyType int

const (
	// KeyTypeByte is the type of a rune key.
	KeyTypeByte KeyType = iota
	// KeyTypeSpecial is the type of a special key.
	KeyTypeSpecial
)

type ByteKey byte

func NewByteKey(r byte) ByteKey {
	return ByteKey(r)
}

func (k ByteKey) Type() KeyType {
	return KeyTypeByte
}

func (k ByteKey) Is(key Key) bool {
	if key.Type() != KeyTypeByte {
		return false
	}
	return rune(k) == rune(key.(ByteKey))
}

type SpecialKey string

func NewSpecialKey(key string) SpecialKey {
	return SpecialKey(key)
}

func (k SpecialKey) Type() KeyType {
	return KeyTypeSpecial
}

func (k SpecialKey) Is(key Key) bool {
	if key.Type() != KeyTypeSpecial {
		return false
	}
	return string(k) == string(key.(SpecialKey))
}

const (
	SpecialKey_Escape SpecialKey = "\x1b"
	SpecialKey_Enter  SpecialKey = "\r"
	SpecialKey_Up     SpecialKey = "\x1b[A"
	SpecialKey_Down   SpecialKey = "\x1b[B"
	SpecialKey_Right  SpecialKey = "\x1b[C"
	SpecialKey_Left   SpecialKey = "\x1b[D"
	SpecialKey_Insert SpecialKey = "\x1b[2~"
	SpecialKey_Delete SpecialKey = "\x1b[3~"
	SpecialKey_Home   SpecialKey = "\x1b[1~"
	SpecialKey_End    SpecialKey = "\x1b[4~"
	SpecialKey_PageUp SpecialKey = "\x1b[5~"
	SpecialKey_F1     SpecialKey = "\x1bOP"
	SpecialKey_F2     SpecialKey = "\x1bOQ"
	SpecialKey_F3     SpecialKey = "\x1bOR"
	SpecialKey_F4     SpecialKey = "\x1bOS"
	SpecialKey_F5     SpecialKey = "\x1b[15~"
	SpecialKey_F6     SpecialKey = "\x1b[17~"
	SpecialKey_F7     SpecialKey = "\x1b[18~"
	SpecialKey_F8     SpecialKey = "\x1b[19~"
	SpecialKey_F9     SpecialKey = "\x1b[20~"
	SpecialKey_F10    SpecialKey = "\x1b[21~"
	SpecialKey_F11    SpecialKey = "\x1b[23~"
	SpecialKey_F12    SpecialKey = "\x1b[24~"
)

type Modifiers byte

const (
	Modifier_None  Modifiers = 0
	Modifier_Shift Modifiers = 1 << iota
	Modifier_Ctrl
	Modifier_Alt
)

// -----------------
//   MouseEvent
// -----------------

type MouseEvent struct {
	// X is the x coordinate of the mouse event.
	X int
	// Y is the y coordinate of the mouse event.
	Y int
	// Button is the button that was pressed.
	Button MouseButton
	// Modifiers is the set of modifiers that were pressed.
	Modifiers Modifiers
}

func (e MouseEvent) Type() EventType {
	return EventType_Mouse
}

func (e MouseEvent) SetModifiers(modifiers Modifiers) MouseEvent {
	e.Modifiers = modifiers
	return e
}

type MouseButton int

const (
	MouseButton_None MouseButton = iota
	MouseButton_Left
	MouseButton_Right
	MouseButton_Middle
)
