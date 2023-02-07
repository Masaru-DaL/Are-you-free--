package db

import (
	"database/sql"
	"log"
	"src/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDatabase() *sql.DB {
	config.LoadConfigForYaml()

	db, err := sql.Open(config.Config.DB.SQLDriver, config.Config.DB.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// エラーは出なかったが、何らかの理由でデータベース接続ができなかった場合
	if db == nil {
		log.Fatal(err)
	}

	/* 接続が可能であることを確認する */
	err = db.Ping()
	if err != nil {
		defer db.Close()
	}

	return db
}
