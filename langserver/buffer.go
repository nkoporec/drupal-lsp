package langserver

import (
	"context"
	"log"
	"sync"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type Buffer struct {
	Documents map[string]Document
	mtx       sync.RWMutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		Documents: make(map[string]Document, 0),
	}
}

func (b *Buffer) UpdateBufferDoc(documentURI string, buf string, ctx context.Context, r *jsonrpc2.Request) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	d := &Document{
		URI:  documentURI,
		Text: buf,
	}
	b.Documents[documentURI] = *d

	diagnostics, err := d.GetDiagnostics()
	if err != nil {
		log.Fatal(err)
	}

	r.Conn().Notify(ctx, lsp.MethodTextDocumentPublishDiagnostics, lsp.PublishDiagnosticsParams{
		URI:         lsp.DocumentURI(uri.File(documentURI)),
		Diagnostics: diagnostics,
	})
}

func (b *Buffer) GetBufferDoc(documentURI string) *Document {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	if doc, ok := b.Documents[documentURI]; ok {
		return &doc
	}

	return nil
}
