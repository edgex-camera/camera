package webcamcmder

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
	FFmpegBaseCmdStr    = "-use_wallclock_as_timestamps 1 -f v4l2 -vcodec mjpeg -i %s "
	FFmpegStreamCmdStr  = "-vcodec h264 -an -f flv %s "
	FFmpegCaptureCmdStr = "-vcodec mjpeg -update 1 -y %s "
	FFmpegVideoCmdStr   = `-flags +global_header 
	-f stream_segment -segment_time %d -segment_format_options 
	movflags=+faststart -reset_timestamps 1 
	-vcodec h264 -q:v 4 -an -r 24 -strftime 1 %s `

	GstBaseCmdStr    = "-e --gst-debug-level=3 v4l2src device=%s ! image/jpeg,width=%d,height=%d,framerate=25/1 ! jpegdec ! queue ! videoconvert ! tee name=t "
	GstCaptureCmdStr = "t. ! queue flush-on-eos=true ! mppjpegenc ! multifilesink location=%s max-files=1 post-messages=true "
	GstH264CmdStr    = "t. ! queue ! mpph264enc vbr=false bitrate=\"800000\" filerate=false ! queue ! h264parse ! tee name=v "
	GstStreamCmdStr  = "v. ! queue ! flvmux streamable=true ! rtmpsink sync=false location=%s "
	GstVideoCmdStr   = "v. ! queue ! splitmuxsink location=%s max-size-time=%d "
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
		cmdStr := ""
		if cc.Height != nil && cc.Width != nil {
			cmdStr = fmt.Sprintf(GstBaseCmdStr, c.inputAddress, *cc.Width, *cc.Height)
		} else {
			cmdStr = fmt.Sprintf(GstBaseCmdStr, c.inputAddress, 640, 480)
		}
		if *cc.StreamConfig.Enabled {
			cmdStr += fmt.Sprintf(GstStreamCmdStr, cc.StreamConfig.Address)
		}
		if *cc.CaptureConfig.Enabled || *cc.VideoConfig.Enabled {
			cmdStr += GstH264CmdStr
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
