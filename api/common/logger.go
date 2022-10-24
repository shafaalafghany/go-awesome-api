package common

import (
	"awesome-api/logger"
	"context"

	"github.com/rs/zerolog"
)

type WrapperZlog struct {
	*zerolog.Logger
}

func attachReqId(event *zerolog.Event, ctx context.Context) *zerolog.Event {
	return event.Str(logger.ReqIdKey, logger.ReqID(ctx))
}

func (w *WrapperZlog) Info(ctx context.Context) *zerolog.Event {
	return attachReqId(w.Logger.Info(), ctx)
}

func (w *WrapperZlog) Warn(ctx context.Context) *zerolog.Event {
	return attachReqId(w.Logger.Warn(), ctx)
}

func (w *WrapperZlog) Error(ctx context.Context) *zerolog.Event {
	return attachReqId(w.Logger.Error(), ctx)
}
