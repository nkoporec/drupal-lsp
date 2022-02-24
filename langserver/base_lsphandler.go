package langserver

import (
	"context"
	"fmt"
	"os"

	"go.lsp.dev/jsonrpc2"
)

type logLevel int

var (
	// LogLevelDebug ...
	LogLevelDebug logLevel = 0
	// LogLevelInfo ...
	LogLevelInfo logLevel = 1
	// LogLevelWarn ...
	LogLevelWarn logLevel = 2
	// LogLevelErr ...
	LogLevelErr logLevel = 3
)

type baseLspHandler struct {
	jsonrpc2.EmptyHandler
	LogLevel logLevel
}

func (h *baseLspHandler) replyEither(ctx context.Context, r *jsonrpc2.Request, result interface{}, err error) {
	if err != nil {
		r.Reply(ctx, nil, err)
	} else {
		r.Reply(ctx, result, nil)
	}
}

func (h *baseLspHandler) LogDebug(format string, params ...interface{}) {
	if h.LogLevel <= LogLevelDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format, params...)
	}
}

func (h *baseLspHandler) LogInfo(format string, params ...interface{}) {
	if h.LogLevel <= LogLevelInfo {
		fmt.Fprintf(os.Stderr, "INFO: "+format, params...)
	}
}

func (h *baseLspHandler) LogWarn(format string, params ...interface{}) {
	if h.LogLevel <= LogLevelWarn {
		fmt.Fprintf(os.Stderr, "WARN: "+format, params...)
	}
}

func (h *baseLspHandler) LogError(format string, params ...interface{}) {
	if h.LogLevel <= LogLevelErr {
		fmt.Fprintf(os.Stderr, "ERROR: "+format, params...)
	}
}
