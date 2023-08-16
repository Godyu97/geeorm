package main

import (
	"fmt"
	"github.com/Godyu97/geeorm/gee"
	"github.com/Godyu97/geeorm/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Info("test geeorm!🚀")
	e, _ := gee.NewEngine("sqlite3", "gee.db")
	defer e.Close()
	s := e.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
