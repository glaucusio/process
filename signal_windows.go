// +build windows

package process

import (
	"os"
)

var defaultSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
}
