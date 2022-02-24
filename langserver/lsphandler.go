package langserver

import (
	"context"
	"encoding/json"
	"fmt"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
)

// LspHandler ...
type LspHandler struct {
	initialDiagnostics map[string][]lsp.Diagnostic
	baseLspHandler
	initialized bool
}

var (
	// ErrWalkAbort should be returned if a walk function should abort early
	ErrWalkAbort = fmt.Errorf("OK")
)

// NewLspHandler ...
func NewLspHandler() *LspHandler {
	logLv := LogLevelInfo
	return &LspHandler{
		baseLspHandler: baseLspHandler{
			LogLevel: logLv,
		},
		initialized: false,
	}
}

func completionItemFromSymbol(s Symbol) (lsp.CompletionItem, error) {
	kind, err := completionItemKindForSymbol(s)
	if err != nil {
		return lsp.CompletionItem{}, err
	}
	return lsp.CompletionItem{
		Kind:   kind,
		Label:  s.Name(),
		Detail: s.String(),
		Documentation: lsp.MarkupContent{
			Kind:  lsp.PlainText,
			Value: s.Documentation(),
		},
	}, nil
}

func completionItemKindForSymbol(s Symbol) (lsp.CompletionItemKind, error) {
	switch s.(type) {
	case VariableSymbol:
		return lsp.VariableCompletion, nil
	case ConstantSymbol:
		return lsp.ConstantCompletion, nil
	case FunctionSymbol:
		return lsp.FunctionCompletion, nil
	case ClassSymbol:
		return lsp.ClassCompletion, nil
	case ProtoTypeOrInstanceSymbol:
		return lsp.ClassCompletion, nil
	}
	return lsp.CompletionItemKind(-1), fmt.Errorf("Symbol not found")
}

func (h *LspHandler) handleTextDocumentCompletion(ctx context.Context, params *lsp.CompletionParams) ([]lsp.CompletionItem, error) {
	result := make([]lsp.CompletionItem, 0, 200)
	return result, nil
}

// Deliver ...
func (h *LspHandler) Deliver(ctx context.Context, r *jsonrpc2.Request, delivered bool) bool {
	h.LogDebug("Requested '%s'\n", r.Method)

	switch r.Method {
	case lsp.MethodInitialize:
		h.LogDebug("initialized")
		// @TODO:
		// Run the index.

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
			h.LogError("Error sending response: %s\n", err)
		}

		h.initialized = true

		return true
	}

	// If something went wrong, return false.
	if !h.initialized {
		if !r.IsNotify() {
			r.Reply(
				ctx,
				nil,
				jsonrpc2.Errorf(jsonrpc2.ServerNotInitialized, "Not initialized yet"),
			)
		}

		return false
	}

	// Recover if something bad happens in the handlers...
	defer func() {
		err := recover()
		if err != nil {
			h.LogWarn("Recovered from panic at %s: %v\n", r.Method, err)
		}
	}()

	// Handle the request.
	switch r.Method {
	case lsp.MethodTextDocumentCompletion:
		var params lsp.CompletionParams
		json.Unmarshal(*r.Params, &params)
		items, err := h.handleTextDocumentCompletion(ctx, &params)
		h.replyEither(ctx, r, items, err)

		return h.baseLspHandler.Deliver(ctx, r, false)
	}

	return true
}
