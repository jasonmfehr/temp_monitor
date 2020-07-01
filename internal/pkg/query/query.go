package query

import (
	"sort"
	"temp_monitor/pkg/iface"

	"github.com/pkg/errors"
)

// QuerySorted retrieves data from the weather cloud and sorts it ascending from oldest to newest
func QuerySorted(weatherCloud iface.WeatherCloud) ([]*iface.WeatherDataPoint, error) {
	dataPoints, err := weatherCloud.GetMostRecentData(dataPointsToRetrieve)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(dataPoints) < dataPointsToRetrieve {
		return nil, errors.Errorf("expected '%d' data points but '%d' were returned", dataPointsToRetrieve, len(dataPoints))
	}

	sort.Slice(dataPoints, func(i int, j int) bool {
		return dataPoints[i].TimeStamp.UnixNano() < dataPoints[j].TimeStamp.UnixNano()
	})

	return dataPoints, nil
}
