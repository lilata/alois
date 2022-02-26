package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/exec"
)

func getDb() *sql.DB {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal("Error on opening DB.", err)
		return nil
	}
	return db
}
func checkDb() bool {
	_, err := os.Stat(DBPath)
	return err == nil
}
func initDb() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("sqlite3 %s < schema.sql", DBPath))
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error on initing database.", err)
	}
	return err
}