package ambient_weather

import (
	"fmt"
	"temp_monitor/pkg/ambient_weather"
	"testing"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

func TestNewSuccess(t *testing.T) {
	newCloudCalledCount := 0
	secretSvc := &mockSecretSvc{}
	testCloud := &ambient_weather.Cloud{}
	testFuncID := "test-func"
	testDeviceID := "test-device"
	testApiKey := "test-api-key"
	testAppKey := "test-application-key"

	secretSvc.On("GetSecretValue", mock.MatchedBy(func(actual *secretsmanager.GetSecretValueInput) bool {
		return assert.Equal(t, fmt.Sprintf("/temp_monitor/%s/ambient_weather", testFuncID), *actual.SecretId)
	})).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(fmt.Sprintf(`{"deviceID":"%s","apiKey":"%s","applicationKey":"%s"}`, testDeviceID, testApiKey, testAppKey)),
	}, nil)

	newSvcFunc = func() secretsmanageriface.SecretsManagerAPI {
		return secretSvc
	}

	newCloudFunc = func(deviceID string, apiKey string, applicationKey string) *ambient_weather.Cloud {
		newCloudCalledCount++
		assert.Equal(t, testDeviceID, deviceID)
		assert.Equal(t, testApiKey, apiKey)
		assert.Equal(t, testAppKey, applicationKey)
		return testCloud
	}

	actual, err := New(testFuncID)

	assert.Same(t, actual, testCloud)
	assert.NoError(t, err)
	assert.Equal(t, 1, newCloudCalledCount)
	secretSvc.AssertExpectations(t)
}

func TestNewSecretError(t *testing.T) {
	secretSvc := &mockSecretSvc{}
	testError := "test-error"

	secretSvc.On("GetSecretValue", mock.Anything).Return(nil, errors.New(testError))

	newSvcFunc = func() secretsmanageriface.SecretsManagerAPI {
		return secretSvc
	}

	actual, err := New("")

	assert.Nil(t, actual)
	assert.EqualError(t, err, fmt.Sprintf("could not retrieve secret named '/temp_monitor//ambient_weather': %s", testError))

	secretSvc.AssertExpectations(t)
}

func TestNewSecretMissingDeviceID(t *testing.T) {
	secretSvc := &mockSecretSvc{}

	secretSvc.On("GetSecretValue", mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(`{}`),
	}, nil)

	newSvcFunc = func() secretsmanageriface.SecretsManagerAPI {
		return secretSvc
	}

	actual, err := New("")

	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not find field 'deviceID' in secret '/temp_monitor//ambient_weather'")
	secretSvc.AssertExpectations(t)
}

func TestNewSecretMissingApiKey(t *testing.T) {
	secretSvc := &mockSecretSvc{}

	secretSvc.On("GetSecretValue", mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(`{"deviceID":""}`),
	}, nil)

	newSvcFunc = func() secretsmanageriface.SecretsManagerAPI {
		return secretSvc
	}

	actual, err := New("")

	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not find field 'apiKey' in secret '/temp_monitor//ambient_weather'")
	secretSvc.AssertExpectations(t)
}

func TestNewSecretMissingAppKey(t *testing.T) {
	secretSvc := &mockSecretSvc{}

	secretSvc.On("GetSecretValue", mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(`{"deviceID":"","apiKey":""}`),
	}, nil)

	newSvcFunc = func() secretsmanageriface.SecretsManagerAPI {
		return secretSvc
	}

	actual, err := New("")

	assert.Nil(t, actual)
	assert.EqualError(t, err, "could not find field 'applicationKey' in secret '/temp_monitor//ambient_weather'")
	secretSvc.AssertExpectations(t)
}

func TestNewSvc(t *testing.T) {
	assert.NotNil(t, newSvc())
}
