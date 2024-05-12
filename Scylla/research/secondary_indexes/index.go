package secondary_indexes

import (
	"github.com/gocql/gocql"
	"main/research/secondary_indexes/sql_commands"
	"main/scripts/sql_common"
)

func Run(session *gocql.Session, keySpaceName, tableName, indexName string, needClear bool) error {
	if needClear {
		defer func() {
			_ = sql_commands.DropSecondaryIndex(session, indexName)

			_ = sql_common.DropTable(session, keySpaceName+"."+tableName)

			_ = sql_common.DropKeySpace(session, keySpaceName)
		}()
	}

	if err := sql_common.CreateKeySpace(session, keySpaceName); err != nil {
		return err
	}

	if err := sql_common.CreateTable(session, keySpaceName+"."+tableName); err != nil {
		return err
	}

	if err := sql_common.InsertRandomData(session, keySpaceName+"."+tableName, 1); err != nil {
		return err
	}

	if err := sql_commands.CreateSecondaryIndex(session, keySpaceName+"."+tableName, indexName, "field1"); err != nil {
		return err
	}

	return nil
}
