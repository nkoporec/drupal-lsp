package langserver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nkoporec/drupal-lsp/utils"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
	"go.lsp.dev/uri"
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
		if utils.InSlice(methods, method) {
			for _, def := range parser.GetDefinitions() {
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

func (h *LspHandler) handleGoToDefinition(ctx context.Context, params *lsp.TextDocumentPositionParams, i *Indexer) ([]lsp.Location, error) {
	// @todo: Implement.
	result := make([]lsp.Location, 0, 200)

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
		if utils.InSlice(methods, method) {
			methodParams, err := doc.GetMethodParams(params.Position)
			if err != nil {
				return result, err
			}

			definitions := parser.GetGoToDefinition(methodParams)
			for _, class := range definitions {
				for _, item := range i.PhpClasses {
					if item.Namespace == class {
						// @todo: Implement start/end range.
						result = append(result, lsp.Location{
							URI:   uri.File(item.Path),
							Range: lsp.Range{},
						})
					}
				}
			}
		}
	}

	return result, nil
}

func (h *LspHandler) handleHoverDefinition(ctx context.Context, params *lsp.TextDocumentPositionParams, i *Indexer) (lsp.Hover, error) {
	result := lsp.Hover{}

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
		if utils.InSlice(methods, method) {
			methodParams, err := doc.GetMethodParams(params.Position)
			if err != nil {
				return result, err
			}

			definitions := parser.GetGoToDefinition(methodParams)
			for _, class := range definitions {
				for _, item := range i.PhpClasses {
					if item.Namespace == class {
						result = lsp.Hover{
							Contents: lsp.MarkupContent{
								Kind:  "php",
								Value: item.Namespace + "\n\n" + item.Path + "\n\n" + item.Description,
							},
						}
					}
				}
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
				DefinitionProvider: true,
				HoverProvider:      true,
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
			h.Buffer.UpdateBufferDoc(documentUri, params.TextDocument.Text, ctx, r, h.Indexer)
		}
	case lsp.MethodTextDocumentDidChange:
		var params lsp.DidChangeTextDocumentParams
		json.Unmarshal(*r.Params, &params)
		documentUri := UriToFilename(params.TextDocument.URI)
		if documentUri != "" && len(params.ContentChanges) > 0 {
			h.Buffer.UpdateBufferDoc(documentUri, params.ContentChanges[0].Text, ctx, r, h.Indexer)
		}
	case lsp.MethodTextDocumentDidSave:
		var params lsp.DidSaveTextDocumentParams
		json.Unmarshal(*r.Params, &params)
		documentUri := UriToFilename(params.TextDocument.URI)
		if documentUri != "" {
			h.Buffer.UpdateBufferDoc(documentUri, params.Text, ctx, r, h.Indexer)
		}
	case lsp.MethodTextDocumentCompletion:
		var params lsp.CompletionParams
		json.Unmarshal(*r.Params, &params)
		items, err := h.handleTextDocumentCompletion(ctx, &params)
		r.Reply(ctx, items, err)
	case lsp.MethodTextDocumentDefinition:
		var params lsp.TextDocumentPositionParams
		json.Unmarshal(*r.Params, &params)
		found, err := h.handleGoToDefinition(ctx, &params, h.Indexer)
		r.Reply(ctx, found, err)
	case lsp.MethodTextDocumentHover:
		var params lsp.TextDocumentPositionParams
		json.Unmarshal(*r.Params, &params)
		found, err := h.handleHoverDefinition(ctx, &params, h.Indexer)
		r.Reply(ctx, found, err)
	}

	return true
}
