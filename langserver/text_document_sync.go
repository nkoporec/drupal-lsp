package langserver

import (
	"context"
	"log"

	"go.lsp.dev/jsonrpc2"
)

type TextDocumentSyncHandler struct {
	jsonrpc2.EmptyHandler
}

// Deliver ...
func (h *TextDocumentSyncHandler) Deliver(ctx context.Context, r *jsonrpc2.Request, delivered bool) bool {
	log.Println(r.Method)
	return true
}
