package src

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func getMysqlDB(dsn string) (db *sql.DB) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Sprintf("open mysql err: %+v", err)
	}
	return db
}

func getTable(db *sql.DB) []byte {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE %s", *tableName))
	if err != nil {
		fmt.Sprintf("get table err: %+v\n", err)
	}
	var (
		name   string
		sqlStr string
	)
	for rows.Next() {
		rows.Scan(&name, &sqlStr)
	}
	return []byte(sqlStr)
}
