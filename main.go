package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"drupal-lsp/langserver"

	"go.lsp.dev/jsonrpc2"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	lspHandler := langserver.NewLspHandler()
	connectLanguageServer(os.Stdin, os.Stdout, lspHandler).Run(ctx)
}

func connectLanguageServer(in io.Reader, out io.Writer, handlers ...jsonrpc2.Handler) *jsonrpc2.Conn {
	bufStream := jsonrpc2.NewStream(in, out)
	rootConn := jsonrpc2.NewConn(bufStream)

	for _, h := range handlers {
		rootConn.AddHandler(h)
	}
	return rootConn
}
