package randomGck

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"

	database "github.com/LCBHSStudent/xfw-core/src/database"
)

const (
	addressTable = "gck_address"
	descriptionTable = "gck_description"
)

func init() {
	gckTableStr := 
`
	ID          BIGINT           	NOT NULL    PRIMARY KEY   AUTO_INCREMENT,
	FROM_GROUP	BIGINT				NOT NULL,
	ADD_DATE	DATE				NOT NULL,
	DATA		BLOB				NOT NULL
`
	database.CreateTableIfNotExist(addressTable, gckTableStr)
	database.CreateTableIfNotExist(descriptionTable, gckTableStr)
}

func SaveAddress(groupID int64, msg string) {
	if len(msg) == 0 {
		return
	}

	database.ExecWriteSql(addressTable, "FROM_GROUP,ADD_DATE,DATA", []interface{} {
		groupID, time.Now().Format("2006-01-02 15:04:05"), msg,
	})
}

func SaveDescription(groupID int64, msg string) {
	if len(msg) == 0 {
		return
	}

	database.ExecWriteSql(descriptionTable, "FROM_GROUP,ADD_DATE,DATA", []interface{} {
		groupID, time.Now().Format("2006-01-02 15:04:05"), msg,
	})
}

func RemoveDescription(groupID int64, msg string) {
	if len(msg) == 0 {
		return
	}

}

func RemoveAddress(groupID, msg string) {
	if len(msg) == 0 {
		return
	}
	
}

func GenerateSpeech() string {
	var ret string
	
	addressCount := database.GetTableRowCount(addressTable)
	descriptionCount := database.GetTableRowCount(descriptionTable)

	if addressCount <= 0 {
		return ""
	}
	temp, err := rand.Int(rand.Reader, big.NewInt(addressCount))
	if err != nil {
		log.Fatal(err)
	}
	addressIdx := temp.Int64()
	result, err := database.ExecReadSql(addressTable, "DATA", "ID=?", []interface{}{addressIdx+1})
	if err != nil {
		log.Println(err)
		return ""
	}
	ret += result[0]["DATA"].(string)


	if descriptionCount <= 0 {
		return ""
	}
	temp, err = rand.Int(rand.Reader, big.NewInt(descriptionCount))
	if err != nil {
		log.Fatal(err)
	}
	descriptionIdx := temp.Int64()
	result, err = database.ExecReadSql(descriptionTable, "DATA", "ID=?", []interface{}{descriptionIdx+1})

	if err != nil {
		log.Println(err)
		return ""
	}
	ret += result[0]["DATA"].(string)


	return ret
}
