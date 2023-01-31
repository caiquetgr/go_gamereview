package database

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

type BaseDbModel struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	CreatedAt time.Time `bun:",notnull,default:current_timestamp`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp`
}

func (b *BaseDbModel) BeforeAppendmodel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		b.UpdatedAt = time.Now()
	}
	return nil
}
