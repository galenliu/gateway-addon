package addon

import "fmt"

type OnProertyValueUpdateFunc func(propName string, value interface{})

type DeviceProxy struct {
	*Device
	adx                  *AdapterProxy
	OnUpdataProertyValue OnProertyValueUpdateFunc
}

func NewDeviceProxy(id, name string) *DeviceProxy {
	deivce := &DeviceProxy{
		Device: NewDevice(id, name),
	}
	return deivce
}

func (proxy *DeviceProxy) GetPropertyByType(typ string) *Property {
	for _, prop := range proxy.Device.Properties {
		if prop.AtType == typ {
			return prop
		}
	}
	return nil
}

func (proxy *DeviceProxy) SetPin(pin PIN) error {
	fmt.Print("set pin: %v", pin)
	return nil
}

func (proxy *DeviceProxy) setProperty(propName string, value interface{}) {

}

func (proxy *DeviceProxy) ConnectedNotify(connencted bool) {
	proxy.adx.managerProxy.sendConnectedStateNotification(proxy.Device, connencted)
}

func (proxy *DeviceProxy) NotifyPropertyChanged(property *Property) {
	proxy.adx.managerProxy.sendPropertyChangedNotification(proxy.AdapterId, property)
}
