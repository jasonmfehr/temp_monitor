package temp_monitor

import (
	"fmt"
	"os"
	"temp_monitor/pkg/iface"
	"temp_monitor/pkg/iface_test"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

const testFuncID = "test-func"
const testSNSTopicARN = "topic-bar"

func TestHandleNotifyHappyPath(t *testing.T) {
	mock := &iface_test.MockWeatherCloud{}
	testDataPoints := []*iface.WeatherDataPoint{&iface.WeatherDataPoint{TimeStamp: time.Now()}}

	setupEnv()

	ambientweatherNew = func(funcID string) (iface.WeatherCloud, error) {
		assert.Equal(t, testFuncID, funcID)
		return mock, nil
	}

	funcQuerySorted = func(weatherCloud iface.WeatherCloud) ([]*iface.WeatherDataPoint, error) {
		return testDataPoints, nil
	}

	funcNotify = func(dataPoints []*iface.WeatherDataPoint, snsTopicARN string) error {
		assert.Equal(t, testDataPoints, dataPoints)
		assert.Equal(t, testSNSTopicARN, snsTopicARN)
		return nil
	}

	assert.Nil(t, HandleRequest(nil))
}

func TestHandleRequestInvalidENV(t *testing.T) {
	os.Unsetenv("SOURCE_WEATHER_CLOUD")
	err := HandleRequest(nil)

	assert.EqualError(t, err, "required key SOURCE_WEATHER_CLOUD missing value")
}

func TestHandleRequestCloudNotFound(t *testing.T) {
	setupEnv()
	os.Setenv("SOURCE_WEATHER_CLOUD", "idontexist")
	err := HandleRequest(nil)

	assert.EqualError(t, err, "invalid source weather cloud: 'idontexist'")
}

func TestHandleRequestErrorInitiateCloud(t *testing.T) {
	testErr := "test-error"
	setupEnv()

	ambientweatherNew = func(funcID string) (iface.WeatherCloud, error) {
		assert.Equal(t, testFuncID, funcID)
		return nil, errors.New(testErr)
	}

	err := HandleRequest(nil)

	assert.EqualError(t, err, fmt.Sprintf("could not instantiate Ambient Weather cloud for function id '%s': %s", testFuncID, testErr))
}

func TestHandleRequestQueryError(t *testing.T) {
	mock := &iface_test.MockWeatherCloud{}
	testErr := "test-error"

	setupEnv()

	ambientweatherNew = func(funcID string) (iface.WeatherCloud, error) {
		assert.Equal(t, testFuncID, funcID)
		return mock, nil
	}

	funcQuerySorted = func(weatherCloud iface.WeatherCloud) ([]*iface.WeatherDataPoint, error) {
		assert.Same(t, mock, weatherCloud)
		return nil, errors.New(testErr)
	}

	err := HandleRequest(nil)

	assert.EqualError(t, err, testErr)
}

func TestHandleNotifyError(t *testing.T) {
	mock := &iface_test.MockWeatherCloud{}
	testErr := "test-error"

	setupEnv()

	ambientweatherNew = func(funcID string) (iface.WeatherCloud, error) {
		assert.Equal(t, testFuncID, funcID)
		return mock, nil
	}

	funcQuerySorted = func(weatherCloud iface.WeatherCloud) ([]*iface.WeatherDataPoint, error) {
		return []*iface.WeatherDataPoint{}, nil
	}

	funcNotify = func(dataPoints []*iface.WeatherDataPoint, snsTopicARN string) error {
		assert.Equal(t, testSNSTopicARN, snsTopicARN)
		return errors.New(testErr)
	}

	err := HandleRequest(nil)

	assert.EqualError(t, err, testErr)
}

func setupEnv() {
	os.Setenv("SOURCE_WEATHER_CLOUD", "AmbientWeather")
	os.Setenv("FUNCTION_ID", testFuncID)
	os.Setenv("SNS_TOPIC_ARN", testSNSTopicARN)
}
