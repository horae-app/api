package Device

import (
	"encoding/json"
	"net/http"
)

func DeviceForm(r *http.Request) (Device, error) {
	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		return device, err
	}

	return device, device.Save()
}
