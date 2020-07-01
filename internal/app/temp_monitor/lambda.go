package temp_monitor

import (
	"context"
	"temp_monitor/internal/pkg/ambient_weather"
	"temp_monitor/internal/pkg/notify"
	"temp_monitor/internal/pkg/query"
	"temp_monitor/pkg/iface"

	"github.com/pkg/errors"

	"github.com/kelseyhightower/envconfig"
)

// ambientweatherNew enable mocking ambient_weather.New
var ambientweatherNew = ambient_weather.New

// funcQuerySorted enable mocking query.QuerySorted
var funcQuerySorted = query.QuerySorted

// funcNotify enable mocking notify.Do
var funcNotify = notify.Do

func HandleRequest(ctx context.Context) error {
	var env EnvConfig
	var weatherCloud iface.WeatherCloud

	err := envconfig.Process("", &env)
	if err != nil {
		return errors.WithStack(err)
	}

	switch env.SourceWeatherCloud {
	case string(AmbientWeather):
		weatherCloud, err = ambientweatherNew(env.FunctionID)
		if err != nil {
			return errors.Wrapf(err, "could not instantiate Ambient Weather cloud for function id '%s'", env.FunctionID)
		}
	default:
		return errors.Errorf("invalid source weather cloud: '%s'", env.SourceWeatherCloud)
	}

	// retrieve data and sort it
	dataPoints, err := funcQuerySorted(weatherCloud)
	if err != nil {
		return err
	}

	// determine if temperatures have crossed and notify accordingly
	return funcNotify(dataPoints, env.SNSTopicARN)
}
