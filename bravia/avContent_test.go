package bravia

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

// TestAudioVideoInput does the following:
// 1. turn on the display
// 2. set the input to hdmi 1
// 3. set the input to hdmi 2
// 4. turn off the display
// Steps 2 and 3 are verified to make sure that the correct input
// is returned.
func TestAudioVideoInput(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, true))

	input := "hdmi?port=3"
	is.NoErr(disp.SetAudioVideoInput(ctx, "", input))

	inputs, err := disp.AudioVideoInputs(ctx)
	is.NoErr(err)
	is.True(inputs[""] == input)

	input = "hdmi?port=2"
	is.NoErr(disp.SetAudioVideoInput(ctx, "", input))

	inputs, err = disp.AudioVideoInputs(ctx)
	is.NoErr(err)
	is.True(inputs[""] == input)

	is.NoErr(disp.SetPower(ctx, false))
}

func TestCurrentExternalInputsStatus(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	list, err := disp.getCurrentExternalInputsStatus(ctx)
	is.NoErr(err)

	for _, item := range list {
		t.Logf("%+v", item)
	}
}
