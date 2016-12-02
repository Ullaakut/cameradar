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

func (t *Tester) dropDB() bool {
	dsn := t.ServiceConf.User + ":" + t.ServiceConf.Password + "@" + "tcp(" + t.ServiceConf.Host + ":" + strconv.Itoa(t.ServiceConf.Port) + ")/" + t.ServiceConf.DbName + "?charset=utf8"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	q := "DROP DATABASE " + t.ServiceConf.DbName + ";"
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------ Dropped Cameradar Database -------")
	return true
}
