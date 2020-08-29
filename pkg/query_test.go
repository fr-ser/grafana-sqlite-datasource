package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/pretty"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func createTmpDB(seedSQL string) (dbPath string, cleanup func()) {
	dir, _ := ioutil.TempDir("", "test-check-db")
	dbPath = filepath.Join(dir, "data.db")
	db, _ := sql.Open("sqlite3", dbPath)
	db.Exec(seedSQL)
	db.Close()
	cleanup = func() { os.RemoveAll(dir) }

	return dbPath, cleanup
}

func getDataQuery(targetModel queryModel) backend.DataQuery {
	jsonData, _ := json.Marshal(targetModel)
	return backend.DataQuery{JSON: jsonData}
}

// TestCTETableQuery tests against a query with no provided type information from the frontend
func TestCTETableQuery(t *testing.T) {
	dbPath, cleanup := createTmpDB("SELECT 1 -- create db")
	defer cleanup()

	dataQuery := getDataQuery(queryModel{QueryText: `
		WITH some_tmp_table(time, value, name) AS (
			SELECT * FROM (VALUES (1, 1.1, 'one'), (2, 2.2, 'two'), (3, 3.3, 'three'))
		)
		SELECT * FROM some_tmp_table
	`})

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []int64{1, 2, 3}),
		data.NewField("value", nil, []float64{1.1, 2.2, 3.3}),
		data.NewField("name", nil, []string{"one", "two", "three"}),
	)

	diff := pretty.Compare(expectedFrame, response.Frames[0])
	if diff != "" {
		t.Error(diff)
	}
}

// TestSimpleTableQuery tests against a tablewith no provided type information from the frontend
func TestSimpleTableQuery(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21, 21.1, 'one'), (22, 22.2, 'two'), (23, 23.3, 'three');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{QueryText: "SELECT * FROM test"})

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []int64{21, 22, 23}),
		data.NewField("value", nil, []float64{21.1, 22.2, 23.3}),
		data.NewField("name", nil, []string{"one", "two", "three"}),
	)

	diff := pretty.Compare(expectedFrame, response.Frames[0])
	if diff != "" {
		t.Error(diff)
	}
}

func TestEmptyQuery(t *testing.T) {
	dbPath, cleanup := createTmpDB(`SELECT 1`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{QueryText: "-- not a query"})

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame("response")

	diff := pretty.Compare(expectedFrame, response.Frames[0])
	if diff != "" {
		t.Error(diff)
	}
}

func TestNoResults(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21, 21.1, 'one'), (22, 22.2, 'two'), (23, 23.3, 'three');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{QueryText: "SELECT * FROM test WHERE false"})

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	expectedFrame := data.NewFrame(
		"response",
		data.NewField("time", nil, []int64{}),
		data.NewField("value", nil, []float64{}),
		data.NewField("name", nil, []string{}),
	)

	diff := pretty.Compare(expectedFrame, response.Frames[0])
	if diff != "" {
		t.Error(diff)
	}
}

func TestInvalidQuery(t *testing.T) {
	dbPath, cleanup := createTmpDB("SELECT 1 -- create db")
	defer cleanup()

	dataQuery := backend.DataQuery{JSON: []byte(`not even json`)}

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error == nil {
		t.Errorf("Expected unmarshal error but got nothing. Response: %+v", response)
	}
}
