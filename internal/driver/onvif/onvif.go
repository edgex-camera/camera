package onvif

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/yakovlevdmv/goonvif"
	"github.com/yakovlevdmv/goonvif/PTZ"
	"github.com/yakovlevdmv/goonvif/xsd"
	"github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type onvifDevice interface {
	CallMethod(method interface{}) (*http.Response, error)
}

type onvifCamera struct {
	device    onvifDevice
	lc        logger.LoggingClient
	address   string
	stopTimer *time.Timer
}

func NewOnvif(lc logger.LoggingClient, address string) (cam Onvif, err error) {
	return &onvifCamera{
		lc:      lc,
		device:  nil,
		address: address,
	}, nil
}

func (c *onvifCamera) connect() (err error) {
	defer func() {
		if r := recover(); r != nil {
			c.lc.Error(fmt.Sprint("Init Onvif camera failed, Recovered in ", r))
			err = fmt.Errorf("Init Onvif camera failed")
		}
	}()

	if c.device != nil { // already connected
		return nil
	}
	device, err := goonvif.NewDevice(c.address)
	if err != nil {
		c.lc.Error("onvif camera connect error: %v", err)
		return err
	}
	device.Authenticate("admin", "admin")
	c.device = device
	return nil
}

func (c *onvifCamera) callMethod(method interface{}) error {
	err := c.connect()
	if err != nil {
		return err
	}
	_, err = c.device.CallMethod(method)
	if err != nil {
		return err
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(response.Body)
	// c.lc.Info(fmt.Sprintf("onvif callMethod response: %s", buf.String()))
	return nil
}

func (c *onvifCamera) ContinuousMove(timeout time.Duration, moveSpeed Move) error {
	c.lc.Info("camera move started")
	req := PTZ.ContinuousMove{
		ProfileToken: "Profile_000",
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{
				X:     moveSpeed.PanTiltSpeed.X,
				Y:     moveSpeed.PanTiltSpeed.Y,
				Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/PanTiltSpaces/GenericSpeedSpace"),
			},
			Zoom: onvif.Vector1D{
				X:     moveSpeed.Zoom,
				Space: xsd.AnyURI("http://www.onvif.org/ver10/tptz/ZoomSpaces/ZoomGenericSpeedSpace"),
			},
		},
		// timeout not working
	}

	if c.stopTimer != nil {
		c.stopTimer.Stop()
		c.stopTimer = nil
	}
	err := c.callMethod(req)
	c.stopTimer = time.AfterFunc(timeout, func() { _ = c.Stop() })
	return err
}

func (c *onvifCamera) Stop() error {
	if c.stopTimer != nil {
		c.stopTimer.Stop()
		c.stopTimer = nil
	}

	c.lc.Info("camera move stopped")
	req := PTZ.Stop{
		ProfileToken: "Profile_000",
		PanTilt:      true,
		Zoom:         true,
	}

	return c.callMethod(req)
}

func (c *onvifCamera) SetHomePosition() error {
	c.lc.Info("camera move reset")
	req := PTZ.SetPreset{
		ProfileToken: "Profile_000",
		PresetToken:  "1",
	}
	return c.callMethod(req)
}

func (c *onvifCamera) Reset() error {
	c.lc.Info("camera move reset")

	if c.stopTimer != nil {
		c.stopTimer.Stop()
		_ = c.Stop()
	}

	req := PTZ.GotoPreset{
		ProfileToken: "Profile_000",
		PresetToken:  "1",
	}
	return c.callMethod(req)
}
