package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mysql `*sql.DB`

func initDB() {
	mysql, err := sql.Open("mysql", "root:apfhd0403@tcp(localhost:3306)/game_poker")
	if err != nil {
		panic(err)
	}
	defer mysql.Close()
	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)
}

func ExecuteSelectQuery(query string) (*sql.Rows, error) {
	mysql, err := sql.Open("mysql", "root:apfhd0403@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return nil, err
	}
	defer mysql.Close()

	rows, err := mysql.Query("SELECT * FROM t_users_gameinfo")
	if err != nil {
		log.Println("Failed to execute : ", err)
		return nil, err
	}
	defer rows.Close()
	return rows, err
}

func ExecuteUpdateQuery(query string) error {
	mysql, err := sql.Open("mysql", "root:apfhd0403@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return err
	}
	defer mysql.Close()

	result, err := mysql.Exec("UPDATE t_users_gameinfo SET id = 'test4' WHERE id = 'tes4'")
	if err != nil {
		log.Println("Failed to execute : ", err)
		return err
	}

	_, err = result.RowAffected()
	if err != nil {
		log.Println("Failed to execute : ", err)
		return err
	}
	return nil
}
