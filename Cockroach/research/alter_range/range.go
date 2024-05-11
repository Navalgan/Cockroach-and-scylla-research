package alter_range

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func SetNumReplicasForRange(ctx context.Context, tx pgx.Tx, rangeName string, numReplicas int) error {
	const sqlCommand = "ALTER RANGE ? CONFIGURE ZONE USING num_replicas = ?"

	log.Print("getting id of cluster nodes...")
	_, err := tx.Exec(ctx, sqlCommand, rangeName, numReplicas)
	if err != nil {
		return err
	}

	return nil
}

func GetNodeIDs(ctx context.Context, tx pgx.Tx) ([]int, error) {
	const sqlCommand = "SELECT store_id FROM crdb_internal.kv_store_status"

	log.Print("getting id of cluster nodes...")
	result, err := tx.Query(ctx, sqlCommand)
	if err != nil {
		return nil, err
	}

	res := make([]int, 0)
	for result.Next() {
		var id int
		result.Scan(&id)
		res = append(res, id)
	}

	return res, nil
}

func GetRangeInfo(ctx context.Context, tx pgx.Tx, tableName string) ([]int, error) {
	const sqlCommand = "WITH user_info AS (SHOW RANGES FROM TABLE ?) SELECT range_id, lease_holder, lease_holder_locality FROM user_info"

	log.Print("getting id of cluster nodes...")
	result, err := tx.Query(ctx, sqlCommand)
	if err != nil {
		return nil, err
	}

	res := make([]int, 0)
	for result.Next() {
		var id int
		result.Scan(&id)
		res = append(res, id)
	}

	return res, nil
}
