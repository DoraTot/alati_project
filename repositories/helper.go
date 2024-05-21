package repositories

import "fmt"

const (
	configs      = "configs/%s/v%.1f"
	configGroups = "configGroups/%s/v%.1f"
)

func constructKey(name string, version float32) string {
	return fmt.Sprintf(configs, name, version)
}

func constructKeyForGroup(name string, version float32) string {
	return fmt.Sprintf(configGroups, name, version)
}
