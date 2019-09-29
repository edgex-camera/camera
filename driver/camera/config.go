package camera

import "gitlab.jiangxingai.com/applications/edgex/device-service/camera/internal/utils"

type StreamConfig struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Address string `json:"address,omitempty"`
}

type CaptureConfig struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Path    string `json:"path,omitempty"`
}

type VideoConfig struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	Path       string `json:"directory,omitempty"`
	Length     *int   `json:"length,omitempty"`
	KeepRecord *int   `json:"keep,omitempty"`
}

type CameraConfig struct {
	Height        *int
	Width         *int
	StreamConfig  *StreamConfig  `json:"stream,omitempty"`
	CaptureConfig *CaptureConfig `json:"capture,omitempty"`
	VideoConfig   *VideoConfig   `json:"video,omitempty"`
}

func DefaultCameraConfig() CameraConfig {
	length := 60 * 10
	return CameraConfig{
		StreamConfig: &StreamConfig{
			Enabled: utils.BoolPointer(false),
		},
		CaptureConfig: &CaptureConfig{
			Enabled: utils.BoolPointer(false),
			Path:    "capture.jpg",
		},
		VideoConfig: &VideoConfig{
			Enabled: utils.BoolPointer(false),
			Length:  &length,
			Path:    "./file-%05d.mp4",
		},
	}
}

// TODO: reuse merge config
func MergeConfigure(old *CameraConfig, new *CameraConfig) (updated bool) {
	updated = MergeStreamConfig(old.StreamConfig, new.StreamConfig) || updated
	updated = MergeCaptureConfig(old.CaptureConfig, new.CaptureConfig) || updated
	updated = MergeVideoConfig(old.VideoConfig, new.VideoConfig) || updated
	return updated
}

func MergeStreamConfig(dst *StreamConfig, src *StreamConfig) (updated bool) {
	updated = false
	if src == nil {
		return false
	}
	if src.Enabled != nil {
		dst.Enabled = src.Enabled
		updated = true
	}
	if src.Address != "" {
		dst.Address = src.Address
		updated = true
	}

	return updated
}

func MergeCaptureConfig(dst *CaptureConfig, src *CaptureConfig) (updated bool) {
	updated = false
	if src == nil {
		return false
	}
	if src.Enabled != nil {
		dst.Enabled = src.Enabled
		updated = true
	}
	if src.Path != "" {
		dst.Path = src.Path
		updated = true
	}

	return updated
}

func MergeVideoConfig(dst *VideoConfig, src *VideoConfig) (updated bool) {
	updated = false
	if src == nil {
		return false
	}
	if src.Enabled != nil {
		dst.Enabled = src.Enabled
		updated = true
	}
	if src.Path != "" {
		dst.Path = src.Path
		updated = true
	}
	if src.Length != nil {
		dst.Length = src.Length
		updated = true
	}
	if src.KeepRecord != nil {
		dst.KeepRecord = src.KeepRecord
		updated = true
	}

	return updated
}
