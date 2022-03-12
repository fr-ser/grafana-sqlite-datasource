package plugin

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// TestJsonSupport tests against a query with a json function to make sure, that the plugin is
// built with the JSON extension
func TestJsonSupport(t *testing.T) {
	dbPath, cleanup := createTmpDB("SELECT 1 -- create db")
	defer cleanup()

	baseQuery := `SELECT json_array_length('[1,2,3,4]') as value;`
	dataQuery := getDataQuery(queryModel{QueryText: baseQuery})

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
		data.NewField("value", nil, []*int64{intPointer(4)}),
	)
	expectedFrame.Meta = &data.FrameMeta{ExecutedQueryString: baseQuery}

	if diff := cmp.Diff(expectedFrame, response.Frames[0], cmpOption...); diff != "" {
		t.Error(diff)
	}
}
