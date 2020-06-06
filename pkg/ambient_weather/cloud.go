package ambient_weather

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"temp_monitor/pkg/iface"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

// httpGet enable mocking http.Get
var httpGet = http.Get

// convertBodyFunc enable mocking convertBody
var convertBodyFunc = convertBody

func (c *Cloud) GetMostRecentData(count int) ([]*iface.WeatherDataPoint, error) {
	if count < 1 {
		return nil, errors.New("count must be greater than 0")
	}

	cloudURL := fmt.Sprintf("https://api.ambientweather.net/v1/devices/%s?apiKey=%s&applicationKey=%s&limit=%d", c.deviceID, c.apiKey, c.applicationKey, count)

	resp, err := httpGet(cloudURL)
	if err != nil {
		return nil, errors.Wrap(err, "could not get data")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return convertBodyFunc(body)
}

// convertBody converts the raw AmbientWeather response into
func convertBody(data []byte) ([]*iface.WeatherDataPoint, error) {
	returnData := []*iface.WeatherDataPoint{}
	var iterErr error

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, errParm error) {
		convertedDataPoint := &iface.WeatherDataPoint{}

		v, err := jsonparser.GetFloat(value, externalTempKey)
		if err != nil {
			iterErr = errors.Wrapf(err, "could not get field '%s'", externalTempKey)
			return
		}
		convertedDataPoint.ExternalTempF = v

		v, err = jsonparser.GetFloat(value, externalTempFeelsKey)
		if err != nil {
			iterErr = errors.Wrapf(err, "could not get field '%s'", externalTempFeelsKey)
			return
		}
		convertedDataPoint.ExternalTempFeelsLikeF = v

		v, err = jsonparser.GetFloat(value, internalTempKey)
		if err != nil {
			iterErr = errors.Wrapf(err, "could not get field '%s'", internalTempKey)
			return
		}
		convertedDataPoint.InternalTempF = v

		v, err = jsonparser.GetFloat(value, internalTempFeelsKey)
		if err != nil {
			iterErr = errors.Wrapf(err, "could not get field '%s'", internalTempFeelsKey)
			return
		}
		convertedDataPoint.InternalTempFeelsLikeF = v

		timeStamp, err := jsonparser.GetInt(value, dateKey)
		if err != nil {
			iterErr = errors.Wrapf(err, "could not get field '%s'", dateKey)
			return
		}
		convertedDataPoint.TimeStamp = time.Unix(timeStamp/1000, 0)

		returnData = append(returnData, convertedDataPoint)
	})

	if iterErr != nil {
		return nil, iterErr
	}

	if err != nil {
		return nil, errors.Wrap(err, "jsonparser.ArrayEach error")
	}

	return returnData, nil
}
