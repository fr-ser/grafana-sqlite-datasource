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

// this struct holds a full query result column (including data)
// the main benefit is type safety
type sqlColumn struct {
	// Name is the name of the column
	Name string

	// Type is the data type of the column.
	// This determines which row property to use
	Type string

	// TimeData contains time values (if Type == "TIME")
	TimeData []time.Time
	// IntData contains integer values (if Type == "INTEGER")
	IntData []int64

	// FloatData contains float values (if Type == "FLOAT")
	FloatData []float64

	// StringData contains string values (if Type == "STRING")
	StringData []string
}

func transformRow(rows *sql.Rows, columns []*sqlColumn) error {
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
		default:
			log.DefaultLogger.Warn(
				"Scanned row value type was unexpected",
				"value", values[i], "type", fmt.Sprintf("%T", values[i]),
			)
			valueType = "UNKNOWN"
		}

		if column.Type == "UNKNOWN" && valueType != "UNKNOWN" {
			column.Type = valueType
		}

		if valueType == "INTEGER" && column.Type == "TIME" {
			columns[i].TimeData = append(columns[i].TimeData, time.Unix(intV, 0))
		} else if valueType == "FLOAT" && column.Type == "TIME" {
			columns[i].TimeData = append(columns[i].TimeData, time.Unix(int64(floatV), 0))
		} else if column.Type == "TIME" {
			val := fmt.Sprintf("%v", values[i])
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				// try parsing the string as a number
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					t = time.Unix(int64(f), 0)
				} else {
					log.DefaultLogger.Warn(
						"Could parse (RFC3339) value to timestamp", "value", val,
					)
				}
			}
			columns[i].TimeData = append(columns[i].TimeData, t)
		} else if valueType == "INTEGER" && column.Type == "INTEGER" {
			columns[i].IntData = append(columns[i].IntData, intV)
		} else if column.Type == "INTEGER" {
			if v, err := strconv.ParseInt(string(stringV), 10, 64); err == nil {
				columns[i].IntData = append(columns[i].IntData, v)
			} else {
				log.DefaultLogger.Warn("Could not convert numeric to float", "value", stringV)
			}
		} else if valueType == "FLOAT" && column.Type == "FLOAT" {
			columns[i].FloatData = append(columns[i].FloatData, floatV)
		} else if column.Type == "FLOAT" {
			if v, err := strconv.ParseFloat(string(stringV), 64); err == nil {
				columns[i].FloatData = append(columns[i].FloatData, v)
			} else {
				log.DefaultLogger.Warn("Could not convert numeric to float", "value", stringV)
			}
		} else if column.Type == "STRING" && valueType == "INTEGER" {
			columns[i].StringData = append(columns[i].StringData, fmt.Sprintf("%d", intV))
		} else if column.Type == "STRING" && valueType == "FLOAT" {
			columns[i].StringData = append(columns[i].StringData, fmt.Sprintf("%f", floatV))
		} else {
			columns[i].StringData = append(columns[i].StringData, fmt.Sprintf("%v", values[i]))
		}
	}

	return nil
}

func fetchData(dbPath string, qm queryModel) (columns []*sqlColumn, err error) {

	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		log.DefaultLogger.Error("Could not open database", "err", err)
		return columns, err
	}

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
		case "REAL":
			columns[idx].Type = "FLOAT"
		case "NULL", "TEXT", "BLOB":
			columns[idx].Type = "STRING"
		case "", "UNKNOWN":
			columns[idx].Type = "UNKNOWN"
		default:
			log.DefaultLogger.Debug(
				"Unknown database type", "type", columnTypes[idx].DatabaseTypeName(),
			)
			columns[idx].Type = "STRING"
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

	// add the frames to the response
	response.Frames = append(response.Frames, frame)

	return response
}
