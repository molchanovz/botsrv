package rpc

import (
	"net/http"

	"botsrv/pkg/db"

	"github.com/vmkteam/embedlog"
	zm "github.com/vmkteam/zenrpc-middleware"
	"github.com/vmkteam/zenrpc/v2"
)

var (
	slackWebhook = "https://hooks.slack.com/services/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"
	token        = "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"
)

var (
	ErrNotImplemented = zenrpc.NewStringError(http.StatusInternalServerError, "not implemented")
	ErrInternal       = zenrpc.NewStringError(http.StatusInternalServerError, "internal error")
)

var allowDebugFn = func() zm.AllowDebugFunc {
	return func(req *http.Request) bool {
		return req != nil && req.FormValue("__level") == "5"
	}
}

//go:generate zenrpc

// New returns new zenrpc Server.
func New(dbo db.DB, logger embedlog.Logger, isDevel bool) zenrpc.Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
	})

	rpc.Use(
		zm.WithDevel(isDevel),
		zm.WithHeaders(),
		zm.WithSentry(zm.DefaultServerName),
		zm.WithNoCancelContext(),
		zm.WithMetrics(zm.DefaultServerName),
		zm.WithTiming(isDevel, allowDebugFn()),
		zm.WithSQLLogger(dbo.DB, isDevel, allowDebugFn(), allowDebugFn()),
	)

	rpc.Use(
		zm.WithSLog(logger.Print, zm.DefaultServerName, nil),
		zm.WithErrorSLog(logger.Print, zm.DefaultServerName, nil),
	)

	// services
	rpc.RegisterAll(map[string]zenrpc.Invoker{
		// "sample": NewSampleService(db, logger),
	})

	return rpc
}

//nolint:unused
func newInternalError(err error) *zenrpc.Error {
	return zenrpc.NewError(http.StatusInternalServerError, err)
}
