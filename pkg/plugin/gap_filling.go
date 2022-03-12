package plugin

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func fillGaps(columns []*sqlColumn, queryConfig *queryConfigStruct) error {
	if queryConfig.FillValuesTimeColumnIndex == -1 {
		return fmt.Errorf("no time column found to use for gap filling")
	}

	timeColumn := &columns[queryConfig.FillValuesTimeColumnIndex].TimeData
	if len(*timeColumn) < 2 {
		// gaps cannot be filled. This a not an error but a noop
		return nil
	}

	gapFilledColumns := make([]*sqlColumn, len(columns))
	for idx, column := range columns {
		gapFilledColumns[idx] = &sqlColumn{Name: column.Name, Type: column.Type}
	}

	// substract one interval to start the gap filling after the first value
	currentRowTime := (*timeColumn)[0].Add(
		time.Second * -1 * time.Duration(queryConfig.FillInterval),
	)

	for idx, timeCell := range *timeColumn {
		if timeCell == nil {
			return fmt.Errorf("received NULL value in time column")
		}
		if idx > 0 && (*timeColumn)[idx-1].After(*timeCell) {
			return fmt.Errorf(
				"received unordered time value. The time column needs to be sorted ascending",
			)

		}
		currentRowTime = currentRowTime.Add(time.Second * time.Duration(queryConfig.FillInterval))

		for currentRowTime.Before(*timeCell) {
			err := addGapToColumns(
				currentRowTime, gapFilledColumns, columns, *queryConfig, len(*timeColumn),
			)
			if err != nil {
				return err
			}
			currentRowTime = currentRowTime.Add(
				time.Second * time.Duration(queryConfig.FillInterval),
			)
		}

		err := addValueToColumns(idx, gapFilledColumns, columns, len(*timeColumn))
		if err != nil {
			return err
		}
		currentRowTime = *timeCell
	}

	for idx := range columns {
		columns[idx] = gapFilledColumns[idx]
	}
	return nil
}

func addGapToColumns(
	fillTime time.Time,
	newColumns []*sqlColumn,
	originalColumns []*sqlColumn,
	queryConfig queryConfigStruct,
	originalRowCount int,
) error {
	for columnIndex := range originalColumns {
		if columnIndex == queryConfig.FillValuesTimeColumnIndex {
			newColumns[columnIndex].TimeData = append(newColumns[columnIndex].TimeData, &fillTime)
			continue
		}

		switch originalRowCount {
		case len(originalColumns[columnIndex].StringData):
			newColumns[columnIndex].StringData = append(newColumns[columnIndex].StringData, nil)
		case len(originalColumns[columnIndex].FloatData):
			newColumns[columnIndex].FloatData = append(newColumns[columnIndex].FloatData, nil)
		case len(originalColumns[columnIndex].IntData):
			newColumns[columnIndex].IntData = append(newColumns[columnIndex].IntData, nil)
		case len(originalColumns[columnIndex].TimeData):
			newColumns[columnIndex].TimeData = append(newColumns[columnIndex].TimeData, nil)
		default:
			log.DefaultLogger.Error(
				"could not find column type to fill gap for", "rowCount", originalRowCount,
			)
			return fmt.Errorf("error filling gaps: undetermined gap row type")
		}

	}
	return nil
}

func addValueToColumns(
	rowIndex int, newColumns []*sqlColumn, originalColumns []*sqlColumn, originalRowCount int,
) error {
	for columnIndex := range originalColumns {
		switch originalRowCount {
		case len(originalColumns[columnIndex].StringData):
			newColumns[columnIndex].StringData = append(
				newColumns[columnIndex].StringData,
				originalColumns[columnIndex].StringData[rowIndex],
			)
		case len(originalColumns[columnIndex].FloatData):
			newColumns[columnIndex].FloatData = append(
				newColumns[columnIndex].FloatData,
				originalColumns[columnIndex].FloatData[rowIndex],
			)
		case len(originalColumns[columnIndex].IntData):
			newColumns[columnIndex].IntData = append(
				newColumns[columnIndex].IntData,
				originalColumns[columnIndex].IntData[rowIndex],
			)
		case len(originalColumns[columnIndex].TimeData):
			newColumns[columnIndex].TimeData = append(
				newColumns[columnIndex].TimeData,
				originalColumns[columnIndex].TimeData[rowIndex],
			)
		default:
			log.DefaultLogger.Error(
				"could not find column type to set value (gap filling) for",
				"rowCount",
				originalRowCount,
			)
			return fmt.Errorf("error filling gaps: undetermined value row type")
		}

	}
	return nil
}
