package sql_commands

import (
	"log"
	"strings"

	"github.com/gocql/gocql"
)

func CreateSecondaryIndex(session *gocql.Session, tableName, indexName, fieldName string) error {
	const cqlCommand = "CREATE INDEX ? ON ? (?)"

	log.Println("Creating the keyspace...")

	return session.Query(strings.Replace(strings.Replace(strings.Replace(cqlCommand, "?", indexName, 1), "?", tableName, 1), "?", fieldName, 1)).Exec()
}

//func CreateSecondaryIndex(ctx context.Context, tx pgx.Tx, tableName, indexName, fieldName string) error {
//	const sqlCommand = "CREATE INDEX ? ON ? (?);"
//
//	log.Printf("Creating secondary index...")
//	if _, err := tx.Exec(ctx, strings.Replace(strings.Replace(strings.Replace(sqlCommand, "?", indexName, 1), "?", tableName, 1), "?", fieldName, 1)); err != nil {
//		return err
//	}
//	return nil
//}

func DropSecondaryIndex(session *gocql.Session, indexName string) error {
	const cqlCommand = "DROP INDEX ?"

	return session.Query(strings.Replace(cqlCommand, "?", indexName, 1)).Exec()
}
