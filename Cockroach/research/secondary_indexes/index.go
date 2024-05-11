package secondary_indexes

import (
	"context"
	"fmt"
	"log"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"

	"main/scripts/pebble_engine"
	"main/scripts/sql_commands"
)

func SecondaryIndex(ctx context.Context, conn *pgx.Conn, tableName, secretKey string, needDelete, showAll bool) error {
	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.CreateTable(context.Background(), tx, tableName)
	}); err != nil {
		return err
	}

	if err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.InsertRow(ctx, tx, tableName, 2281337, secretKey)
	}); err != nil {
		return err
	}

	log.Print("Before creating a secondary index")
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

	if err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.CreateSecondaryIndex(context.Background(), tx, tableName, "research_index_1", "field2")
	}); err != nil {
		return err
	}

	log.Print("After creating a secondary index")
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

	if needDelete {
		if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return sql_commands.DeleteTable(context.Background(), tx, tableName)
		}); err != nil {
			return err
		}
	}

	return nil
}
