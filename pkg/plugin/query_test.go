package plugin

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
	db, _ := sql.Open("sqlite", dbPath)
	_, _ = db.Exec(seedSQL)
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

func TestNoResultsTable(t *testing.T) {
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

func TestNoResultsTimeSeriesWithUnknownColumns(t *testing.T) {
	var longToWideCalled bool
	mockableLongToWide = func(a *data.Frame, b *data.FillMissing) (*data.Frame, error) {
		longToWideCalled = true
		return data.NewFrame("response"), nil
	}

	dbPath, cleanup := createTmpDB(`SELECT 1`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{
		QueryText:   "SELECT 1 as time, 2 as value WHERE FALSE",
		TimeColumns: []string{"time"},
	})
	dataQuery.QueryType = "time series"

	response := query(dataQuery, pluginConfig{Path: dbPath})
	if response.Error != nil {
		t.Errorf("Unexpected error - %s", response.Error)
	}

	if len(response.Frames) != 1 {
		t.Errorf(
			"Expected one frame but got - %d: Frames %+v", len(response.Frames), response.Frames,
		)
	}

	if longToWideCalled {
		t.Errorf("Expected to not call 'longToWide' for wide time series queries")
	}

	expectedFrame := data.NewFrame(
		"value",
		data.NewField("time", nil, []*time.Time{}),
		data.NewField("value", nil, []*float64{}),
	)
	expectedFrame.Meta = &data.FrameMeta{
		ExecutedQueryString: "SELECT 1 as time, 2 as value WHERE FALSE",
	}

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

func TestReplaceToAndFromVariables(t *testing.T) {
	dbPath, cleanup := createTmpDB(`SELECT 1`)
	defer cleanup()

	dataQuery := getDataQuery(queryModel{QueryText: "SELECT $__from AS a, $__to AS b"})
	dataQuery.TimeRange.From = time.Unix(123, 0)
	dataQuery.TimeRange.To = time.Unix(456, 0)

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
		data.NewField("a", nil, []*int64{intPointer(123000)}),
		data.NewField("b", nil, []*int64{intPointer(456000)}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: "SELECT 123000 AS a, 456000 AS b"}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Errorf(diff)
	}
}
