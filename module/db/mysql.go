package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Mysql *sql.DB

func InitDB() {
	var err error
	Mysql, err = sql.Open("mysql", "poker:fake@tcp(172.17.0.2:3306)/game_poker")
	if err != nil {
		panic(err)
	}
	Mysql.SetConnMaxLifetime(time.Minute * 3)
	Mysql.SetMaxOpenConns(10)
	Mysql.SetMaxIdleConns(10)
}

func SelectQueryRows(query string) (*sql.Rows, error) {
	/*mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return nil, err
	}
	defer mysql.Close()*/

	rows, err := Mysql.Query(query)
	if err != nil {
		log.Println("Failed to execute : ", err)
		return nil, err
	}
	return rows, err
}

func SelectQueryRow(query string) (*sql.Row, error) {
	/*mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return nil, err
	}
	defer mysql.Close()*/

	row := Mysql.QueryRow(query)
	return row, nil
}

func ExecuteQuery(query string) (int64, error) {
	/*mysql, err := sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
	if err != nil {
		log.Println("Failed to open db : ", err)
		return 0, err
	}
	defer mysql.Close()*/

	result, err := Mysql.Exec(query)
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

func CheckPing() {
	var err error
	if Mysql == nil {
		Mysql, err = sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
		if err != nil {
			panic(err)
		}
	}
	err = Mysql.Ping()
	if err != nil {
		Mysql, err = sql.Open("mysql", "poker:fake@tcp(localhost:3306)/game_poker")
		if err != nil {
			panic(err)
		}
	}

}
