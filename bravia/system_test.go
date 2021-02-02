package bravia

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

// TestPower turns the tv on and then off, verifying that
// it works after each step.
func TestPower(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, true))

	pow, err := disp.Power(ctx)
	is.NoErr(err)
	is.True(pow)

	is.NoErr(disp.SetPower(ctx, false))

	pow, err = disp.Power(ctx)
	is.NoErr(err)
	is.True(!pow)
}

func TestBlank(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, true))
	is.NoErr(disp.SetBlank(ctx, true))

	blanked, err := disp.Blank(ctx)
	is.NoErr(err)
	is.True(blanked)

	is.NoErr(disp.SetBlank(ctx, false))

	blanked, err = disp.Blank(ctx)
	is.NoErr(err)
	is.True(!blanked)

	is.NoErr(disp.SetPower(ctx, false))
}

func TestInfo(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := disp.Info(ctx)
	is.NoErr(err)
	is.True(info != nil)

	t.Logf("%+v", info)
}

func TestHealth(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := disp.Healthy(ctx)
	is.NoErr(err)
}
