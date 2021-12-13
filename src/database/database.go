package database

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	util "github.com/LCBHSStudent/xfw-core/util"
	_ "github.com/go-sql-driver/mysql"
)

const MySqlConnKey = "mysql-connect"

var db *sql.DB

// sync.Mutex

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

func checkTableExist(table string) bool {
	query := "SELECT nsp FROM " + table
	_, err := db.Query(query)
	if err != nil &&
		err.Error() == "Error 1054: Unknown column 'nsp' in 'field list'" {
		return true
	} else {
		return false
	}
}

func CreateTableIfNotExist(table string, tableFormat string) {
	if checkTableExist(table) {
		return
	}

	create := `CREATE TABLE IF NOT EXISTS ` + table + "(" + tableFormat + ");"
	_, err := db.Exec(create)
	if err != nil {
		panic(err)
	}
}

func GetTableRowCount(table string) int64 {
	totalRow, err := db.Query("SELECT COUNT(*) FROM " + table)
	if err != nil {
		log.Println(err)
		return -1
	}

	total := int64(-1)
	for totalRow.Next() {
		err := totalRow.Scan(
			&total,
		)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return total
}

func ExecWriteSql(table string, fields string, args []interface{}) {
	valuesStr := "VALUES("

	for i := 0; i < len(args); i++ {
		valuesStr += "?"
		if i != len(args)-1 {
			valuesStr += ","
		} else {
			valuesStr += ")"
		}
	}

	stmt, err := db.Prepare("INSERT INTO " + table + "(" + fields + ")" + valuesStr)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Println(err)
		return
	}
}

func ExecReadSql(table string, fields string, cond string, args []interface{}) (list []map[string]interface{}, err error) {
	sqlStr := "SELECT " + fields + " FROM " + table

	if len(cond) != 0 {
		sqlStr += " WHERE " + cond
	}

	rows, err := db.Query(sqlStr, args...)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err)
			return
		}

		ret := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				ret[columns[i]] = nil
			} else {
				switch val := (*scanArgs[i].(*interface{})).(type) {
				case byte:
					ret[columns[i]] = val
					break
				case []byte:
					v := string(val)
					switch v {
					case "\x00":
						ret[columns[i]] = 0
					case "\x01":
						ret[columns[i]] = 1
					default:
						ret[columns[i]] = v
						break
					}
					break
				case time.Time:
					if val.IsZero() {
						ret[columns[i]] = nil
					} else {
						ret[columns[i]] = val.Format("2006-01-02 15:04:05")
					}
					break
				default:
					ret[columns[i]] = val
				}
			}
		}
		list = append(list, ret)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
