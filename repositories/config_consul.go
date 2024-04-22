package repositories

import "projekat/model"

type ConfigConsulRepository struct {
}

func (c ConfigConsulRepository) GetConfig(name string, version float32) (*model.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) AddConfig(config *model.Config) error { panic("implement me") }

func (c ConfigConsulRepository) DeleteConfig(name string, version float32) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) AddToConfigGroup(config *model.Config, groupName string) error {
	panic("implement me")
}

func (c ConfigConsulRepository) DeleteFromConfigGroup(config *model.Config, groupName string) error {
	panic("implement me")
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigConsulRepository() model.ConfigRepository {
	return ConfigConsulRepository{}
}
