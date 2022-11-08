package infrastructure

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/service/interfaces"
)

type SqliteHandler struct {
	db *sql.DB
}

func (h *SqliteHandler) Close() error {
	return h.db.Close()
}

func (h *SqliteHandler) Execute(statement string, args ...interface{}) (sql.Result, error) {
	return h.db.Exec(statement, args...)
}

func (h *SqliteHandler) Query(statement string, args ...interface{}) (interfaces.Rows, error) {
	rows, err := h.db.Query(statement, args...)
	if err != nil {
		return nil, fmt.Errorf("db.query [%s]: %w", statement, err)
	}

	row := new(SqliteRow)
	row.Rows = rows
	return row, nil
}

type SqliteRow struct {
	Rows *sql.Rows
}

func (r SqliteRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}

func (r SqliteRow) Close() error {
	return r.Rows.Close()
}

func NewSqliteHandler(dbfileName string) (*SqliteHandler, error) {

	db, err := sql.Open("sqlite3", dbfileName)
	if err != nil {
		return nil, fmt.Errorf("sql.Open sqlite3 [%s] err: %w", dbfileName, err)
	}

	sqliteHandler := new(SqliteHandler)
	sqliteHandler.db = db

	return sqliteHandler, nil

}
