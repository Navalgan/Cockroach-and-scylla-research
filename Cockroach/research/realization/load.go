package realization

import (
	"context"
	"log"
	"sync"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"

	"main/scripts/sql_commands"
)

func SelectLoad(ctx context.Context, conn *pgx.Conn, connStr, tableName string, connCount, loadTime int, needDelete bool) error {
	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.CreateTable(context.Background(), tx, tableName)
	}); err != nil {
		return err
	}

	log.Print("starting select load...")

	wg := sync.WaitGroup{}

	wg.Add(connCount)
	for i := 0; i < connCount; i++ {
		go func() {
			conn, err := pgx.Connect(ctx, connStr)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close(context.Background())

			for j := 0; j < loadTime; j++ {
				if err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
					return sql_commands.ArtificialSelect(context.Background(), tx, tableName)
				}); err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if needDelete {
		if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return sql_commands.DeleteTable(context.Background(), tx, tableName)
		}); err != nil {
			return err
		}
	}

	return nil
}

func InsertLoad(ctx context.Context, conn *pgx.Conn, connStr, tableName string, connCount, loadTime int, needDelete bool) error {
	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_commands.CreateTable(context.Background(), tx, tableName)
	}); err != nil {
		return err
	}

	log.Print("starting insert load...")

	wg := sync.WaitGroup{}

	wg.Add(connCount)
	for i := 0; i < connCount; i++ {
		go func() {
			conn, err := pgx.Connect(ctx, connStr)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close(context.Background())

			for j := 0; j < loadTime; j++ {
				if err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
					return sql_commands.InsertRandomData(context.Background(), tx, tableName, 10)
				}); err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if needDelete {
		if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return sql_commands.DeleteTable(context.Background(), tx, tableName)
		}); err != nil {
			return err
		}
	}

	return nil
}
