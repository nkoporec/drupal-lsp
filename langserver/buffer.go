package langserver

import (
	"sync"
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

func (b *Buffer) UpdateBufferDoc(documentURI string, buf string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	d := &Document{
		URI:  documentURI,
		Text: buf,
	}

	b.Documents[documentURI] = *d
}

func (b *Buffer) GetBufferDoc(documentURI string) *Document {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	if doc, ok := b.Documents[documentURI]; ok {
		return &doc
	}

	return nil
}
