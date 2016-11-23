// Copyright 2016 Etix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDB contains the MySQL configuration
type MysqlDB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

func (t *Tester) dropDB() bool {
	dsn := t.DB.User + ":" + t.DB.Password + "@" + "tcp(" + t.DB.Host + ":" + strconv.Itoa(t.DB.Port) + ")/" + t.DB.DbName + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	q := "DROP DATABASE " + t.DB.DbName + ";"
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("------ Dropped Cameradar Database -------")
	return true
}

func (t *Tester) configureDatabase(DataBase *MysqlDB) bool {
	var db MysqlDB

	db.Host = t.Cameradar.DbHost
	db.Port = t.Cameradar.DbPort
	db.User = t.Cameradar.DbUser
	db.Password = t.Cameradar.DbPassword
	db.DbName = t.Cameradar.DbName

	*DataBase = db
	return true
}
