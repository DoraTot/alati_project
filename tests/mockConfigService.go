// tests/mock_config_service.go
package tests

import (
	"context"
	"github.com/stretchr/testify/mock"
	"projekat/model"
	"projekat/services"
)

type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) Hello() {
	m.Called()
}

func (m *MockConfigService) AddConfig(name string, version float32, parameters map[string]string, ctx context.Context) error {
	args := m.Called(name, version, parameters, ctx)
	return args.Error(0)
}

func (m *MockConfigService) GetConfig(name string, version float32, ctx context.Context) (*model.Config, error) {
	args := m.Called(name, version, ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Config), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockConfigService) DeleteConfig(name string, version float32, ctx context.Context) error {
	args := m.Called(name, version, ctx)
	return args.Error(0)
}

// Ensure MockConfigService implements services.ConfigServiceInterface
var _ services.ConfigServiceInterface = (*MockConfigService)(nil)
