package notify

import (
	"fmt"
	"temp_monitor/pkg/iface"
	"testing"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHappyPathOutdoorRising(t *testing.T) {
	oldestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 64.0,
		InternalTempFeelsLikeF: 65.0,
	}
	newestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 65.1,
		InternalTempFeelsLikeF: 65.0,
	}

	runHappyPathTest(t, "Close", oldestDataPoint, newestDataPoint)
}

func TestHappyPathOutdoorFalling(t *testing.T) {
	oldestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 65.1,
		InternalTempFeelsLikeF: 65.0,
	}
	newestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 64.0,
		InternalTempFeelsLikeF: 65.0,
	}

	runHappyPathTest(t, "Open", oldestDataPoint, newestDataPoint)
}

func runHappyPathTest(t *testing.T, expectedMsgPrefix string, oldestDataPoint *iface.WeatherDataPoint, newestDataPoint *iface.WeatherDataPoint) {
	testTopicARN := "test-topic-arn"

	mockSNSSvc := &mockSNS{}
	newSnsFunc = func() snsiface.SNSAPI {
		return mockSNSSvc
	}

	mockSNSSvc.On("Publish", mock.MatchedBy(func(actual *sns.PublishInput) bool {
		return assert.NotNil(t, actual) &&
			assert.Equal(t, testTopicARN, *actual.TopicArn) &&
			assert.Equal(t, fmt.Sprintf("%s the windows. Outdoor feels like %.1f, indoor feels like %.1f", expectedMsgPrefix, newestDataPoint.ExternalTempFeelsLikeF, newestDataPoint.InternalTempFeelsLikeF), *actual.Message)
	})).Return(&sns.PublishOutput{}, nil)

	assert.NoError(t, Do([]*iface.WeatherDataPoint{
		oldestDataPoint,
		newestDataPoint,
	}, testTopicARN))

	mockSNSSvc.AssertExpectations(t)
}

func TestSNSError(t *testing.T) {
	testErr := "test-error"
	oldestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 65.1,
		InternalTempFeelsLikeF: 65.0,
	}
	newestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 64.0,
		InternalTempFeelsLikeF: 65.0,
	}

	mockSNSSvc := &mockSNS{}
	newSnsFunc = func() snsiface.SNSAPI {
		return mockSNSSvc
	}

	mockSNSSvc.On("Publish", mock.Anything).Return(nil, errors.New(testErr))

	assert.EqualError(t, Do([]*iface.WeatherDataPoint{oldestDataPoint, newestDataPoint}, ""), fmt.Sprintf("could not publish message to SNS topic '': %s", testErr))

	mockSNSSvc.AssertExpectations(t)
}

func TestTempUnchanged(t *testing.T) {
	oldestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 65.1,
		InternalTempFeelsLikeF: 65.0,
	}
	newestDataPoint := &iface.WeatherDataPoint{
		ExternalTempFeelsLikeF: 65.1,
		InternalTempFeelsLikeF: 65.0,
	}

	mockSNSSvc := &mockSNS{}
	newSnsFunc = func() snsiface.SNSAPI {
		return mockSNSSvc
	}

	assert.Nil(t, Do([]*iface.WeatherDataPoint{oldestDataPoint, newestDataPoint}, ""))

	mockSNSSvc.AssertExpectations(t)
}

func TestNewSns(t *testing.T) {
	assert.NotNil(t, newSns())
}
