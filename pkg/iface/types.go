package iface

import (
	"time"
)

type WeatherCloud interface {
	// GetMostRecentData returns the most recent datapoints, any maximums are up to the implementers
	//                   if the returned error is not nil, then the contents of the returned slice is
	//                   up to the implementer
	GetMostRecentData(count int) ([]*WeatherDataPoint, error)
}

type WeatherDataPoint struct {
	// InternalTempF the internal temperature in Farenheit
	InternalTempF float64

	// InternalTempFeelsLikeF what the the internal temperature actually feels like in Farenheit
	InternalTempFeelsLikeF float64

	// ExternalTempF the external temperature in Farenheit
	ExternalTempF float64

	// ExternalTempFeelsLikeF what the the external temperature actually feels like in Farenheit
	ExternalTempFeelsLikeF float64

	// TimeStamp the date/time of the data point
	TimeStamp time.Time
}
