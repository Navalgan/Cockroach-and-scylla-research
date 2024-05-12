package views

import (
	"main/research/views/sql_commands"
	"main/scripts/sql_common"

	"context"
	"fmt"
	"strconv"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

func Views(ctx context.Context, conn *pgx.Conn, tableName, viewName string, needClear bool) error {
	if needClear {
		defer func() {
			_ = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
				return sql_commands.DropView(ctx, tx, viewName)
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

	const dataCount = 10
	for i := 0; i < dataCount; i++ {
		if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return sql_common.InsertRow(ctx, tx, tableName, i, tableName+strconv.Itoa(i))
		}); err != nil {
			return err
		}
	}

	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.CreateView(ctx, tx, viewName, tableName, "field1")
	}); err != nil {
		return err
	}

	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		results, err := sql_common.SelectAll(ctx, tx, viewName)
		defer results.Close()

		i := 0
		for results.Next() {
			var field string
			err := results.Scan(&field)
			if err != nil {
				return err
			}
			fmt.Printf("%d row -- %s\n", i, field)
			i++
		}

		return err
	}); err != nil {
		return err
	}

	fmt.Println("insert extra row")
	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_common.InsertRow(ctx, tx, tableName, 2281337, "extra row")
	}); err != nil {
		return err
	}

	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		results, err := sql_common.SelectAll(ctx, tx, viewName)
		defer results.Close()

		i := 0
		for results.Next() {
			var field string
			err := results.Scan(&field)
			if err != nil {
				return err
			}
			fmt.Printf("%d row -- %s\n", i, field)
			i++
		}

		return err
	}); err != nil {
		return err
	}

	return nil

}
