package dbutils

import (
	"log"
	"os"
)

func GetSQLDSN() string {
	deflink := "postgresql://root:root@auth_micro_service_db:5431/auth_micro_service?sslmode=disable"
	// read ENV variable SQL_DSN
	sqldsnenv := os.Getenv("SQL_DSN")
	if sqldsnenv != "" {
		deflink = sqldsnenv
	}
	log.Default().Printf("returning sql dsn %s", deflink)
	return deflink
}
