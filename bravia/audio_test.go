package bravia

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

func TestVolumeInformation(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, true))

	list, err := disp.getVolumeInformation(ctx)
	is.NoErr(err)

	for _, item := range list {
		t.Logf("%+v", item)
	}

	is.NoErr(disp.SetPower(ctx, false))
}

func TestVolume(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, true))

	for i := 0; i < 3; i++ {
		vol := rand.Intn(101)
		is.NoErr(disp.SetVolume(ctx, "speaker", vol))

		vols, err := disp.Volumes(ctx, []string{"speaker"})
		is.NoErr(err)

		v, ok := vols["speaker"]
		is.True(ok)
		is.True(v == vol)
	}

	is.NoErr(disp.SetPower(ctx, false))
}

func TestMute(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	test := func(muted bool) {
		is.NoErr(disp.SetMute(ctx, "speaker", muted))

		mutes, err := disp.Mutes(ctx, []string{"speaker"})
		is.NoErr(err)

		m, ok := mutes["speaker"]
		is.True(ok)
		is.True(m == muted)
	}

	is.NoErr(disp.SetPower(ctx, true))

	test(true)
	test(false)

	is.NoErr(disp.SetPower(ctx, false))
}

func TestOffVolumes(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, false))

	vols, err := disp.Volumes(ctx, []string{"speaker"})
	is.NoErr(err)
	is.True(len(vols) == 0)
}

func TestOffMutes(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	is.NoErr(disp.SetPower(ctx, false))

	mutes, err := disp.Mutes(ctx, []string{"speaker"})
	is.NoErr(err)
	is.True(len(mutes) == 0)
}
