package query

import (
	"temp_monitor/pkg/iface"
	"temp_monitor/pkg/iface_test"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHappyPath(t *testing.T) {
	mockCloud := &iface_test.MockWeatherCloud{}
	newestDataPoint := &iface.WeatherDataPoint{TimeStamp: time.Unix(1, 0)}
	oldestDataPoint := &iface.WeatherDataPoint{TimeStamp: time.Unix(0, 0)}

	mockCloud.On("GetMostRecentData", 2).Return([]*iface.WeatherDataPoint{
		newestDataPoint,
		oldestDataPoint,
	}, nil)

	actual, err := QuerySorted(mockCloud)
	assert.NoError(t, err)
	assert.Len(t, actual, 2)
	assert.Same(t, oldestDataPoint, actual[0])
	assert.Same(t, newestDataPoint, actual[1])

	mockCloud.AssertExpectations(t)
}

func TestQueryError(t *testing.T) {
	testErr := "test-error"
	mockCloud := &iface_test.MockWeatherCloud{}

	mockCloud.On("GetMostRecentData", mock.Anything).Return(nil, errors.New(testErr))

	actual, err := QuerySorted(mockCloud)

	assert.Nil(t, actual)
	assert.EqualError(t, err, testErr)
	mockCloud.AssertExpectations(t)
}

func TestTooFewDataPoints(t *testing.T) {
	mockCloud := &iface_test.MockWeatherCloud{}

	mockCloud.On("GetMostRecentData", mock.Anything).Return([]*iface.WeatherDataPoint{}, nil)

	actual, err := QuerySorted(mockCloud)
	assert.Nil(t, actual)
	assert.EqualError(t, err, "expected '2' data points but '0' were returned")

	mockCloud.AssertExpectations(t)
}
