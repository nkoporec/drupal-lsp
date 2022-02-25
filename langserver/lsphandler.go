package langserver

import (
	"context"
	"encoding/json"
	"log"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
)

// LspHandler ...
type LspHandler struct {
	jsonrpc2.EmptyHandler
	rootUri string
}

// InitializeParams
type InitializeParams struct {
	ProcessID int    `json:"processId,omitempty"`
	RootURI   string `json:"rootUri,omitempty"`
}

// NewLspHandler ...
func NewLspHandler() *LspHandler {
	return &LspHandler{}
}

func (h *LspHandler) handleTextDocumentCompletion(ctx context.Context, params *lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	result := make([]lsp.CompletionItem, 0, 200)
	return result, nil
}

// Deliver ...
func (h *LspHandler) Deliver(ctx context.Context, r *jsonrpc2.Request, delivered bool) bool {

	switch r.Method {
	case lsp.MethodInitialize:
		// Get params.
		var params InitializeParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			log.Fatal(err)
		}

		// Set rootUri.
		h.rootUri = params.RootURI

		// Run the index.
		indexer := NewIndexer(h.rootUri)
		go func() {
			indexer.Run()
		}()

		// Send back the response.
		err := r.Reply(ctx, lsp.InitializeResult{
			Capabilities: lsp.ServerCapabilities{
				CompletionProvider: &lsp.CompletionOptions{
					TriggerCharacters: []string{"."},
				},
				DefinitionProvider: false,
				HoverProvider:      false,
				SignatureHelpProvider: &lsp.SignatureHelpOptions{
					TriggerCharacters: []string{"(", ","},
				},
			},
		}, nil)

		if err != nil {
			panic(err)
		}

		return true
	}

	// Handle the request.
	switch r.Method {
	case lsp.MethodTextDocumentCompletion:
		var params lsp.CompletionParams
		json.Unmarshal(*r.Params, &params)
		items, err := h.handleTextDocumentCompletion(ctx, &params)
		r.Reply(ctx, items, err)
	}

	return true
}
