package langserver

import (
	"context"
	"drupal-lsp/langserver/parser"
	"encoding/json"
	"log"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
)

// LspHandler ...
type LspHandler struct {
	jsonrpc2.EmptyHandler
	rootUri                 string
	Indexer                 *Indexer
	Buffer                  *Buffer
	TextDocumentSyncHandler jsonrpc2.Handler
}

// InitializeParams
type InitializeParams struct {
	ProcessID int    `json:"processId,omitempty"`
	RootURI   string `json:"rootUri,omitempty"`
}

// NewLspHandler ...
func NewLspHandler() *LspHandler {
	return &LspHandler{
		TextDocumentSyncHandler: &TextDocumentSyncHandler{},
	}
}

func (h *LspHandler) handleTextDocumentCompletion(ctx context.Context, params *lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	result := make([]lsp.CompletionItem, 0, 200)

	for _, item := range h.Indexer.Services {
		serviceCompletion, err := parser.CompletionItemForService(item)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, serviceCompletion)
	}

	return result, nil
}

// Deliver ...
func (h *LspHandler) Deliver(ctx context.Context, r *jsonrpc2.Request, delivered bool) bool {

	log.Println("Method: ", r.Method)
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
		// @todo make it async.
		indexer := NewIndexer(h.rootUri)
		indexer.Run()

		// Set the indexer.
		h.Indexer = indexer

		// Buffer.
		h.Buffer = NewBuffer()

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
	case lsp.MethodTextDocumentDidOpen:
		var params lsp.DidOpenTextDocumentParams
		json.Unmarshal(*r.Params, &params)
		documentUri := UriToFilename(params.TextDocument.URI)
		if documentUri != "" {
			h.Buffer.UpdateBufferDoc(documentUri, params.TextDocument.Text)
		}
	case lsp.MethodTextDocumentDidChange:
		var params lsp.DidChangeTextDocumentParams
		json.Unmarshal(*r.Params, &params)
		documentUri := UriToFilename(params.TextDocument.URI)
		if documentUri != "" && len(params.ContentChanges) > 0 {
			h.Buffer.UpdateBufferDoc(documentUri, params.ContentChanges[0].Text)
		}
	case lsp.MethodTextDocumentDidSave:
		var params lsp.DidSaveTextDocumentParams
		json.Unmarshal(*r.Params, &params)
		documentUri := UriToFilename(params.TextDocument.URI)
		if documentUri != "" {
			h.Buffer.UpdateBufferDoc(documentUri, params.Text)
		}
	case lsp.MethodTextDocumentCompletion:
		var params lsp.CompletionParams
		json.Unmarshal(*r.Params, &params)
		items, err := h.handleTextDocumentCompletion(ctx, &params)
		r.Reply(ctx, items, err)
	}

	return true
}
