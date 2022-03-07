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
	Indexer *Indexer
	Buffer  *Buffer
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

	// Get the doc.
	doc := h.Buffer.GetBufferDoc(UriToFilename(params.TextDocument.URI))
	method, err := doc.GetMethodCall(params.Position)
	if err != nil {
		return result, err
	}

	// Get all parsers.
	parsers := h.Indexer.Parsers
	for _, parser := range parsers {
		// Get the method call.
		methods := parser.Methods()
		if inSlice(methods, method) {
			for _, def := range parser.GetDefinitions() {
				log.Println(def.Name)
				completion, err := parser.CompletionItem(def)
				if err != nil {
					return result, err
				}

				result = append(result, completion)
			}
		}
	}

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
				CompletionProvider: &lsp.CompletionOptions{},
				DefinitionProvider: false,
				HoverProvider:      false,
				TextDocumentSync: lsp.TextDocumentSyncOptions{
					Change:    float64(lsp.Full),
					OpenClose: true,
					Save: &lsp.SaveOptions{
						IncludeText: true,
					},
				},
			},
		}, nil)

		if err != nil {
			log.Fatal(err)
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

// @todo Move it to the appropriate place.
func inSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
