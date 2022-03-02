package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"drupal-lsp/langserver"

	"go.lsp.dev/jsonrpc2"
)

var (
	mode         = flag.String("mode", "stdio", "communication mode (stdio|tcp|websocket)")
	logfile      = flag.String("logfile", "", "log to this file (in addition to stderr)")
	printVersion = flag.Bool("version", false, "print version and exit")
)

// Update this when we do a release.
const VERSION = "0.0.1"

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Printf("drupal-lsp version: %s", VERSION)
		os.Exit(0)
	}

	// @todo: implement other modes.
	if *mode != "stdio" {
		fmt.Println("We only support stdio mode")
		os.Exit(0)
	}

	if *logfile != "" {
		log.Println("Logging enabled")
		log.Println(fmt.Sprintf("Logging to file: %s", *logfile))
		f, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open logfile: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	log.Printf(
		fmt.Sprintf("Starting Drupal Language Server in %s mode ...", *mode),
	)

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
