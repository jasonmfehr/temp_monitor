package ambient_weather

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/mock"
)

type mockSecretSvc struct {
	mock.Mock
	secretsmanageriface.SecretsManagerAPI
}

func (m *mockSecretSvc) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	args := m.Called(input)

	output := args.Get(0)
	err := args.Error(1)

	if output == nil {
		return nil, err
	}

	return output.(*secretsmanager.GetSecretValueOutput), err
}
