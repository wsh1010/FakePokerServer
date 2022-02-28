package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//var mysql *sql.DB

func InitDB() {
	mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		panic(err)
	}
	defer mysql.Close()
	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)
}

func SelectQueryRows(query string) (*sql.Rows, error) {
	mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return nil, err
	}
	defer mysql.Close()

	rows, err := mysql.Query(query)
	if err != nil {
		log.Println("Failed to execute : ", err)
		return nil, err
	}
	return rows, err
}

func SelectQueryRow(query string) (*sql.Row, error) {
	mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return nil, err
	}
	defer mysql.Close()
	row := mysql.QueryRow(query)
	return row, nil
}

func ExecuteQuery(query string) (int64, error) {
	mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return 0, err
	}
	defer mysql.Close()

	result, err := mysql.Exec(query)
	if err != nil {
		log.Println("Failed to execute : ", err)
		return 0, err
	}
	success, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to execute : ", err)
		return 0, err
	}
	if success == 0 {
		log.Println("Not Excute : 0")
		return 0, err
	}

	return success, nil
}
