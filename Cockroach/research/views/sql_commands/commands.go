package sql_commands

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

func CreateView(ctx context.Context, tx pgx.Tx, viewName, tableName, fieldName string) error {
	const sqlCommand = "CREATE VIEW ? AS SELECT ? FROM ?"

	log.Printf("Creating view...")
	if _, err := tx.Exec(ctx, strings.Replace(strings.Replace(strings.Replace(sqlCommand, "?", viewName, 1), "?", fieldName, 1), "?", tableName, 1)); err != nil {
		return err
	}

	return nil
}

func DropView(ctx context.Context, tx pgx.Tx, viewName string) error {
	const sqlCommand = "DROP VIEW ? CASCADE"

	log.Printf("Dropping view...")
	if _, err := tx.Exec(ctx, strings.Replace(sqlCommand, "?", viewName, 1)); err != nil {
		return err
	}

	return nil
}
