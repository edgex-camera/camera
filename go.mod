module github.com/edgex-camera/camera

go 1.13

require (
	github.com/beevik/etree v1.1.0 // indirect
	github.com/edgex-camera/device-sdk-go v1.0.1-0.20200329120140-a908c998ccc9

	github.com/edgex-camera/edgex-utils v0.0.0-20200328191541-b71e1e0b64ec
	github.com/edgexfoundry/go-mod-core-contracts v0.1.52
	github.com/elgs/gostrgen v0.0.0-20161222160715-9d61ae07eeae // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/pkg/errors v0.9.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/yakovlevdmv/Golang-iso8601-duration v0.0.0-20180403125811-e5db0413b903 // indirect
	github.com/yakovlevdmv/WS-Discovery v0.0.0-20180512141937-16170c6c3677 // indirect
	github.com/yakovlevdmv/goonvif v0.0.0-20180517145634-8181eb3ef2fb
	github.com/yakovlevdmv/gosoap v0.0.0-20180512142237-299a954b1c6d // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	gopkg.in/evanphx/json-patch.v4 v4.5.0

)

replace github.com/edgexfoundry/go-mod-core-contracts v0.1.52 => github.com/edgexfoundry/go-mod-core-contracts v0.1.14

replace github.com/edgexfoundry/go-mod-registry => github.com/edgexfoundry/go-mod-registry v0.1.9

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
