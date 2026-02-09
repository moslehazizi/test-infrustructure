package repository

import "context"

type DatabaseMetadata interface {
	GetAll(ctx context.Context) ([]string, error)
	GetTablesByDBName(ctx context.Context, dbName string) ([]string, error)
}
