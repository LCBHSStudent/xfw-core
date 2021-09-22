package database

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	sqlUserName = "root"
	sqlPassword = "password"
	sqlAddr     = "127.0.0.1"
	sqlPort     = "3306"
	sqlDbName   = "xfw"
)

var db *sql.DB

func init() {
	var err error
	connectStr := strings.Join([]string{
		sqlUserName, ":", sqlPassword,
		"@tcp(", sqlAddr, ":", sqlPort, ")/", sqlDbName,
	}, "")

	db, err = sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
