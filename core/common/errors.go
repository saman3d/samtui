package common

import "errors"

var (
	ErrUnexpectedEndTag = NewError("unexpected end tag")
)

func NewError(text string) error {
	return errors.New(text)
}
