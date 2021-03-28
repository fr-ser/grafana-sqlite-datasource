package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var cmpOption = data.FrameTestCompareOptions()

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

func strPointer(x string) *string {
	return &x
}

func floatPointer(x float64) *float64 {
	return &x
}

func intPointer(x int64) *int64 {
	return &x
}

func timePointer(x time.Time) *time.Time {
	return &x
}

func unixTimePointer(x int64) *time.Time {
	return timePointer(time.Unix(x, 0))
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
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "-- not a query"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Errorf(diff)
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
		data.NewField("time", nil, []*int64{}),
		data.NewField("value", nil, []*float64{}),
		data.NewField("name", nil, []*string{}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test WHERE false"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
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
