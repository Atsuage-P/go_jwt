package config

import (
	"database/sql"
	"log"
	"time"

	"go_oauth/env"

	"github.com/go-sql-driver/mysql"
)

func ConnectDB(cnf *env.DBConfig) *sql.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("DB connect Error: %v", err)
	}
	c := mysql.Config{
		DBName:    cnf.Name,
		User:      cnf.User,
		Passwd:    cnf.Password,
		Addr:      cnf.Host + ":" + cnf.Port,
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}
	db, err := sql.Open(cnf.MS, c.FormatDSN())
	if err != nil {
		log.Fatalf("DB connect Error: %v", err)
	}
	db.SetConnMaxLifetime(time.Duration(cnf.MaxLifeTimeMin))
	db.SetMaxOpenConns(cnf.MaxOpenConns)
	db.SetMaxIdleConns(cnf.MaxOpenIdleConns)
	return db
}
