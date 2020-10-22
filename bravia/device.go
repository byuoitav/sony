package bravia

import "go.uber.org/zap"

type TV struct {
	Address string
	PSK     string
	Log     *zap.Logger
}
