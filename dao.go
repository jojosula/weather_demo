package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/mattn/go-sqlite3/driver"
)

const (
	DB_DRIVER        = "sqlite3"
	DB_PATH          = "./demo.db"
	DEFAULT_USERNAME = "admin"
	DEFAULT_PASSWORD = "admin"

	CREATE_USER_TABLE_SQL = `CREATE TABLE IF NOT EXISTS userinfo (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		username VARCHAR(50) NOT NULL UNIQUE, 
		password VARCHAR(50),
		created_at DATETIME default CURRENT_TIMESTAMP
	);`

	CREATE_TOKEN_TABLE_SQL = `CREATE TABLE IF NOT EXISTS usertoken (
		user_id INTEGER NOT NULL PRIMARY KEY, 
		access_token VARCHAR(50),
		created_at DATETIME default CURRENT_TIMESTAMP
	);`

	CREATE_WEATHER_TABLE_SQL = `CREATE TABLE IF NOT EXISTS weather (
		location VARCHAR(50) NOT NULL PRIMARY KEY, 
		data TEXT,
		created_at DATETIME default CURRENT_TIMESTAMP
	);`

	INSERT_USER_SQL = `INSERT INTO userinfo (username, password) 
		values ('%s', '%s')`

	INSERT_USER_TOKEN_SQL = `INSERT INTO usertoken (user_id, access_token) 
		values (%d, '%s')`

	INSERT_WEATHER_SQL = `INSERT INTO weather (location, data) 
		values ('%s', '%s')`

	COUNT_ALL_USER_SQL    = `SELECT COUNT(*) as count FROM userinfo`
	QUERY_USER_SQL        = `SELECT * FROM userinfo WHERE username = '%s'`
	QUERY_USER_TOKEN_SQL  = `SELECT * FROM usertoken WHERE user_id = %d`
	QUERY_TOKEN_SQL       = `SELECT * FROM usertoken WHERE access_token = '%s'`
	QUERY_WEATHER_SQL     = `SELECT * FROM weather WHERE location = '%s'`
	UPDATE_WEATHER_SQL    = `UPDATE weather SET data = '%s' WHERE location = '%s'`
	DELETE_USER_TOKEN_SQL = `DELETE FROM usertoken WHERE access_token = '%s'`
)

var db *sql.DB
var once sync.Once

func getDBConnection() *sql.DB {
	var err error
	// make sure one instance
	once.Do(func() {
		db, err = sql.Open(DB_DRIVER, DB_PATH)
		checkErr(err)
	})
	return db
}

func DBInit() {
	// open db
	db = getDBConnection()
	db.SetMaxOpenConns(1)

	// create table
	stmt, _ := db.Prepare(CREATE_USER_TABLE_SQL)
	stmt.Exec()
	stmt, _ = db.Prepare(CREATE_TOKEN_TABLE_SQL)
	stmt.Exec()
	stmt, _ = db.Prepare(CREATE_WEATHER_TABLE_SQL)
	stmt.Exec()

	// insert default user
	insertSQL := fmt.Sprintf(INSERT_USER_SQL, DEFAULT_USERNAME, GetMD5(DEFAULT_PASSWORD))
	stmt, _ = db.Prepare(insertSQL)
	stmt.Exec()
}

func DBClose() {
	if db != nil {
		db.Close()
	}
}

func CreateUser(username, password string) {
	// insert new user
	insertSQL := fmt.Sprintf(INSERT_USER_SQL, username, GetMD5(password))
	stmt, _ := db.Prepare(insertSQL)
	stmt.Exec()
}

func GetUser(username, password string) *UserInfo {
	// query
	user := UserInfo{}
	querySQL := fmt.Sprintf(QUERY_USER_SQL, username)
	rows, _ := db.Query(querySQL)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
		// password incorrect
		if user.Password != GetMD5(password) {
			return &UserInfo{}
		}
	}
	return &user
}

func CreateUserToken(userID int) *UserToken {
	// generate new token
	userToken := UserToken{}
	userToken.AccessToken = TokenGenerate()
	userToken.UserID = userID
	userToken.CreatedAt = time.Now()

	// insert
	insertSQL := fmt.Sprintf(INSERT_USER_TOKEN_SQL, userID, userToken.AccessToken)
	stmt, _ := db.Prepare(insertSQL)
	stmt.Exec()
	return &userToken
}

func GetUserToken(userID int) *UserToken {
	// query
	userToken := UserToken{}
	querySQL := fmt.Sprintf(QUERY_USER_TOKEN_SQL, userID)
	rows, _ := db.Query(querySQL)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&userToken.UserID, &userToken.AccessToken, &userToken.CreatedAt)
	}
	return &userToken
}

func UserTokenIsValid(token string) bool {
	// query
	userToken := UserToken{}
	querySQL := fmt.Sprintf(QUERY_TOKEN_SQL, token)
	rows, _ := db.Query(querySQL)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&userToken.UserID, &userToken.AccessToken, &userToken.CreatedAt)
	}
	return (userToken.UserID != 0)
}

func DeleteUserToken(accessToken string) {
	// delete
	deleteSQL := fmt.Sprintf(DELETE_USER_TOKEN_SQL, accessToken)
	stmt, _ := db.Prepare(deleteSQL)
	stmt.Exec()
}

func CreateWeatherData(location, data string) {
	// insert weather data
	insertSQL := fmt.Sprintf(INSERT_WEATHER_SQL, location, data)
	stmt, _ := db.Prepare(insertSQL)
	stmt.Exec()
}

func UpdateWeatherData(location, data string) {
	// update weather data
	updateSQL := fmt.Sprintf(UPDATE_WEATHER_SQL, data, location)
	stmt, _ := db.Prepare(updateSQL)
	stmt.Exec()
}

func GetWeatherData(location string) *WeatherData {
	// query
	data := WeatherData{}
	querySQL := fmt.Sprintf(QUERY_WEATHER_SQL, location)
	rows, _ := db.Query(querySQL)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&data.Location, &data.Data, &data.CreatedAt)
	}
	return &data
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
