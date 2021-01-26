package bravia

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestSupportedAPIInfo(t *testing.T) {
	is := is.New(t)

	d := &Display{
		Address:      "ITB-2033-D1.byu.edu",
		PreSharedKey: _preSharedKey,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := d.getSupportedAPIInfo(ctx)
	is.NoErr(err)
}
