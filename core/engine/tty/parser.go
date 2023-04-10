package tty

// --------------------
//    Input Parser
// --------------------

type InputParser interface {
	// Parse parses the given bytes and returns a slice of events.
	Parse([]byte) []Event
}
