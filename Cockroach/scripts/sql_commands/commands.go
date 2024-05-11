package sql_commands

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"main/scripts"
	"math/rand"
)

func CreateTable(ctx context.Context, tx pgx.Tx, tableName string) error {
	const sqlCommand = "CREATE TABLE IF NOT EXISTS ? (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), field1 INT, field2 STRING)"

	log.Println("Creating the table...")
	if _, err := tx.Exec(ctx, strings.Replace(sqlCommand, "?", tableName, 1)); err != nil {
		return err
	}
	return nil
}

func DeleteTable(ctx context.Context, tx pgx.Tx, tableName string) error {
	const sqlCommand = "DROP TABLE ?"

	log.Print("Deleting the table...")
	if _, err := tx.Exec(ctx, strings.Replace(sqlCommand, "?", tableName, 1)); err != nil {
		return err
	}
	return nil
}

func InsertRow(ctx context.Context, tx pgx.Tx, tableName string, field1 int, field2 string) error {
	const sqlCommand = "INSERT INTO ? (field1, field2) VALUES ($1, $2)"

	log.Println("Creating new row...")
	if _, err := tx.Exec(ctx, strings.Replace(sqlCommand, "?", tableName, 1), field1, field2); err != nil {
		return err
	}
	return nil
}

func DeleteRow(ctx context.Context, tx pgx.Tx, tableName string, id uuid.UUID) error {
	const sqlCommand = "DELETE FROM ? WHERE id IN ($1)"

	log.Printf("Deleting rows with IDs %s...", id)
	if _, err := tx.Exec(ctx, strings.Replace(sqlCommand, "?", tableName, 1), id); err != nil {
		return err
	}
	return nil
}

func CreateSecondaryIndex(ctx context.Context, tx pgx.Tx, tableName, indexName, fieldName string) error {
	const sqlCommand = "CREATE INDEX ? ON ? (?);"

	log.Printf("Creating secondary index...")
	if _, err := tx.Exec(ctx, strings.Replace(strings.Replace(strings.Replace(sqlCommand, "?", indexName, 1), "?", tableName, 1), "?", fieldName, 1)); err != nil {
		return err
	}
	return nil
}

func InsertRandomData(ctx context.Context, tx pgx.Tx, tableName string, count int) error {
	const stringLen = 16
	for i := 0; i < count; i++ {
		if err := InsertRow(ctx, tx, tableName, rand.Int(), scripts.RandString(stringLen)); err != nil {
			return err
		}
	}
	return nil
}

func ArtificialSelect(ctx context.Context, tx pgx.Tx, tableName string) error {
	const sqlCommand = "SELECT field1, field2 FROM ? WHERE field1 > 0;"

	results, err := tx.Query(ctx, strings.Replace(sqlCommand, "?", tableName, 1))
	if err != nil {
		return err
	}

	for results.Next() {
		var field1 int
		var field2 string

		err = results.Scan(&field1, &field2)
		if err != nil {
			return err
		}
	}

	return nil
}
