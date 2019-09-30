package camera

import (
	"os"
	"os/exec"
)

// CameraCmder produce command for stream and image
type CameraCmder interface {
	GetCmdProducers(cc CameraConfig) []func() *exec.Cmd
}

// Camera stores config and provide REST interface
type Camera interface {
	Enable()
	Disable(wait bool)
	Refresh()
	CapturePhotoJPG() (*os.File, error)
	GetCapturePath() string
	GetVideoPaths() []string
	Configure(CameraConfig)
	MergeConfigure(CameraConfig)
	GetConfigure() CameraConfig
}

// support actions in: https://www.onvif.org/ver20/ptz/wsdl/ptz.wsdl
type PTZer interface {
}
