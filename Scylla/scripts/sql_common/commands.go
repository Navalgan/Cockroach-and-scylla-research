package sql_common

import (
	"log"
	"main/scripts"
	"strings"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func CreateKeySpace(session *gocql.Session, keySpaceName string) error {
	const cqlCommand = "CREATE KEYSPACE IF NOT EXISTS ? WITH REPLICATION = {'class': 'NetworkTopologyStrategy'}"

	log.Println("Creating the keyspace...")

	return session.Query(strings.Replace(cqlCommand, "?", keySpaceName, 1)).Exec()
}

func DropKeySpace(session *gocql.Session, keySpaceName string) error {
	const cqlCommand = "DROP KEYSPACE ?"

	log.Println("Dropping the keyspace...")

	return session.Query(strings.Replace(cqlCommand, "?", keySpaceName, 1)).Exec()
}

func CreateTable(session *gocql.Session, tableName string) error {
	const cqlCommand = "CREATE TABLE IF NOT EXISTS ? (myID UUID, field1 INT, field2 TEXT, PRIMARY KEY (myID))"

	log.Println("Creating the table...")

	return session.Query(strings.Replace(cqlCommand, "?", tableName, 1)).Exec()
}

func DropTable(session *gocql.Session, tableName string) error {
	const cqlCommand = "DROP TABLE ?"

	log.Print("Dropping the table...")

	return session.Query(strings.Replace(cqlCommand, "?", tableName, 1)).Exec()
}

func InsertRow(session *gocql.Session, tableName, id string, field1 int, field2 string) error {
	const cqlCommand = "INSERT INTO ? (myID, field1, field2) VALUES (?, ?, ?)"

	log.Println("Creating new row...")

	return session.Query(strings.Replace(cqlCommand, "?", tableName, 1), id, "1", field2).Exec()
}

func DeleteRow(session *gocql.Session, tableName string, id string) error {
	const sqlCommand = "DELETE FROM ? WHERE myID IN ($1)"

	log.Printf("Deleting rows with IDs %s...", id)

	return nil
}

func InsertRandomData(session *gocql.Session, tableName string, count int) error {
	const stringLen = 16

	for i := 0; i < count; i++ {
		if err := InsertRow(session, tableName, uuid.NewString(), i, scripts.RandString(stringLen)); err != nil {
			return err
		}
	}

	return nil
}

func SelectAll(session *gocql.Session, tableName string) error {
	const sqlCommand = "SELECT * FROM ?"

	return nil
}
