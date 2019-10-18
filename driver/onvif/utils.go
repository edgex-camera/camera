package onvif

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/Media"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type Envelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Body    Body
}

type Body struct {
	XMLName             xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
	GetProfilesResponse Response
}

type Response struct {
	XMLName  xml.Name
	Profiles Profiles
}

type Profiles struct {
	XMLName xml.Name
	Token   onvif.ReferenceToken `xml:"token,attr"`
}

func getToken(address string) onvif.ReferenceToken {
	device, _ := goonvif.NewDevice(address)
	device.Authenticate("admin", "admin")
	req := Media.GetProfiles{}
	res, _ := device.CallMethod(req)
	body, _ := ioutil.ReadAll(res.Body)

	response := &Envelope{}
	xml.Unmarshal(body, &response)
	return response.Body.GetProfilesResponse.Profiles.Token
}
