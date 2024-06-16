package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Country7/backend-captaincode-mysql/util"

	"github.com/Country7/backend-captaincode-mysql/api"
	db "github.com/Country7/backend-captaincode-mysql/db/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	fmt.Println(">> main Loaded Config: ", config) // Выводим загруженную конфигурацию

	// Преобразуем URL-подобную строку подключения в DSN формат только для mysql
	dsnDBSource, err := parseDSN(config.DBSource)
	if err != nil {
		log.Fatal(">> main invalid DB source:", err)
	}
	fmt.Println(">> main DSN DB Connection Source: ", dsnDBSource) // Выводим строку подключения

	conn, err := sql.Open(config.DBDriver, dsnDBSource)
	if err != nil {
		log.Fatal("cannot connect to the db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

// parseDSN преобразует URL-подобную строку подключения в формат DSN
func parseDSN(dbSource string) (string, error) {
	// fmt.Println(">> main.parseDSN dbSource: ", dbSource) // отладка
	if !strings.HasPrefix(dbSource, "mysql://") {
		return "", fmt.Errorf("invalid dbSource prefix")
	}
	// Убираем "mysql://" префикс
	dsn := strings.TrimPrefix(dbSource, "mysql://")
	// fmt.Println(">> main.parseDSN dsn: ", dsn) // отладка
	return dsn, nil
}
