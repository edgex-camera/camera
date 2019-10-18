package camera

import (
	"encoding/json"

	jsonpatch "gopkg.in/evanphx/json-patch.v4"
)

type StreamConfig struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"-"`
}

type CaptureConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"-"`
}

type VideoConfig struct {
	Enabled    bool   `json:"enabled"`
	Path       string `json:"-"`
	Length     int    `json:"length"`
	KeepRecord int    `json:"keep"`
}

// 画质选项
type QualityConfig struct {
	ImageHeight int `json:"-"`
	ImageWidth  int `json:"-"`
	VideoHeight int `json:"-"`
	VideoWidth  int `json:"-"`
	FrameRate   int `json:"-"`
}

type CameraConfig struct {
	Enabled       bool   `json:"enabled"`
	InputAddr     string `json:"-"`
	StreamConfig  `json:"stream"`
	CaptureConfig `json:"capture"`
	VideoConfig   `json:"video"`
	QualityConfig `json:"quality"`
}

var defaultConf = CameraConfig{
	Enabled:       true,
	StreamConfig:  StreamConfig{Enabled: true},
	CaptureConfig: CaptureConfig{Enabled: true},
	VideoConfig: VideoConfig{
		Enabled:    true,
		Length:     600,
		KeepRecord: 3,
	},
}

var DefaultConf []byte

func init() {
	DefaultConf, _ = json.Marshal(defaultConf)
}

func (c *camera) MergeConfig(configPatch []byte) error {
	if c.IsEnabled() {
		c.Disable(true)
	}

	old, err := json.Marshal(c.CameraConfig)
	if err != nil {
		return err
	}
	new, err := jsonpatch.MergePatch(old, configPatch)
	if err != nil {
		return err
	}
	err = json.Unmarshal(new, &c.CameraConfig)
	if err != nil {
		return err
	}

	if c.CameraConfig.Enabled {
		c.Enable()
	}
	return nil
}

func (c *camera) GetConfigure() []byte {
	config, _ := json.Marshal(c.CameraConfig)
	return config
}
