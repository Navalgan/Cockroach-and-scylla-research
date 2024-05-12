package main

import (
	"github.com/gocql/gocql"
	"github.com/spf13/pflag"
	"log"
	"main/research/secondary_indexes"
	"time"
)

const connStr = "postgresql://root@localhost:26250/defaultdb?sslmode=disable"

var (
	flagTableName    = pflag.String("table-name", "ResearchTable1", "")
	flagKeySpaceName = pflag.String("key-space-name", "ResearchKeySpace1", "")
	flagIndexName    = pflag.String("index-name", "Field1Index1", "")
	flagScriptName   = pflag.String("script-name", "secondary_indexes", "")
)

func main() {
	pflag.Parse()

	cluster := gocql.NewCluster("localhost")

	retryPolicy := &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        10 * time.Second,
		NumRetries: 5,
	}
	cluster.Timeout = 5 * time.Second
	cluster.RetryPolicy = retryPolicy
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("unable to connect to scylla with error: %s", err.Error())
	}
	defer session.Close()

	switch *flagScriptName {
	case "secondary_indexes":
		err = secondary_indexes.Run(session, *flagKeySpaceName, *flagTableName, *flagIndexName, false)
	}

	if err != nil {
		log.Fatal(err)
	}
}
