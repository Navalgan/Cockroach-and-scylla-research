package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/pflag"

	"main/research/realization"
	"main/research/secondary_indexes"
)

const connStr = "postgresql://root@localhost:26250/defaultdb?sslmode=disable"

var (
	flagTableName         = pflag.String("table-name", "research", "")
	flagScriptName        = pflag.String("script-name", "secondary_indexes", "")
	flagSecretKey         = pflag.String("secret-key", "my secret key 1337", "")
	flagIsNeedDeleteTable = pflag.Bool("is-need-delete-table", true, "")
	flagIsShowAll         = pflag.Bool("is-show-all", false, "")
	flagLoadConnCount     = pflag.Int("load-conn-count", 3, "")
	flagLoadTime          = pflag.Int("load-conn-count", 100, "")
)

func main() {
	pflag.Parse()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	switch *flagScriptName {
	case "alter_range":
	case "insert-load":
		err = realization.InsertLoad(ctx, conn, connStr, *flagTableName, *flagLoadConnCount, *flagLoadTime, *flagIsNeedDeleteTable)
	case "select-load":
		err = realization.InsertLoad(ctx, conn, connStr, *flagTableName, *flagLoadConnCount, *flagLoadTime, *flagIsNeedDeleteTable)
	case "secondary_indexes":
		err = secondary_indexes.SecondaryIndex(ctx, conn, *flagTableName, *flagSecretKey, *flagIsNeedDeleteTable, *flagIsShowAll)
	case "views":
	}

	if err != nil {
		log.Fatal(err)
	}
}
