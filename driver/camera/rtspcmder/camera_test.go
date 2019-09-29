package rtspcmder

import (
	"testing"
	"time"

	"gitlab.jiangxingai.com/applications/edgex/device-service/camera/internal/driver/camera"
	"gitlab.jiangxingai.com/applications/edgex/device-service/camera/internal/utils"
	"gitlab.jiangxingai.com/applications/edgex/edgex-utils/logger"
)

var lc = logger.NewPrintClient()
var rtspAddr = "rtsp://10.54.128.132:554/mpeg4"

func TestRtspStream(t *testing.T) {
	cc := camera.CameraConfig{
		StreamConfig: &camera.StreamConfig{
			Enabled: utils.BoolPointer(true),
			Address: "rtmp://10.54.128.94/live/test2",
		},
	}
	cmder := NewCmder(rtspAddr)
	cam := camera.NewCamera(lc, cmder)
	cam.Configure(cc)
	time.Sleep(10 * time.Second)
	cam.Disable(true)
}

func TestRtspCapture(t *testing.T) {
	cc := camera.CameraConfig{
		CaptureConfig: &camera.CaptureConfig{
			Enabled: utils.BoolPointer(true),
			Path:    "./capture.jpg",
		},
	}
	cmder := NewCmder(rtspAddr)
	cam := camera.NewCamera(lc, cmder)
	cam.Configure(cc)
	time.Sleep(10 * time.Second)
	cam.Disable(true)
}
