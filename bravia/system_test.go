package bravia

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

var _preSharedKey = os.Getenv("BRAVIA_PSK")

// TestPower turns the tv on and then off, verifying that
// it works after each step.
func TestPower(t *testing.T) {
	is := is.New(t)

	d := &Display{
		Address:      "ITB-2033-D1.byu.edu",
		PreSharedKey: _preSharedKey,
		Log:          zaptest.NewLogger(t),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	is.NoErr(d.SetPower(ctx, true))

	pow, err := d.Power(ctx)
	is.NoErr(err)
	is.True(pow)

	is.NoErr(d.SetPower(ctx, false))

	pow, err = d.Power(ctx)
	is.NoErr(err)
	is.True(!pow)
}

func TestBlank(t *testing.T) {
	is := is.New(t)

	d := &Display{
		Address:      "ITB-2033-D1.byu.edu",
		PreSharedKey: _preSharedKey,
		Log:          zaptest.NewLogger(t),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	is.NoErr(d.SetPower(ctx, true))
	is.NoErr(d.SetBlank(ctx, true))

	blanked, err := d.Blank(ctx)
	is.NoErr(err)
	is.True(blanked)

	is.NoErr(d.SetBlank(ctx, false))

	blanked, err = d.Blank(ctx)
	is.NoErr(err)
	is.True(!blanked)

	is.NoErr(d.SetPower(ctx, false))
}
