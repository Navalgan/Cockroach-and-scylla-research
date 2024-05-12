package realization

import (
	"main/scripts/sql_common"

	"context"
	"log"
	"sync"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

func InsertLoad(ctx context.Context, connStr, tableName string, connCount, loadTime int, wg *sync.WaitGroup) {
	log.Print("starting insert load...")

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
					return sql_common.InsertRandomData(context.Background(), tx, tableName, 10)
				}); err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
}

func SelectLoad(ctx context.Context, connStr, tableName string, connCount, loadTime int, wg *sync.WaitGroup) {
	log.Print("starting select load...")

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
					return sql_common.ArtificialSelect(context.Background(), tx, tableName)
				}); err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
}

func Load(ctx context.Context, conn *pgx.Conn, connStr, tableName string, insertCount, selectCount, loadTime int, needDelete bool) error {
	if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return sql_common.CreateTable(context.Background(), tx, tableName)
	}); err != nil {
		return err
	}

	// TODO
	// snapshot from gdb

	wg := sync.WaitGroup{}

	InsertLoad(ctx, connStr, tableName, insertCount, loadTime, &wg)
	SelectLoad(ctx, connStr, tableName, selectCount, loadTime, &wg)

	wg.Wait()

	if needDelete {
		if err := crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return sql_common.DropTable(context.Background(), tx, tableName)
		}); err != nil {
			return err
		}
	}

	return nil
}
