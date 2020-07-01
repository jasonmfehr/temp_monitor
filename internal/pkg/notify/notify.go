package notify

import (
	"fmt"
	"temp_monitor/pkg/iface"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/pkg/errors"
)

// newSnsFunc enable mocking newSns
var newSnsFunc = newSns

// Do determine if temperatures have crossed and notify accordingly
func Do(dataPoints []*iface.WeatherDataPoint, snsTopicARN string) error {
	var msg string

	older := dataPoints[0].ExternalTempFeelsLikeF > dataPoints[0].InternalTempFeelsLikeF
	newer := dataPoints[1].ExternalTempFeelsLikeF > dataPoints[1].InternalTempFeelsLikeF

	if older == true && newer == false {
		// outdoor temp is falling
		msg = fmt.Sprintf("Open the windows. Outdoor feels like %.1f, indoor feels like %.1f", dataPoints[1].ExternalTempFeelsLikeF, dataPoints[1].InternalTempFeelsLikeF)
	} else if older == false && newer == true {
		// outdoor temp is rising
		msg = fmt.Sprintf("Close the windows. Outdoor feels like %.1f, indoor feels like %.1f", dataPoints[1].ExternalTempFeelsLikeF, dataPoints[1].InternalTempFeelsLikeF)
	}

	if older != newer {
		//temperatures have crossed, send notification
		_, err := newSnsFunc().Publish(&sns.PublishInput{
			TopicArn: &snsTopicARN,
			Message:  &msg,
		})

		return errors.Wrapf(err, "could not publish message to SNS topic '%s'", snsTopicARN)
	}

	return nil
}

func newSns() snsiface.SNSAPI {
	return sns.New(session.New())
}
