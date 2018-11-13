package util

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	connPool        *sql.DB // database connection pool
	maxOpenConn     = 100
	maxIdleConn     = 30
	maxConnLifeTime = time.Minute * 10
)

func init() {
	// change to your database configuration file path, see dbconfig-sample.json
	connStr, err := LoadDBConf("dbconfig.json")
	if err != nil {
		fmt.Printf("cannot open database configuration file, %v", err)
		os.Exit(1)
	}
	connPool, err = sql.Open("mysql", connStr)
	if err != nil {
		fmt.Printf("cannot create database connection pool, %v", err)
		os.Exit(1)
	}
	err = connPool.Ping()
	if err != nil {
		fmt.Printf("cannot access database, %v", err)
		os.Exit(1)
	}

	connPool.SetMaxOpenConns(maxOpenConn)
	connPool.SetMaxIdleConns(maxIdleConn)
	connPool.SetConnMaxLifetime(maxConnLifeTime)
}

func CloseDatabase() {
	if connPool != nil {
		connPool.Close()
	}
}

func InsertVideo(video Video) error {
	tx, err := connPool.Begin()
	if err != nil{
		return errors.New("transaction begin failed : " + err.Error())
	}

	var pubdate interface{}
	if video.Pubdate.IsZero(){
		pubdate = nil
	}else{
		pubdate = video.Pubdate
	}

	_, err = tx.Exec(
		"INSERT INTO video VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		video.Aid,
		video.Status,
		video.Title,
		pubdate,
		video.Owner,
		video.Duration,
		video.View,
		video.Dannmaku,
		video.Reply,
		video.Favorite,
		video.Coin,
		video.Share,
		video.Now_rank,
		video.His_rank,
		video.Support,
		video.Dislike,
		video.No_reprint,
		video.Copyright)
	if err != nil {
		return rollback(tx, err)
	}

	for _, page := range video.Pages{
		_, err = tx.Exec(
			"INSERT INTO pages VALUES(?, ?, ?, ?, ?);",
			video.Aid,
			page.PageNo,
			page.Chatid,
			page.Duration,
			page.Subtitle,
		)
		if err != nil {
			return rollback(tx, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("transaction Commit failed : " + err.Error())
	}
	return nil
}

func rollback(tx *sql.Tx, oldErr error) error{
	err := tx.Rollback()
	if err != nil{
		return errors.New(oldErr.Error() + " : " + err.Error())
	}
	return oldErr
}
