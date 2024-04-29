package plugin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

const macroRegex = `\$__([_a-zA-Z0-9]+)\(([^\)]*)\)`

func applyMacros(queryConfig *queryConfigStruct) error {
	compiledRegex, err := regexp.Compile(macroRegex)
	if err != nil {
		log.DefaultLogger.Error("Could create macro regex", "err", err)
		return err
	}

	newQuery := ""
	lastReplacedIndex := 0

	for _, match := range compiledRegex.FindAllSubmatchIndex([]byte(queryConfig.FinalQuery), -1) {
		groups := []string{}

		for i := 0; i < len(match); i += 2 {
			groups = append(groups, queryConfig.FinalQuery[match[i]:match[i+1]])
		}

		var replacedString string
		switch groups[1] {
		case "timeFilter":
			replacedString, err = timeFilter(queryConfig, strings.Split(groups[2], ","))
			if err != nil {
				return err
			}
		case "unixEpochGroupSeconds":
			replacedString, err = unixEpochGroupSeconds(queryConfig, strings.Split(groups[2], ","))
			if err != nil {
				return err
			}
		default:
			replacedString = groups[0]
		}

		newQuery += queryConfig.FinalQuery[lastReplacedIndex:match[0]] + replacedString
		lastReplacedIndex = match[1]
	}

	queryConfig.FinalQuery = newQuery + queryConfig.FinalQuery[lastReplacedIndex:]

	return nil
}

func timeFilter(queryConfig *queryConfigStruct, arguments []string) (string, error) {
	if len(arguments) != 1 {
		return "", fmt.Errorf(
			"unsupported number of arguments (%d) for timeFilter", len(arguments),
		)
	}

	return fmt.Sprintf(
		// TODO: Check if this is faster
		// `CASE
		// 	WHEN typeof(%[1]s) = 'integer' THEN %[1]s BETWEEN %[4]d AND %[5]d
		// 	ELSE %[1]s BETWEEN %[2]s AND %[3]s
		// END`,
		`CAST(CASE WHEN typeof(%[1]s) = 'integer' THEN %[1]s ELSE strftime("%%s", %[1]s) END AS INTEGER) BETWEEN %[2]d AND %[3]d`,
		arguments[0],
		queryConfig.TimeRange.From.Unix(),
		queryConfig.TimeRange.To.Unix(),
	), nil
}

func unixEpochGroupSeconds(queryConfig *queryConfigStruct, arguments []string) (string, error) {
	if len(arguments) < 2 || len(arguments) > 3 {
		return "", fmt.Errorf(
			"unsupported number of arguments (%d) for unixEpochGroupSeconds", len(arguments),
		)
	}
	var err error
	queryConfig.FillInterval, err = strconv.Atoi(strings.Trim(arguments[1], " "))
	if err != nil {
		log.DefaultLogger.Error(
			"Could not convert grouping interval to an integer",
			"macro",
			"unixEpochGroupSeconds",
			"err",
			err,
		)
		return "", fmt.Errorf(
			"could not convert '%s' to an integer grouping interval", arguments[1],
		)
	}

	// the gap filling value
	if len(arguments) == 3 {
		if strings.ToLower(strings.Trim(arguments[2], " ")) != "null" {
			return "", fmt.Errorf("unsupported gap filling value of: `%s`", arguments[2])
		}
		queryConfig.ShouldFillValues = true
	}

	return fmt.Sprintf(
		"cast((%s / %d) as int) * %d",
		arguments[0],
		queryConfig.FillInterval,
		queryConfig.FillInterval,
	), nil
}
