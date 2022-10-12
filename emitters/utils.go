package emitters

import (
	"strings"
)

func getNickNameFromUUIDs(deviceUUIDList string, ID string) string {
	deviceUuids := strings.Split(deviceUUIDList, ",")
	nickname := make([]string, 1)
	for _, uuidkv := range deviceUuids {
		uuidMap := strings.Split(uuidkv, "=")
		//panic(fmt.Sprintf("%+v -> %+v", uuidMap[0], ID))
		if uuidMap[0] == ID {
			nickname[0] = uuidMap[1]
		}
	}
	if nickname[0] == "" {
		nickname[0] = ID
	}
	return nickname[0]
}
