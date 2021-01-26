package bravia

import (
	"os"
	"testing"
	"time"
)

var disp *Display

func TestMain(m *testing.M) {
	disp = &Display{
		Address:      "ITB-2033-D1.byu.edu",
		PreSharedKey: os.Getenv("BRAVIA_PSK"),
		RequestDelay: 300 * time.Millisecond,
	}

	os.Exit(m.Run())
}
