package onvif

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/Media"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

func getToken(address string) onvif.ReferenceToken {
	device, _ := goonvif.NewDevice(address)
	device.Authenticate("admin", "admin")
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
