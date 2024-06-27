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

func TestGetConfigGroup(t *testing.T) {
	configGroupName := "db_config_group"
	configGroupVersion := float32(1.0)
	expectedConfigGroup := &model.ConfigGroup{
		Name:    configGroupName,
		Version: configGroupVersion,
		Configurations: []model.ConfigForGroup{
			{Name: "config1", Labels: map[string]string{"key2": "value2"}, Parameters: map[string]string{"key1": "value1"}},
			{Name: "config2", Labels: map[string]string{"key4": "value4"}, Parameters: map[string]string{"key2": "value2"}},
		},
	}

	mockRepo := new(repositories.MockConfigRepository)

	mockRepo.On("GetConfigGroup", configGroupName, configGroupVersion, mock.Anything).Return(expectedConfigGroup, nil)

	service := services.NewConfigGroupService(mockRepo)
	retrievedConfigGroup, err := service.GetConfigGroup(configGroupName, configGroupVersion, context.Background())

	assert.NoError(t, err)
	assert.Equal(t, retrievedConfigGroup, expectedConfigGroup)

	mockRepo.AssertExpectations(t)
}

func TestAddConfigGroup(t *testing.T) {
	configGroupName := "db_config_group"
	configGroupVersion := float32(1.0)
	expectedConfigurations := []model.ConfigForGroup{
		{Name: "config1", Labels: map[string]string{"key2": "value2"}, Parameters: map[string]string{"key1": "value1"}},
		{Name: "config2", Labels: map[string]string{"key4": "value4"}, Parameters: map[string]string{"key2": "value2"}},
	}

	mockRepo := new(repositories.MockConfigRepository)
	mockRepo.On("AddConfigGroup", mock.Anything, mock.Anything).Return(nil)
	service := services.NewConfigGroupService(mockRepo)
	err := service.AddConfigGroup(configGroupName, configGroupVersion, expectedConfigurations, context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteConfigGroup(t *testing.T) {
	configGroupName := "db_config_group"
	configGroupVersion := float32(1.0)

	mockRepo := new(repositories.MockConfigRepository)
	mockRepo.On("DeleteConfigGroup", mock.Anything, configGroupName, configGroupVersion).Return(nil)
	service := services.NewConfigGroupService(mockRepo)
	err := service.DeleteConfigGroup(configGroupName, configGroupVersion, context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
