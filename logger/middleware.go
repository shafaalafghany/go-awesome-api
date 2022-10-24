package logger

import "context"

type reqIdCtxKey struct{}

const ReqIdKey = "requestId"

func ReqID(ctx context.Context) string {
	if id, ok := ctx.Value(reqIdCtxKey{}).(string); ok {
		return id
	}
	return ""
}
