package mydb

import (
	"database/sql"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func (db *DB) AddUser(username string, password string) error {
	_, err := db.Exec("insert into go_test_task.users (name, password, token) values (?, ?, '')", username, password)
	return err
}

func (db *DB) DeleteUser(username string) error {
	_, err := db.Exec("delete from go_test_task.users where name = ?", username)
	return err
}

func (db *DB) CheckPassword(username, password string) bool {
	row, err := db.Query("select password from go_test_task.users where name = ? AND password = ?", username, password)
	if err != err {
		return false
	}

	return row.Next()
}

func (db *DB) SetToken(username, token string) error {
	_, err := db.Exec("update go_test_task.users set token = ? where name = ?", token, username)
	return err
}

func (db *DB) GetToken(username string) (string, error) {
	row, err := db.Query("select token from go_test_task.users where name = ?", username)
	if err != err {
		return "", err
	}
	token := ""
	row.Next()
	err = row.Scan(&token)
	return token, err
}

var (
	database *DB
	once     sync.Once
)

func Init() (*DB, error) {
	var err error
	once.Do(func() {
		var db *sql.DB
		db, err = sql.Open("mysql", "go_test_task_user:password@/go_test_task")
		_, err = db.Exec(`CREATE TABLE go_test_task.users(name VARCHAR(50) PRIMARY KEY, 
		password VARCHAR(20), token varchar(100))`)
		database = &DB{db}
	})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, err
	}
	return database, nil
}
