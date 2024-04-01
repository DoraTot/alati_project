package repositories

import "projekat/model"

type ConfigConsulRepository struct {
}

func (c ConfigConsulRepository) GetConfig(name string, version float32) (*model.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) AddConfig(config *model.Config) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) DeleteConfig(name string) error {
	//TODO implement me
	panic("implement me")
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigConsulRepository() model.ConfigRepository {
	return ConfigConsulRepository{}
}
