package notify

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/stretchr/testify/mock"
)

type mockSNS struct {
	mock.Mock
	snsiface.SNSAPI
}

func (m *mockSNS) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	args := m.Called(input)

	output := args.Get(0)
	err := args.Error(1)

	if output == nil {
		return nil, err
	}

	return output.(*sns.PublishOutput), err
}
