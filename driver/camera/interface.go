package camera

import (
	"os"
	"os/exec"
)

// Camera stores config and provide REST interface
type Camera interface {
	Enable()
	Disable(wait bool)
	IsEnabled() bool

	CapturePhotoJPG() (*os.File, error)
	GetCapturePath() string
	GetVideoPaths() []string

	// use Json config
	MergeConfig(configPatch []byte) error
	GetConfigure() []byte
}

// CameraCmder produce command for stream and image
type CameraCmder interface {
	GetCmdProducers(cc CameraConfig) []func() *exec.Cmd
}
