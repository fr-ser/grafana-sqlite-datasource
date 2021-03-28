package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var mockableLongToWide = data.LongToWide

// commented out as currently unused
// const timeSeriesType = "time series"
const tableType = "table"

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

func transformRow(rows *sql.Rows, columns []*sqlColumn) (err error) {
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
			log.DefaultLogger.Warn(
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
				value = time.Unix(int64(floatV), 0)
			} else if valueType != "NULL" {
				val := fmt.Sprintf("%v", values[i])
				value, err = time.Parse(time.RFC3339, val)
				if err != nil {
					// try parsing the string as a number
					if f, err := strconv.ParseFloat(val, 64); err == nil {
						value = time.Unix(int64(f), 0)
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
					log.DefaultLogger.Warn("Could not convert value to int", "value", stringV)
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
					log.DefaultLogger.Warn("Could not convert value to float", "value", stringV)
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

func fetchData(dbPath string, qm queryModel) (columns []*sqlColumn, err error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.DefaultLogger.Error("Could not open database", "err", err)
		return columns, err
	}
	defer db.Close()

	rows, err := db.Query(qm.QueryText)
	if err != nil {
		log.DefaultLogger.Error("Could execute query", "query", qm.QueryText, "err", err)
		return columns, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.DefaultLogger.Error("Could not get column types", "err", err)
		return columns, err
	}
	columnCount := len(columnTypes)
	columns = make([]*sqlColumn, columnCount)

	if columnCount == 0 {
		return columns, nil
	}

	for idx := range columns {
		columns[idx] = &sqlColumn{Name: columnTypes[idx].Name()}

		switch columnTypes[idx].DatabaseTypeName() {
		case "INTEGER":
			columns[idx].Type = "INTEGER"
		case "REAL", "NUMERIC":
			columns[idx].Type = "FLOAT"
		case "NULL", "TEXT", "BLOB":
			columns[idx].Type = "STRING"
		default:
			log.DefaultLogger.Debug(
				"Unknown database type", "type", columnTypes[idx].DatabaseTypeName(),
			)
			columns[idx].Type = "UNKNOWN"
		}

		for _, timeColumnName := range qm.TimeColumns {
			if columns[idx].Name == timeColumnName {
				columns[idx].Type = "TIME"
				break
			}
		}
	}

	for rows.Next() {
		err := transformRow(rows, columns)
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

	columns, err := fetchData(config.Path, qm)
	if err != nil {
		response.Error = err
		return response
	}

	frame := data.NewFrame("response")

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
		default:
			frame.Fields = append(
				frame.Fields, data.NewField(column.Name, nil, column.StringData),
			)
		}
	}

	// default "table" case. Return whatever SQL we received
	if dataQuery.QueryType == tableType || dataQuery.QueryType == "" {
		response.Frames = append(response.Frames, frame)

		return response
	}
	// as the QueryType currently only has two options no further "if check" is required

	if frame.TimeSeriesSchema().Type != data.TimeSeriesTypeWide {
		frame, err = mockableLongToWide(frame, nil)
		if err != nil {
			response.Error = err
			return response
		}
	}

	// some plugins do not play well with the "wide format" of a time series
	// therefore we transform into individual frames
	// https://github.com/fr-ser/grafana-sqlite-datasource/issues/16
	tsSchema := frame.TimeSeriesSchema()
	for _, field := range frame.Fields {
		if field.Type().Time() {
			continue
		}
		partialFrame := data.NewFrame(
			fmt.Sprintf("%s %s", field.Name, field.Labels["name"]),
			frame.Fields[tsSchema.TimeIndex],
			field,
		)

		response.Frames = append(response.Frames, partialFrame)
	}

	return response
}
