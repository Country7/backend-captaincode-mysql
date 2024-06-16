package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Country7/backend-captaincode-mysql/util"

	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	fmt.Println("Initializing test suite")

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSourceTestsqlc)
	if err != nil {
		log.Fatal("cannot connect to the db:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
