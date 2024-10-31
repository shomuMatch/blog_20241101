package test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func TestApi(t *testing.T) {
	db, err := gorm.Open("postgres", "user=postgres port=5432 host=localhost password=password sslmode=disable")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer db.Close()
	db.Exec("CREATE EXTENSION IF NOT EXISTS pg_raise_error")
	testCases := []struct {
		count int
		level string
		body  string
	}{
		{-1, "info", "OK\n"},
		{1, "error", "Internal Server Error: first: pq: COUNT IS ZERO!!\n"},
		{2, "error", "Internal Server Error: second: pq: COUNT IS ZERO!!\n"},
		{3, "error", "Internal Server Error: third: pq: COUNT IS ZERO!!\n"},
		{4, "error", "Internal Server Error: fourth: pq: COUNT IS ZERO!!\n"},
		{5, "error", "OK\n"},
	}
	for caseNo, testCase := range testCases {
		t.Run(fmt.Sprint(caseNo), func(t *testing.T) {
			db.Exec("select set_error_trigger(?,?)", testCase.count, testCase.level)
			url := "http://localhost:8888/api"
			res, _ := http.Get(url)
			resBody, _ := io.ReadAll(res.Body)
			db.Exec("select clear_error_trigger();")

			if testCase.body != string(resBody) {
				t.Errorf("expect: %v, actual: %v", testCase.body, string(resBody))
			}
		})
	}
}
