package camera

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/reiver/go-v4l2"
)

func EnsureUsbExists(address string) error {
	if strings.HasPrefix(address, "rtsp") {
		return nil
	}

	device, err := v4l2.Open(address)
	if nil != err {
		return err
	}
	defer device.Close()

	busInfo, err := device.BusInfo()
	if err != nil {
		return err
	}

	busInfos := strings.Split(busInfo, "-")
	valuePath := fmt.Sprintf("/sys/devices/platform/%v/usb5/5-1/5-%v/bConfigurationValue", busInfos[1], busInfos[2])
	_, err = os.Stat(valuePath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(valuePath, []byte("1"), 0644)
}
