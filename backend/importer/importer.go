package importer

import "context"

// dummy thingy

type Importer struct {
	ctx context.Context
}

func NewImporter(ctx context.Context) *Importer {
	return &Importer{ctx: ctx}
}
