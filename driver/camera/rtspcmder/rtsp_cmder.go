package rtspcmder

import (
	"fmt"
	"os/exec"
	"strings"

	"gitlab.jiangxingai.com/applications/edgex/device-service/camera/internal/driver/camera"
)

type cmder struct {
	inputAddress string
}

const (
	FFmpegBaseCmdStr    = "-rtsp_transport tcp -i %s "
	FFmpegStreamCmdStr  = "-vcodec copy -an -f flv %s "
	FFmpegCaptureCmdStr = "-vcodec mjpeg -update 1 -y %s "
	FFmpegVideoCmdStr   = `-flags +global_header 
	-f stream_segment -segment_time %d -segment_format_options 
	movflags=+faststart -reset_timestamps 1 
	-vcodec copy -q:v 4 -an -r 24 -strftime 1 %s `

	GstBaseCmdStr    = "-e --gst-debug-level=3 rtspsrc location=%s ! rtph264depay ! h264parse ! tee name=t "
	GstStreamCmdStr  = "t. ! queue ! flvmux streamable=true ! rtmpsink sync=false location=%s "
	GstCaptureCmdStr = "t. ! queue ! avdec_h264 ! queue flush-on-eos=true ! jpegenc ! multifilesink post-messages=true location=%s max-files=1 "
	GstVideoCmdStr   = "t. ! queue ! splitmuxsink location=%s max-size-time=%d "
)

func NewCmder(inputAddress string) camera.CameraCmder {
	return &cmder{
		inputAddress: inputAddress,
	}
}

func checkGstAvail() bool {
	_, err := exec.LookPath("gst-launch-1.0")
	return err == nil
}

func checkFFmpegAvail() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

func (c *cmder) GetCmdProducers(cc camera.CameraConfig) []func() *exec.Cmd {
	if checkGstAvail() {
		return c.GetGstCmdProducers(cc)
	}

	if checkFFmpegAvail() {
		return c.GetFFmpegCmdProducers(cc)
	}

	return nil
}

func (c *cmder) GetFFmpegCmdProducers(cc camera.CameraConfig) []func() *exec.Cmd {
	producers := []func() *exec.Cmd{}
	if !*cc.StreamConfig.Enabled && !*cc.CaptureConfig.Enabled {
		return producers
	}

	func1 := func() *exec.Cmd {
		cmdStr := fmt.Sprintf(FFmpegBaseCmdStr, c.inputAddress)
		if *cc.StreamConfig.Enabled {
			cmdStr += fmt.Sprintf(FFmpegStreamCmdStr, cc.StreamConfig.Address)
		}
		if *cc.CaptureConfig.Enabled {
			cmdStr += fmt.Sprintf(FFmpegCaptureCmdStr, cc.CaptureConfig.Path)
		}
		if *cc.VideoConfig.Enabled {
			cmdStr += fmt.Sprintf(FFmpegVideoCmdStr, *cc.VideoConfig.Length, cc.VideoConfig.Path)
		}
		return exec.Command("ffmpeg", strings.Fields(cmdStr)...)
	}
	return []func() *exec.Cmd{func1}
}

func (c *cmder) GetGstCmdProducers(cc camera.CameraConfig) []func() *exec.Cmd {
	producers := []func() *exec.Cmd{}
	if !*cc.StreamConfig.Enabled && !*cc.CaptureConfig.Enabled {
		return producers
	}

	func1 := func() *exec.Cmd {
		cmdStr := fmt.Sprintf(GstBaseCmdStr, c.inputAddress)
		if *cc.StreamConfig.Enabled {
			cmdStr += fmt.Sprintf(GstStreamCmdStr, cc.StreamConfig.Address)
		}
		if *cc.CaptureConfig.Enabled {
			cmdStr += fmt.Sprintf(GstCaptureCmdStr, cc.CaptureConfig.Path)
		}
		if *cc.VideoConfig.Enabled {
			cmdStr += fmt.Sprintf(GstVideoCmdStr, cc.VideoConfig.Path, *cc.VideoConfig.Length*1000000000)
		}
		return exec.Command("gst-launch-1.0", strings.Fields(cmdStr)...)
	}
	return []func() *exec.Cmd{func1}
}
