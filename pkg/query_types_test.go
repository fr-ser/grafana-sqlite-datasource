package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

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

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

// TestSimpleTableQuery tests against a table with no provided type information from the frontend
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

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func pointerOf(x float64) *float64 {
	return &x
}
