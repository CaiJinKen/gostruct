package source

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DbSource struct {
	DSN       string
	TableName string
}

func (d *DbSource) GetData() (data []byte, err error) {
	if d.DSN == "" {
		return
	}

	db := d.getMysqlDB(d.DSN)
	if db == nil {
		fmt.Printf("get db connection err.\n")
		return
	}
	defer db.Close()
	data = d.getTable(db)

	return
}

func (d *DbSource) getMysqlDB(dsn string) (db *sql.DB) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("open mysql err: %+v\n", err)
	}
	return db
}

func (d *DbSource) getTable(db *sql.DB) []byte {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE %s", d.TableName))
	if err != nil {
		fmt.Printf("get table err: %+v\n", err)
		return nil
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
