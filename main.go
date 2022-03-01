package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"drupal-lsp/langserver"

	"go.lsp.dev/jsonrpc2"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	// @todo: Add a flag to allow the user to specify the file
	f, err := os.OpenFile("/home/nkoporec/personal/drupal-lsp/drupal-lsp.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("Starting Drupal Language Server ...")

	lspHandler := langserver.NewLspHandler()
	connectLanguageServer(os.Stdin, os.Stdout, lspHandler.TextDocumentSyncHandler, lspHandler).Run(ctx)
}

func connectLanguageServer(in io.Reader, out io.Writer, handlers ...jsonrpc2.Handler) *jsonrpc2.Conn {
	bufStream := jsonrpc2.NewStream(in, out)
	rootConn := jsonrpc2.NewConn(bufStream)

	for _, h := range handlers {
		rootConn.AddHandler(h)
	}

	return rootConn
}
