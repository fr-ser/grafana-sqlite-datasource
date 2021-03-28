package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestEpochGroupSecondsShouldBeReplacedInTheFinalQueryForTables(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (34, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10) as window, value FROM test`,
		TimeColumns: []string{},
	})
	dataQuery.QueryType = tableType

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
		data.NewField("window", nil, []*int64{intPointer(0), intPointer(10), intPointer(30)}),
		data.NewField("value", nil, []*int64{intPointer(1), intPointer(2), intPointer(4)}),
	)

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestEpochGroupSecondsShouldBeReplacedInTheFinalQueryForTimeSeries(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value INTEGER);
		INSERT INTO test(time, value)
		VALUES (4, 1), (13, 2), (34, 4);
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   `SELECT $__unixEpochGroupSeconds("time", 10) as window, value FROM test`,
		TimeColumns: []string{"window"},
	})
	dataQuery.QueryType = timeSeriesType

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
		"value",
		data.NewField("window", nil, []*time.Time{
			timePointer(time.Unix(0, 0)),
			timePointer(time.Unix(10, 0)),
			timePointer(time.Unix(30, 0)),
		}),
		data.NewField("value", nil, []*int64{intPointer(1), intPointer(2), intPointer(4)}),
	)

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}
