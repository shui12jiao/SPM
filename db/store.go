package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
}

// SQLStore提供了执行SQL查询和事务的所有函数
type SQLStore struct {
	*Queries
	db *sql.DB
}

// 返回一个新的Store对象
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx 执行一个函数，这个函数会传入一个Queries对象执行数据库操作
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
