package repositories

import (
	"context"
	"github.com/stretchr/testify/mock"
	"projekat/model"
)

type MockConfigRepository struct {
	mock.Mock
}

func NewMockConfigRepository() *MockConfigRepository {
	return &MockConfigRepository{}
}

func (m *MockConfigRepository) AddConfig(config *model.Config, ctx context.Context) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockConfigRepository) DeleteConfig(name string, version float32, ctx context.Context) error {
	args := m.Called(ctx, name, version)
	return args.Error(0)
}

func (m *MockConfigRepository) GetConfig(name string, version float32, ctx context.Context) (*model.Config, error) {
	args := m.Called(name, version, ctx)
	return args.Get(0).(*model.Config), args.Error(1)
}

func (m *MockConfigRepository) GetConfigGroup(name string, version float32, ctx context.Context) (*model.ConfigGroup, error) {
	args := m.Called(name, version, ctx)
	return args.Get(0).(*model.ConfigGroup), args.Error(1)
}

func (m *MockConfigRepository) AddConfigGroup(config *model.ConfigGroup, ctx context.Context) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockConfigRepository) DeleteConfigGroup(name string, version float32, ctx context.Context) error {
	args := m.Called(ctx, name, version)
	return args.Error(0)
}
