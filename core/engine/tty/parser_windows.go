//g0:build windows
// +build windows
 
package tty

import (
	"fmt"
)

// --------------------
//  Windows Event Parser
// --------------------

type windowsInputParser struct {}

func newInputParser() EventParser {
	return &windowsInputParser{}
}

func (p *windowsInputParser) Parse(b []byte) []Event {
	// TODO(maia): implement
	return nil
}


