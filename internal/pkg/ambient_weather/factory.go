package ambient_weather

import (
	"fmt"
	"temp_monitor/pkg/ambient_weather"
	"temp_monitor/pkg/iface"

	"github.com/buger/jsonparser"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// newSvcFunc enable mocking newSvc
var newSvcFunc = newSvc

// newCloudFunc enable mocking ambient_weather.NewCloud
var newCloudFunc = ambient_weather.NewCloud

// New configures a new implementation of an iface.WeatherCloud that points to Ambient Weather
//     three config items are needed: device id, api key, and application key,
//     attempts to read the necessary config from an AWS secret named
//     "/temp_monitor/{{funcID}}/ambient_weather"
func New(funcID string) (iface.WeatherCloud, error) {
	secretID := fmt.Sprintf("/temp_monitor/%s/ambient_weather", funcID)

	output, err := newSvcFunc().GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve secret named '%s'", secretID)
	}

	data := []byte(*output.SecretString)

	deviceID, err := jsonparser.GetString(data, deviceIDKey)
	if err != nil {
		return nil, errors.Errorf("could not find field '%s' in secret '%s'", deviceIDKey, secretID)
	}

	apiKey, err := jsonparser.GetString(data, apiKeyKey)
	if err != nil {
		return nil, errors.Errorf("could not find field '%s' in secret '%s'", apiKeyKey, secretID)
	}

	appKey, err := jsonparser.GetString(data, appKeyKey)
	if err != nil {
		return nil, errors.Errorf("could not find field '%s' in secret '%s'", appKeyKey, secretID)
	}

	return newCloudFunc(deviceID, apiKey, appKey), nil
}

func newSvc() secretsmanageriface.SecretsManagerAPI {
	return secretsmanager.New(session.New())
}
