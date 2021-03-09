package addon

type DeviceHandler interface {
}

type DeviceProxy struct {
	*Device
	Handler    DeviceHandler
	Properties map[string]*PropertyProxy
	AdapterId  string
}

func NewDeviceProxy(id, name string) *DeviceProxy {
	deivce := &DeviceProxy{
		Device: NewDevice(id, name),
	}
	return deivce
}
