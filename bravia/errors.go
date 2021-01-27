package bravia

import "fmt"

const (
	_displayOff = 40005
)

type Error struct {
	code   int
	reason string
}

func (e *Error) Error() string {
	if e.code < 0 {
		return e.reason
	}

	return fmt.Sprintf("%v: %v", e.code, e.reason)
}
