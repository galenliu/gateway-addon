package addon

type DeviceProxy struct {
	*Device
	Properties map[string]*PropertyProxy
	AdapterId  string
}
