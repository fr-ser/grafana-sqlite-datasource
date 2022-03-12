package plugin

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestQueryWithTimeColumn(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, value REAL, name TEXT);
		INSERT INTO test(time, value, name)
		VALUES (21.0, 21.1, 'one'), (22.3, 22.2, 'two'), (23, 23.3, 'three');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
	})

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
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(22, 300000000)),
			timePointer(time.Unix(23, 0)),
		}),
		data.NewField("value", nil, []*float64{
			floatPointer(21.1), floatPointer(22.2), floatPointer(23.3),
		}),
		data.NewField("name", nil, []*string{
			strPointer("one"), strPointer("two"), strPointer("three"),
		}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

func TestQueryWithTimeStringColumn(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time INTEGER, time_string TEXT);
		INSERT INTO test(time, time_string)
		VALUES	(21, '1970-01-01T00:00:21Z'),
				(22.3, '1970-01-01T00:00:22.300Z'),
				(23, '1970-01-01T00:00:23Z');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time", "time_string"},
	})

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
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(22, 300000000)),
			timePointer(time.Unix(23, 0)),
		}),
		data.NewField("time_string", nil, []*time.Time{
			timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(22, 300000000)),
			timePointer(time.Unix(23, 0)),
		}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}

// TestUnixTimestampAsString tests that unix timestamps as strings are also formatted correctly
// this case is common in SQLite as the only way to get a UNIX timestamp dynamically requires
// to convert to a string, e.g. strftime('%s', 'now')
func TestUnixTimestampAsString(t *testing.T) {
	dbPath, cleanup := createTmpDB(`
		CREATE TABLE test(time TEXT);
		INSERT INTO test VALUES	('21'), ('22.3'), ('23');
	`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText: "SELECT * FROM test", TimeColumns: []string{"time"},
	})

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
		data.NewField("time", nil, []*time.Time{
			timePointer(time.Unix(21, 0)),
			timePointer(time.Unix(22, 300000000)),
			timePointer(time.Unix(23, 0)),
		}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT * FROM test"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}
