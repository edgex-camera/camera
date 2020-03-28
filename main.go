package main

import (
	"github.com/edgex-camera/camera/driver/camera/cmder"
	"github.com/edgex-camera/camera/driver/onvif"
)

func main() {
	_ = cmder.NewCmder("sss")
	_, _ = onvif.NewOnvif(nil, nil)
}
