package addon

import (
	"fmt"
)

type Manager interface {
	HandleDeviceRemoved(device *Device)
	HandleDeviceAdded(device *Device)
}

type Adapter struct {
	ID          string `json:"adapterId"`
	Name        string `json:"name"`
	PackageName string `json:"packageName"`
	manager     Manager
	Devices     map[string]*Device
}

func NewAdapter(manager Manager, adapterId, name, packageName string) *Adapter {
	adapter := &Adapter{}
	adapter.manager = manager
	adapter.PackageName = packageName
	adapter.Name = name
	adapter.ID = adapterId
	adapter.Devices = make(map[string]*Device, 10)
	return adapter
}

func (adapter *Adapter) HandleDeviceAdded(device *Device) {
	if device == nil {
		fmt.Sprint("device addon invaild")
		return
	}
	device.AdapterId = adapter.ID
	adapter.Devices[device.ID] = device
	adapter.manager.HandleDeviceAdded(device)
}

func (adapter *Adapter) HandleDeviceRemoved(device *Device) {
	delete(adapter.Devices, device.ID)
	adapter.manager.HandleDeviceAdded(device)
}

func (adapter *Adapter) GetDevice(deviceId string) *Device {
	device, ok := adapter.Devices[deviceId]
	if !ok {
		return nil
	}
	return device
}

func (adapter *Adapter) FindDevice(deviceId string) (*Device, error) {
	device, ok := adapter.Devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("devices id:(%s) invaild", deviceId)
	}
	return device, nil
}
