package addon

import (
	"fmt"
)

type AdapterProxy struct {
	*Adapter
	addonManager    *AddonManagerProxy
	OnPairing       func(float64)
	OnCancelPairing func()
}

func NewAdapterProxy(id, name, packageName string) *AdapterProxy {
	manager := NewAddonManagerProxy(packageName)
	adp := &AdapterProxy{Adapter: NewAdapter(manager, id, name, packageName)}
	adp.addonManager = manager
	adp.addonManager.handleAdapterAdded(adp)

	return adp
}

func (adapter *AdapterProxy) AddDevice(device *Device) {
	device.AdapterId = adapter.ID
	adapter.HandleDeviceAdded(device)
}

func (adapter *AdapterProxy) pairing(timeout float64) {
	if adapter.OnPairing != nil {
		adapter.OnPairing(timeout)
	}
}

func (adapter *AdapterProxy) cancelPairing() {
	if adapter.OnCancelPairing != nil {
		adapter.OnCancelPairing()
	}
}

func (adapter *AdapterProxy) HandleDeviceSaved(devId string, dev interface{}) {
	fmt.Print("on devices saved on the gateway")
}

func (adapter *AdapterProxy) Run() {
}

func (adapter *AdapterProxy) Unload() {
	fmt.Printf("adapter unload, AdapterId:%v", adapter.ID)
}

func (adapter *AdapterProxy) getID() string {
	return adapter.ID
}

func (adapter *AdapterProxy) getPackageName() string {
	return adapter.PackageName
}

func (adapter *AdapterProxy) CloseProxy() {
	fmt.Print("do some thing while adapter close")
}

func (adapter *AdapterProxy) ProxyRunning() bool {
	//return adapter.manager.Running()
	return true
}

func (adapter *AdapterProxy) HandlerDevicePropertyChanged(property *Property) {
	adapter.addonManager.sendPropertyChangedNotification(property)
}
