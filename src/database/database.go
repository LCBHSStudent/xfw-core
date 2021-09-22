package database

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	util "github.com/LCBHSStudent/xfw-core/util"
	_ "github.com/go-sql-driver/mysql"
)

const MySqlConnKey = "mysql-connect"

var db *sql.DB

func init() {
	var err error

	mysqlConf := util.GetObjectByKey(MySqlConnKey).(map[interface{}]interface{})

	connectStr := strings.Join([]string{
		mysqlConf["username"].(string),
		":",
		mysqlConf["password"].(string),
		"@tcp(",
			mysqlConf["ipv4"].(string), ":", mysqlConf["port"].(string),
		")/", 
		mysqlConf["db-name"].(string),
	}, "")

	db, err = sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect db\n\t", err)
	} else {
		log.Printf("Succeed to connect db")
	}


	go func() {
		sigs := make(chan os.Signal, 1)

		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	
		<-sigs
        err := db.Close()
		if err != nil {
			log.Printf("%v", err)
		}

		signal.Stop(sigs)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	// defer db.Close()
}
