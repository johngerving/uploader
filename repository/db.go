// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createPartStmt, err = db.PrepareContext(ctx, createPart); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePart: %w", err)
	}
	if q.createUploadStmt, err = db.PrepareContext(ctx, createUpload); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUpload: %w", err)
	}
	if q.findUploadByIdStmt, err = db.PrepareContext(ctx, findUploadById); err != nil {
		return nil, fmt.Errorf("error preparing query FindUploadById: %w", err)
	}
	if q.findUploadPartsByIdStmt, err = db.PrepareContext(ctx, findUploadPartsById); err != nil {
		return nil, fmt.Errorf("error preparing query FindUploadPartsById: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createPartStmt != nil {
		if cerr := q.createPartStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPartStmt: %w", cerr)
		}
	}
	if q.createUploadStmt != nil {
		if cerr := q.createUploadStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUploadStmt: %w", cerr)
		}
	}
	if q.findUploadByIdStmt != nil {
		if cerr := q.findUploadByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing findUploadByIdStmt: %w", cerr)
		}
	}
	if q.findUploadPartsByIdStmt != nil {
		if cerr := q.findUploadPartsByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing findUploadPartsByIdStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                      DBTX
	tx                      *sql.Tx
	createPartStmt          *sql.Stmt
	createUploadStmt        *sql.Stmt
	findUploadByIdStmt      *sql.Stmt
	findUploadPartsByIdStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                      tx,
		tx:                      tx,
		createPartStmt:          q.createPartStmt,
		createUploadStmt:        q.createUploadStmt,
		findUploadByIdStmt:      q.findUploadByIdStmt,
		findUploadPartsByIdStmt: q.findUploadPartsByIdStmt,
	}
}
