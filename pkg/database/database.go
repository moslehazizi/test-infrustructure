package database

import (
	"context"
)

type Database interface {
	Ping() error
	Init() error
	Close() error
	WithContext(ctx context.Context) Database
	Begin() Database
	Commit() error
	Rollback() error
	Client() Database
}

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const (
	ContextKeyDBTx ContextKey = "tx"
)

type DBInitializerFn func(cfg any) (Database, error)
