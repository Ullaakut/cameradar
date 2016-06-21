package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strconv"
)

func (m *manager) dropDB() bool {
	dsn := m.DB.User + ":" + m.DB.Password + "@" + "tcp(" + m.DB.Host + ":" + strconv.Itoa(m.DB.Port) + ")/" + m.DB.Db_name + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	q := "DROP DATABASE cctv;"
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("------ Dropped CCTV Database -------")
	return true
}