package plugin

import (
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// replaceVariables replaces Grafana Template Variables in the query
// this is mainly used for alert queries, which need time replacement
func replaceVariables(queryConfig *queryConfigStruct, dataQuery backend.DataQuery) error {
	queryConfig.FinalQuery = strings.ReplaceAll(
		queryConfig.FinalQuery,
		"$__from",
		strconv.FormatInt(dataQuery.TimeRange.From.Unix()*1000, 10),
	)
	queryConfig.FinalQuery = strings.ReplaceAll(
		queryConfig.FinalQuery,
		"$__to",
		strconv.FormatInt(dataQuery.TimeRange.To.Unix()*1000, 10),
	)
	return nil
}
