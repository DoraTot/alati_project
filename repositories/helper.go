package repositories

import "fmt"

const (
	configs             = "configs/%s/v%.1f"
	configGroups        = "configGroups/%s/v%.1f"
	configsByLabels     = "configGroup/%s/v%.1f/%s"
	idempotencyRequests = "idempotency_requests/%s/"
)

func constructKey(name string, version float32) string {
	return fmt.Sprintf(configs, name, version)
}

func constructKeyConfigsByLabels(groupName string, groupVersion float32, labels map[string]string) string {
	labelsStr := ""
	for key, value := range labels {
		labelsStr += fmt.Sprintf("%s:%s/", key, value)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[:len(labelsStr)-1]
	}
	return fmt.Sprintf(configsByLabels, groupName, groupVersion, labelsStr)
}

func constructIdempotencyRequestKey(key string) string {
	return fmt.Sprintf(idempotencyRequests, key)
}

func constructKeyForGroup(name string, version float32) string {
	return fmt.Sprintf(configGroups, name, version)
}
