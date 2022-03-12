package plugin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var mockableLongToWide = data.LongToWide

const timeSeriesType = "time series"
const tableType = "table"

type queryConfigStruct struct {
	BaseQuery   string
	TimeColumns []string
	QueryType   string
	FinalQuery  string

	ShouldFillValues          bool
	FillInterval              int
	FillValuesTimeColumnIndex int
}

func (qc *queryConfigStruct) isTableType() bool {
	return qc.QueryType != timeSeriesType
}

// this struct holds a full query result column (including data)
// the main benefit is type safety
type sqlColumn struct {
	// Name is the name of the column
	Name string

	// Type is the data type of the column.
	// This determines which row property to use
	Type string

	// TimeData contains time values (if Type == "TIME")
	TimeData []*time.Time
	// IntData contains integer values (if Type == "INTEGER")
	IntData []*int64

	// FloatData contains float values (if Type == "FLOAT")
	FloatData []*float64

	// StringData contains string values (if Type == "STRING")
	StringData []*string
}

func addTransformedRow(rows *sql.Rows, columns []*sqlColumn) (err error) {
	columnCount := len(columns)
	values := make([]interface{}, columnCount)
	valuePointers := make([]interface{}, columnCount)

	for i := 0; i < columnCount; i++ {
		valuePointers[i] = &values[i]
	}

	if err := rows.Scan(valuePointers...); err != nil {
		log.DefaultLogger.Error("Could not scan row", "err", err)
		return err
	}

	for i, column := range columns {
		var intV int64
		var floatV float64
		var stringV string
		valueType := ""

		switch v := values[i].(type) {
		case int8:
			valueType = "INTEGER"
			intV = int64(v)
		case int16:
			valueType = "INTEGER"
			intV = int64(v)
		case int32:
			valueType = "INTEGER"
			intV = int64(v)
		case int64:
			valueType = "INTEGER"
			intV = v
		case float32:
			valueType = "FLOAT"
			floatV = float64(v)
		case float64:
			valueType = "FLOAT"
			floatV = v
		case []byte:
			valueType = "STRING"
			stringV = string(v)
		case string:
			valueType = "STRING"
			stringV = v
		case nil:
			valueType = "NULL"
		default:
			log.DefaultLogger.Debug(
				"Scanned row value type was unexpected",
				"value", values[i], "type", fmt.Sprintf("%T", values[i]),
				"column", column.Name,
			)
			valueType = "UNKNOWN"
		}

		if column.Type == "UNKNOWN" && valueType != "NULL" {
			// we need to decide on a type for the column as we need to
			// fill a typed list later. Multiple types are not allowed
			if valueType == "UNKNOWN" {
				column.Type = "STRING"
			} else {
				column.Type = valueType
			}
		}

		// variable to indicate whether to explicitly set the value to null
		setNull := false

		if column.Type == "TIME" {
			var value time.Time

			if valueType == "INTEGER" {
				value = time.Unix(intV, 0)
			} else if valueType == "FLOAT" {
				seconds, milliseconds := math.Modf(floatV)
				value = time.Unix(int64(seconds), int64(milliseconds*1000000000))
			} else if valueType != "NULL" {
				val := fmt.Sprintf("%v", values[i])
				value, err = time.Parse(time.RFC3339, val)
				if err != nil {
					// try parsing the string as a number
					if f, err := strconv.ParseFloat(val, 64); err == nil {
						seconds, milliseconds := math.Modf(f)
						value = time.Unix(int64(seconds), int64(milliseconds*1000000000))
					} else {
						log.DefaultLogger.Warn(
							"Could not parse (RFC3339) value to timestamp", "value", val,
						)
						setNull = true
					}
				}
			}

			if setNull || valueType == "NULL" {
				columns[i].TimeData = append(columns[i].TimeData, nil)
			} else {
				columns[i].TimeData = append(columns[i].TimeData, &value)
			}
			continue
		}

		if column.Type == "INTEGER" {
			var value int64

			if valueType == "INTEGER" {
				value = intV
			} else if valueType == "FLOAT" {
				value = int64(floatV)
			} else {
				value, err = strconv.ParseInt(string(stringV), 10, 64)
				if err != nil {
					log.DefaultLogger.Debug("Could not convert value to int", "value", stringV)
					setNull = true
				}
			}

			if setNull || valueType == "NULL" {
				columns[i].IntData = append(columns[i].IntData, nil)
			} else {
				columns[i].IntData = append(columns[i].IntData, &value)
			}
			continue
		}

		if column.Type == "FLOAT" {
			var value float64

			if valueType == "FLOAT" {
				value = floatV
			} else if valueType == "INTEGER" {
				value = float64(intV)
			} else {
				value, err = strconv.ParseFloat(string(stringV), 64)

				if err != nil {
					log.DefaultLogger.Debug("Could not convert value to float", "value", stringV)
					setNull = true
				}
			}

			if setNull || valueType == "NULL" {
				columns[i].FloatData = append(columns[i].FloatData, nil)
			} else {
				columns[i].FloatData = append(columns[i].FloatData, &value)
			}
			continue
		}

		if column.Type == "STRING" {
			var value string

			if valueType == "INTEGER" {
				value = fmt.Sprintf("%d", intV)
			} else if valueType == "FLOAT" {
				value = fmt.Sprintf("%f", floatV)
			} else {
				value = fmt.Sprintf("%v", values[i])
			}

			if valueType == "NULL" {
				columns[i].StringData = append(columns[i].StringData, nil)
			} else {
				columns[i].StringData = append(columns[i].StringData, &value)
			}
			continue
		}

		// column.Type == "UNKNOWN"
		columns[i].TimeData = append(columns[i].TimeData, nil)
		columns[i].IntData = append(columns[i].IntData, nil)
		columns[i].FloatData = append(columns[i].FloatData, nil)
		columns[i].StringData = append(columns[i].StringData, nil)

	}

	return nil
}

