package secondary_indexes

import (
	"main/research/secondary_indexes/sql_commands"
	"main/scripts"
	"main/scripts/pebble_engine"
	"main/scripts/sql_common"

	"context"
	"fmt"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

func SecondaryIndex(ctx context.Context, conn *pgx.Conn, tableName, indexName, secretKey string, needClear, showAll bool) error {
	if needClear {
		defer func() {
			_ = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
				return sql_commands.DropSecondaryIndex(ctx, tx, indexName)
			})

			_ = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
				return sql_common.DropTable(ctx, tx, tableName)
			})
		}()
	}

	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_common.CreateTable(ctx, tx, tableName)
	}); err != nil {
		return err
	}

	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_common.InsertRow(ctx, tx, tableName, 2281337, secretKey)
	}); err != nil {
		return err
	}

	fmt.Print("Before creating a secondary index\n")
	nodesBefore := pebble_engine.CheckAllNodes([]byte(secretKey), showAll)
	for len(nodesBefore) == 0 {
		nodesBefore = pebble_engine.CheckAllNodes([]byte(secretKey), showAll)
	}

	nodesWithKey := make([]string, 0)
	for node, nodeRes := range nodesBefore {
		nodesWithKey = append(nodesWithKey, node)

		fmt.Printf("On %s:\n", node)
		for _, kv := range nodeRes {
			fmt.Printf("Finded key: %s\nValue: %s\n", kv.Key, kv.Value)
		}
	}

	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.CreateSecondaryIndex(ctx, tx, tableName, indexName, "field2")
	}); err != nil {
		return err
	}

	fmt.Print("After creating a secondary index\n")
	nodesAfter := pebble_engine.CheckNodes(nodesWithKey, []byte(secretKey), true)

	for node, nodeRes := range nodesAfter {
		fmt.Printf("On %s:\n", node)
		for _, kv := range nodeRes {
			fmt.Printf("Finded key: %s\nValue: %s\n", kv.Key, kv.Value)
		}
	}

	if len(nodesBefore) != 0 {
		fmt.Printf("Key was before on:\n")
		for node, nodeRes := range nodesBefore {
			if len(nodeRes) != 0 {
				fmt.Printf("%s ", node)
			}
		}
		fmt.Printf("\n")
	}

	if len(nodesAfter) != 0 {
		fmt.Printf("Key was after on:\n")
		for node, nodeRes := range nodesAfter {
			if len(nodeRes) != 0 {
				fmt.Printf("%s ", node)
			}
		}
		fmt.Printf("\n")
	}

	fmt.Print("Data in SST files:\n")
	for _, node := range nodesWithKey {
		files := scripts.GetSSTFromNode(node)
		for _, file := range files {
			data, err := scripts.GetDataFromSST(file, secretKey)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
		}
	}

	return nil
}
