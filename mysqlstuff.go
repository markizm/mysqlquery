package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

//type ReleaseTable struct {
//	relTicket string
//	cmTicket  string
//	startTime string
//	endTime   string
//	progress  string
//	status    string
//}

func main() {
	http.HandleFunc("/releaseTable", releaseTable)
	http.ListenAndServe(":9002", nil)
}

func releaseTable(w http.ResponseWriter, r *http.Request) {
	//open command to connect to db.. need to hide password
	db, err := sql.Open("mysql", "<user>:<pw>@tcp(localhost:3306)/dbname")
	if err != nil {
		fmt.Println(err)
	}
	//basic sql query
	rows, err := db.Query("SELECT * FROM dbtable;")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	//creates object for data in each column
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		//gets hairy here but it seems to work
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	//marshal/unmarshal data from db, will likely break this out into another func
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err)
	}

	var e Events
	json.Unmarshal(jsonData, &e)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(e)
}

type Events []interface{}
