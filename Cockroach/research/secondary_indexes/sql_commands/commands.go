package sql_commands

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

func CreateSecondaryIndex(ctx context.Context, tx pgx.Tx, tableName, indexName, fieldName string) error {
	const sqlCommand = "CREATE INDEX ? ON ? (?);"

	log.Printf("Creating secondary index...")
	if _, err := tx.Exec(ctx, strings.Replace(strings.Replace(strings.Replace(sqlCommand, "?", indexName, 1), "?", tableName, 1), "?", fieldName, 1)); err != nil {
		return err
	}
	return nil
}

func DropSecondaryIndex(ctx context.Context, tx pgx.Tx, indexName string) error {
	const sqlCommand = "DROP INDEX ? CASCADE"

	log.Printf("Dropping secondary index...")
	if _, err := tx.Exec(ctx, strings.Replace(sqlCommand, "?", indexName, 1)); err != nil {
		return err
	}

	return nil
}
