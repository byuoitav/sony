package bravia

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

func TestSupportedAPIInfo(t *testing.T) {
	is := is.New(t)
	disp.Log = zaptest.NewLogger(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := disp.getSupportedAPIInfo(ctx)
	is.NoErr(err)
}
