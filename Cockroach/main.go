package main

import (
	"main/research/realization"
	"main/research/secondary_indexes"
	"main/research/views"

	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/pflag"
)

const connStr = "postgresql://root@localhost:26250/defaultdb?sslmode=disable"

var (
	flagTableName       = pflag.String("table-name", "research", "")
	flagIndexName       = pflag.String("index-name", "researchIndex", "")
	flagViewName        = pflag.String("view-name", "researchView", "")
	flagScriptName      = pflag.String("script-name", "secondary_indexes", "")
	flagSecretKey       = pflag.String("secret-key", "my secret key 1337", "")
	flagIsNeedClear     = pflag.Bool("is-need-clear", true, "")
	flagIsShowAll       = pflag.Bool("is-show-all", false, "")
	flagLoadInsertCount = pflag.Int("load-insert-count", 0, "")
	flagLoadSelectCount = pflag.Int("load-select-count", 0, "")
	flagLoadTime        = pflag.Int("load-time", 100, "")
)

func main() {
	pflag.Parse()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	switch *flagScriptName {
	case "load":
		err = realization.Load(ctx, conn, connStr, *flagTableName, *flagLoadInsertCount, *flagLoadSelectCount, *flagLoadTime, *flagIsNeedClear)
	case "secondary_indexes":
		err = secondary_indexes.SecondaryIndex(ctx, conn, *flagTableName, *flagIndexName, *flagSecretKey, *flagIsNeedClear, *flagIsShowAll)
	case "views":
		err = views.Views(ctx, conn, *flagTableName, *flagViewName, *flagIsNeedClear)
	}

	if err != nil {
		log.Fatal(err)
	}
}
