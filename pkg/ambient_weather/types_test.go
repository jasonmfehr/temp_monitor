package ambient_weather

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCloud(t *testing.T) {
	testDeviceID := "test-device"
	testApiKey := "test-api-key"
	testAppKey := "test-application-key"

	actual := NewCloud(testDeviceID, testApiKey, testAppKey)
	assert.NotNil(t, actual)
	assert.Equal(t, testDeviceID, actual.deviceID)
	assert.Equal(t, testApiKey, actual.apiKey)
	assert.Equal(t, testAppKey, actual.applicationKey)
}

type mockReaderCloser struct {
	mock.Mock
}

func (m *mockReaderCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReaderCloser) Close() error {
	return m.Called().Error(0)
}
