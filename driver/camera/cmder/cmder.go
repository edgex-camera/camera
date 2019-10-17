package cmder

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"gitlab.jiangxingai.com/applications/edgex/device-service/camera/driver/camera"
)

type cmder struct {
	template              processCmdTemplate
	videoLengthMultiplier int
}

func NewCmder(processMethod string) camera.CameraCmder {
	switch {
	case processMethod == "gst-launch-1.0":
		return &cmder{
			template:              gstreamer,
			videoLengthMultiplier: 1000000000,
		}
	case processMethod == "ffmpeg":
		return &cmder{
			template:              ffmpeg,
			videoLengthMultiplier: 1,
		}
	case isGstAvail():
		return &cmder{
			template:              gstreamer,
			videoLengthMultiplier: 1000000000,
		}
	case isFFmpegAvail():
		return &cmder{
			template:              ffmpeg,
			videoLengthMultiplier: 1,
		}
	default:
		panic("no supported video processor")
	}
}

func isGstAvail() bool {
	_, err := exec.LookPath("gst-launch-1.0")
	return err == nil
}

func isFFmpegAvail() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func (c *cmder) GetCmdProducers(cc camera.CameraConfig) []func() *exec.Cmd {
	var template cmdTemplate
	switch {
	case strings.HasPrefix(cc.InputAddr, "rtsp://"):
		template = c.template.rtsp
	case strings.HasPrefix(cc.InputAddr, "/"):
		template = c.template.webcam
	default:
		log.Printf("input not supported: %s", cc.InputAddr)
		return []func() *exec.Cmd{}
	}

	outputCapture := cc.CaptureConfig.Enabled && cc.CaptureConfig.Path != ""
	outputStream := cc.StreamConfig.Enabled && cc.StreamConfig.Address != ""
	outputVideo := cc.VideoConfig.Enabled && cc.VideoConfig.Path != ""
	if !outputCapture && !outputStream && !outputVideo {
		return []func() *exec.Cmd{}
	}
	cmdStr := fmt.Sprintf(template.base, cc.Width, cc.Height, cc.Frame, cc.InputAddr)
	if outputCapture {
		cmdStr += fmt.Sprintf(template.capture, cc.CaptureConfig.Path)
	}
	if outputStream || outputVideo {
		cmdStr += template.h264
	}
	if outputStream {
		cmdStr += fmt.Sprintf(template.stream, cc.StreamConfig.Address)
	}
	if outputVideo {
		cmdStr += fmt.Sprintf(template.video, cc.VideoConfig.Length*c.videoLengthMultiplier, cc.VideoConfig.Path)
	}

	cmdList := strings.Fields(cmdStr)
	func1 := func() *exec.Cmd { return exec.Command(c.template.processor, cmdList...) }
	return []func() *exec.Cmd{func1}
}