func fetchData(
	dbPathPrefix string, dbPath string, dbPathOptions string, queryConfig *queryConfigStruct,
) (columns []*sqlColumn, err error) {
	db, err := sql.Open("sqlite", dbPathPrefix+dbPath+"?"+dbPathOptions)
	if err != nil {
		log.DefaultLogger.Error("Could not open database", "err", err)
		return columns, err
	}
	defer db.Close()

	rows, err := db.Query(queryConfig.FinalQuery)
	if err != nil {
		log.DefaultLogger.Error(
			"Could not execute query", "query", queryConfig.FinalQuery, "err", err,
		)
		return columns, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil && err.Error() == "sql: no Rows available" {
		return make([]*sqlColumn, 0), nil
	} else if err != nil {
		log.DefaultLogger.Error("Could not get column types", "err", err)
		return columns, err
	}
	// closing on an empty set "sql: no Rows available" causes a panic
	defer rows.Close()

	columnCount := len(columnTypes)
	columns = make([]*sqlColumn, columnCount)

	if columnCount == 0 {
		return columns, nil
	}

	for idx := range columns {
		columns[idx] = &sqlColumn{Name: columnTypes[idx].Name()}

		switch strings.ToUpper(columnTypes[idx].DatabaseTypeName()) {
		case "INTEGER", "INT":
			columns[idx].Type = "INTEGER"
		case "REAL", "NUMERIC", "DOUBLE", "FLOAT":
			columns[idx].Type = "FLOAT"
		case "NULL", "TEXT", "BLOB":
			columns[idx].Type = "STRING"
		default:
			log.DefaultLogger.Debug(
				"Unknown database type",
				"type",
				columnTypes[idx].DatabaseTypeName(),
				"column",
				columnTypes[idx].Name(),
			)
			columns[idx].Type = "UNKNOWN"
		}

		for _, timeColumnName := range queryConfig.TimeColumns {
			if columns[idx].Name == timeColumnName {
				columns[idx].Type = "TIME"
				if queryConfig.FillValuesTimeColumnIndex == -1 {
					queryConfig.FillValuesTimeColumnIndex = idx
				}
				break
			}
		}
	}

	for rows.Next() {
		err := addTransformedRow(rows, columns)
		if err != nil {
			return columns, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.DefaultLogger.Error("The row scan finished with an error", "err", err)
		return columns, err
	}

	return columns, nil
}

type queryModel struct {
	QueryText   string   `json:"queryText"`
	TimeColumns []string `json:"timeColumns"`
}

func query(dataQuery backend.DataQuery, config pluginConfig) (response backend.DataResponse) {
	var qm queryModel
	err := json.Unmarshal(dataQuery.JSON, &qm)
	if err != nil {
		log.DefaultLogger.Error("Could not unmarshal query", "err", err)
		response.Error = err
		return response
	}

	queryConfig := queryConfigStruct{
		BaseQuery:                 qm.QueryText,
		FinalQuery:                qm.QueryText,
		TimeColumns:               qm.TimeColumns,
		QueryType:                 dataQuery.QueryType,
		FillValuesTimeColumnIndex: -1,
	}

	err = replaceVariables(&queryConfig, dataQuery)
	if err != nil {
		response.Error = err
		return response
	}
	log.DefaultLogger.Debug("Variables replaced")

	err = applyMacros(&queryConfig)
	if err != nil {
		response.Error = err
		return response
	}
	log.DefaultLogger.Debug("Macros applied")

	columns, err := fetchData(config.PathPrefix, config.Path, config.PathOptions, &queryConfig)
	if err != nil {
		response.Error = err
		return response
	}
	log.DefaultLogger.Debug("Fetched data from database")

	frame := data.NewFrame("response")
	frame.Meta = &data.FrameMeta{ExecutedQueryString: queryConfig.FinalQuery}

	if queryConfig.ShouldFillValues {
		err := fillGaps(columns, &queryConfig)
		if err != nil {
			response.Error = err
			return response
		}
		log.DefaultLogger.Debug("Filled gaps in data according to macro")
	}

	// construct a regular SQL dataframe (for time series this is usually the "long format")
	for _, column := range columns {
		switch column.Type {
		case "TIME":
			frame.Fields = append(
				frame.Fields, data.NewField(column.Name, nil, column.TimeData),
			)
		case "FLOAT":
			frame.Fields = append(
				frame.Fields, data.NewField(column.Name, nil, column.FloatData),
			)
		case "INTEGER":
			frame.Fields = append(
				frame.Fields, data.NewField(column.Name, nil, column.IntData),
			)
		case "STRING":
			frame.Fields = append(
				frame.Fields, data.NewField(column.Name, nil, column.StringData),
			)
		default:
			frame.Fields = append(
				frame.Fields, data.NewField(column.Name, nil, column.FloatData),
			)
		}
	}

	// default case. Return whatever SQL we received
	if queryConfig.isTableType() {
		response.Frames = append(response.Frames, frame)

		return response
	}
	// as the QueryType currently only has two options no further "if check" is required

	if frame.TimeSeriesSchema().Type != data.TimeSeriesTypeWide {
		frame, err = mockableLongToWide(frame, nil)
		if err != nil {
			log.DefaultLogger.Error("Could not convert from long to wide time-series", "err", err)
			response.Error = err
			return response
		}
		frame.Meta = &data.FrameMeta{ExecutedQueryString: queryConfig.FinalQuery}

		log.DefaultLogger.Debug("Initial data converted into wide time-series")

		// bug? if there is row with only null values "NULL" gets added as a "factor"
		// this adds an empty string named null only field to the response. Here we remove it
		emptyFieldIndexes := map[int]bool{}
		for idx, field := range frame.Fields {
			if field.Labels["name"] == "" && fieldHasOnlyNulls(field) {
				emptyFieldIndexes[idx] = true
			}
		}

		if len(emptyFieldIndexes) > 0 {
			filledFields := []*data.Field{}
			for idx, field := range frame.Fields {
				if !emptyFieldIndexes[idx] {
					filledFields = append(filledFields, field)
				}
			}
			frame.Fields = filledFields
			log.DefaultLogger.Debug("Removed null field from generated time-series dataframe")

		}
	}

	// some plugins do not play well with the "wide format" of a time series
	// therefore we transform into individual frames
	// https://github.com/fr-ser/grafana-sqlite-datasource/issues/16
	tsSchema := frame.TimeSeriesSchema()
	for idx, field := range frame.Fields {
		if idx == tsSchema.TimeIndex {
			continue
		}
		partialFrame := data.NewFrame(
			strings.Trim(fmt.Sprintf("%s %s", field.Name, field.Labels["name"]), " "),
			frame.Fields[tsSchema.TimeIndex],
			field,
		)
		partialFrame.Meta = &data.FrameMeta{ExecutedQueryString: queryConfig.FinalQuery}

		response.Frames = append(response.Frames, partialFrame)
	}
	log.DefaultLogger.Debug("Wide time-series converted into multiple frames")

	return response
}

func fieldHasOnlyNulls(field *data.Field) bool {
	for row := 0; row < field.Len(); row++ {
		if _, isNil := field.ConcreteAt(row); isNil {
			return false
		}
	}
	return true
}
