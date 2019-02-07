package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/xackery/eqemuconfig"

	// mysql is used for support
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	config, err := eqemuconfig.GetConfig()
	if err != nil {
		return
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
	if err != nil {
		return
	}
	defer db.Close()

	f, err := os.Open("../Release/Dumps/phinterior3a3_NPC_2019-02-06-16-39-10.csv")
	if err != nil {
		return
	}
	fileName := path.Base(f.Name())
	if len(fileName) < 3 {
		err = fmt.Errorf("file name must be greater than 3 characters long")
		return
	}
	if !strings.Contains(fileName, "_") {
		err = fmt.Errorf("file name should have a _ seperator")
		return
	}
	zoneName := fileName[0:strings.Index(fileName, "_")]
	if len(zoneName) < 1 {
		err = fmt.Errorf("zone name failed to parse")
		return
	}

	zoneID := 0
	row := db.QueryRow("SELECT zoneidnumber FROM zone WHERE short_name = ?", strings.ToLower(zoneName))
	err = row.Scan(&zoneID)
	if err != nil {
		return
	}
	fmt.Println("zone ID:", zoneID)

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return
	}

	//fmt.Print(records)
	if len(records) > 0 {
		fmt.Println("success")
		return
	}
	return
}
