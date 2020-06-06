package ambient_weather

import (
	"fmt"
	"io"
	"net/http"
	"temp_monitor/pkg/iface"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMostRecentDataSuccess(t *testing.T) {
	httpGetCallCount := 0
	convertBodyFuncCallCount := 0
	testDeviceID := "test-device"
	testApiKey := "test-api-key"
	testAppKey := "test-application-key"
	testCount := 2
	mockBody := &mockReaderCloser{}
	testBodyContents := []byte(`{"foo":"bar"}`)
	testRet := []*iface.WeatherDataPoint{&iface.WeatherDataPoint{ExternalTempF: 49.8}}
	fixture := &Cloud{
		apiKey:         testApiKey,
		applicationKey: testAppKey,
		deviceID:       testDeviceID,
	}

	mockBody.On("Close").Return(nil)
	mockBody.On("Read", mock.MatchedBy(func(actual []byte) bool {
		for i, b := range testBodyContents {
			actual[i] = b
		}
		return true
	})).Return(len(testBodyContents), io.EOF)

	httpGet = func(url string) (resp *http.Response, err error) {
		httpGetCallCount++
		assert.Equal(t, fmt.Sprintf("https://api.ambientweather.net/v1/devices/%s?apiKey=%s&applicationKey=%s&limit=%d", testDeviceID, testApiKey, testAppKey, testCount), url)
		return &http.Response{Body: mockBody}, nil
	}

	convertBodyFunc = func(data []byte) ([]*iface.WeatherDataPoint, error) {
		convertBodyFuncCallCount++
		assert.Equal(t, testBodyContents, data)
		return testRet, nil
	}

	actual, err := fixture.GetMostRecentData(testCount)

	assert.Equal(t, testRet, actual)
	assert.NoError(t, err)
	assert.Equal(t, 1, httpGetCallCount)
	assert.Equal(t, 1, convertBodyFuncCallCount)
	mockBody.AssertExpectations(t)
}

func TestGetMostRecentDataInvalidCount(t *testing.T) {
	actual, err := (&Cloud{}).GetMostRecentData(0)
	assert.Nil(t, actual)
	assert.EqualError(t, err, "count must be greater than 0")
}

func TestGetMostRecentDataHttpError(t *testing.T) {
	httpGetCallCount := 0
	testError := "test error"

	httpGet = func(url string) (resp *http.Response, err error) {
		httpGetCallCount++
		return nil, errors.New(testError)
	}

	actual, err := (&Cloud{}).GetMostRecentData(1)
	assert.Nil(t, actual)
	assert.EqualError(t, err, fmt.Sprintf("could not get data: %s", testError))
	assert.Equal(t, 1, httpGetCallCount)
}

func TestGetMostRecentDataBodyReadError(t *testing.T) {
	testError := "test read error"
	mockBody := &mockReaderCloser{}

	mockBody.On("Close").Return(nil)
	mockBody.On("Read", mock.Anything).Return(0, errors.New(testError))

	httpGet = func(url string) (resp *http.Response, err error) {
		return &http.Response{Body: mockBody}, nil
	}

	actual, err := (&Cloud{}).GetMostRecentData(1)

	assert.Nil(t, actual)
	assert.EqualError(t, err, testError)
	mockBody.AssertExpectations(t)
}

func TestConvertBodySuccess(t *testing.T) {
	actual, err := convertBody([]byte(`[{"dateutc":1591394400000,"tempinf":71.4,"tempf":65.9,"solarradiation":163.93,"uv":1,"feelsLike":65.8,"feelsLikein":72.0,"date":"2020-06-05T22:00:00.000Z"},{"dateutc":1591393800000,"tempinf":69.6,"humidityin":52,"baromrelin":29.699,"baromabsin":29.374,"tempf":60.1,"feelsLike":60.4,"feelsLikein":68.7,"date":"2020-06-05T21:50:00.000Z"}]`))
	assert.Len(t, actual, 2)
	assert.Contains(t, actual, &iface.WeatherDataPoint{ExternalTempF: 65.9, ExternalTempFeelsLikeF: 65.8, InternalTempF: 71.4, InternalTempFeelsLikeF: 72.0, TimeStamp: time.Unix(1591394400, 0)})
	assert.Contains(t, actual, &iface.WeatherDataPoint{ExternalTempF: 60.1, ExternalTempFeelsLikeF: 60.4, InternalTempF: 69.6, InternalTempFeelsLikeF: 68.7, TimeStamp: time.Unix(1591393800, 0)})

	assert.NoError(t, err)
}

func TestConvertBodyExternalTempError(t *testing.T) {
	actual, err := convertBody([]byte(`[{"dateutc":1591394400000,"tempinf":71.4,"solarradiation":163.93,"uv":1,"feelsLike":65.8,"feelsLikein":72.0,"date":"2020-06-05T22:00:00.000Z"}]`))
	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not get field 'tempf': Key path not found")
}

func TestConvertBodyExternalTempFeelError(t *testing.T) {
	actual, err := convertBody([]byte(`[{"dateutc":1591394400000,"tempinf":71.4,"tempf":65.9,"solarradiation":163.93,"uv":1,"feelsLikein":72.0,"date":"2020-06-05T22:00:00.000Z"}]`))
	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not get field 'feelsLike': Key path not found")
}

func TestConvertBodyInternalTempError(t *testing.T) {
	actual, err := convertBody([]byte(`[{"dateutc":1591394400000,"tempf":65.9,"solarradiation":163.93,"uv":1,"feelsLike":65.8,"feelsLikein":72.0,"date":"2020-06-05T22:00:00.000Z"}]`))
	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not get field 'tempinf': Key path not found")
}

func TestConvertBodyInternalFeelTempError(t *testing.T) {
	actual, err := convertBody([]byte(`[{"dateutc":1591394400000,"tempinf":71.4,"tempf":65.9,"solarradiation":163.93,"uv":1,"feelsLike":65.8,"date":"2020-06-05T22:00:00.000Z"}]`))
	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not get field 'feelsLikein': Key path not found")
}

func TestConvertBodyTimestampError(t *testing.T) {
	actual, err := convertBody([]byte(`[{"tempinf":71.4,"tempf":65.9,"solarradiation":163.93,"uv":1,"feelsLike":65.8,"feelsLikein":72.0,"date":"2020-06-05T22:00:00.000Z"}]`))
	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not get field 'dateutc': Key path not found")
}

func TestConvertBodyArrayIterError(t *testing.T) {
	actual, err := convertBody([]byte(`{`))
	assert.Nil(t, actual)
	assert.EqualError(t, err, "jsonparser.ArrayEach error: Malformed JSON error")
}
