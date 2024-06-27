package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"testing"
)

func TestAddConfig(t *testing.T) {
	configName := "db_config"
	configVersion := float32(5)
	configParameters := map[string]string{
		"additionalProp1": "param1",
		"additionalProp2": "param2",
		"additionalProp3": "param3",
	}

	mockRepo := new(repositories.MockConfigRepository)

	mockRepo.On("AddConfig", mock.Anything, mock.Anything).Return(nil)

	service := services.NewConfigService(mockRepo)

	err := service.AddConfig(configName, configVersion, configParameters, context.Background())

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetConfig(t *testing.T) {

	configName := "db_config"
	configVersion := float32(5)
	expectedConfig := &model.Config{
		Name:       configName,
		Version:    configVersion,
		Parameters: map[string]string{"param1": "value1", "param2": "value2"},
	}

	mockRepo := new(repositories.MockConfigRepository)

	mockRepo.On("GetConfig", configName, configVersion, mock.Anything).Return(expectedConfig, nil)

	service := services.NewConfigService(mockRepo)

	retrievedConfig, err := service.GetConfig(configName, configVersion, context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, retrievedConfig)

	mockRepo.AssertExpectations(t)
}

func TestDeleteConfig(t *testing.T) {
	configName := "db_config"
	configVersion := float32(5)
	mockRepo := new(repositories.MockConfigRepository)

	mockRepo.On("DeleteConfig", mock.Anything, configName, configVersion).Return(nil)
	service := services.NewConfigService(mockRepo)
	err := service.DeleteConfig(configName, configVersion, context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
