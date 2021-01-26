package bravia

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestVolumeInformation(t *testing.T) {
	t.SkipNow()
	is := is.New(t)

	d := &Display{
		Address:      "ITB-2033-D1.byu.edu",
		PreSharedKey: _preSharedKey,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	is.NoErr(d.SetPower(ctx, true))

	list, err := d.getVolumeInformation(ctx)
	is.NoErr(err)

	for _, item := range list {
		t.Logf("%+v", item)
	}

	is.NoErr(d.SetPower(ctx, false))
}

func TestVolume(t *testing.T) {
	is := is.New(t)

	d := &Display{
		Address:      "ITB-2033-D1.byu.edu",
		PreSharedKey: _preSharedKey,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	is.NoErr(d.SetPower(ctx, true))

	for i := 0; i < 3; i++ {
		vol := rand.Intn(101)
		is.NoErr(d.SetVolume(ctx, "speaker", vol))

		vols, err := d.Volumes(ctx, []string{"speaker"})
		is.NoErr(err)

		v, ok := vols["speaker"]
		is.True(ok)
		is.True(v == vol)
	}

	is.NoErr(d.SetPower(ctx, false))
}
