package helper

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func InjectTX(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, db)
}

func ExtractTX(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return nil
}
