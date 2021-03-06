package onvif

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strconv"

	"github.com/edgex-camera/device-sdk-go"
	"github.com/edgex-camera/device-sdk-go/pkg/camstartup"
	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/Media"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

func getToken(config OnvifConfig) onvif.ReferenceToken {
	device, _ := goonvif.NewDevice(config.Address)
	device.Authenticate(config.Username, config.Password)
	req := Media.GetProfiles{}
	res, _ := device.CallMethod(req)
	body, _ := ioutil.ReadAll(res.Body)

	var response struct {
		Body struct {
			GetProfilesResponse struct {
				Profiles []onvif.Profile
			}
		}
	}
	xml.Unmarshal(body, &response)

	// 取第一个Profile使用
	profile := response.Body.GetProfilesResponse.Profiles[0]
	return profile.Token
}

// 新建预置点配置，1点占用，2-255未占用
func InitPresetsConfig() {
	presets := make(map[int64]bool)
	presets[1] = true
	for i := 2; i < 256; i++ {
		presets[int64(i)] = false
	}
	config, _ := json.Marshal(presets)
	camstartup.PutDriverConfig("presets", config)
}

func getPresets() string {
	return device.DriverConfigs()["presets"]
}

func setPreset(number int64) {
	InitPresetsConfig()
	current := []byte(device.DriverConfigs()["presets"])
	current_map := make(map[int64]bool)
	json.Unmarshal(current, &current_map)

	current_map[number] = true
	new_presets, _ := json.Marshal(current_map)
	camstartup.PutDriverConfig("presets", new_presets)
}

func numberToToken(number int64) onvif.ReferenceToken {
	return onvif.ReferenceToken(strconv.Itoa(int(number)))
}
