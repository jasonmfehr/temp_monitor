package iface_test

import (
	"temp_monitor/pkg/iface"

	"github.com/stretchr/testify/mock"
)

type MockWeatherCloud struct {
	mock.Mock
}

func (m *MockWeatherCloud) GetMostRecentData(count int) ([]*iface.WeatherDataPoint, error) {
	args := m.Called(count)
	dataPoints := args.Get(0)
	err := args.Error(1)

	if dataPoints == nil {
		return nil, err
	}

	return dataPoints.([]*iface.WeatherDataPoint), err
}
