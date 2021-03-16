package devices

import (
	"addon"
	"addon/properties"
	"strconv"
	"strings"
)

type LightHandler interface {
	TurnOn()
	TurnOff()
	SetBrightness(brightness int)
}

type LightBulb struct {
	*addon.DeviceProxy
	On     *properties.OnOffProperty
	Bright *properties.BrightnessProperty
	Color  *properties.ColorProperty
}

func NewLightBulb(id, title string) *LightBulb {

	lightBulb := &LightBulb{}
	lightBulb.DeviceProxy = addon.NewDeviceProxy(id, title)
	lightBulb.On = properties.NewOnOffProperty()
	lightBulb.AddProperty(lightBulb.On.Property)
	lightBulb.AddTypes(addon.Light, addon.OnOffSwitch)
	lightBulb.OnUpdataProertyValue = lightBulb.propertyValueUpdate
	return lightBulb
}

func (light *LightBulb) TurnOn() {
	light.On.SetValue(true)
}

func (light *LightBulb) TurnOff() {
	light.On.SetValue(false)
}

func (light *LightBulb) Toggle() {
	if light.On.Value == true {
		light.TurnOff()
	} else {
		light.TurnOn()
	}
}

func (light *LightBulb) SetBrightness(brightness int) {
	if light.Bright == nil {
		return
	}
	if brightness == 0 && light.On.Value == true {
		light.TurnOff()
	} else if brightness > 0 && light.On.Value == false {
		light.TurnOn()
	}
	light.Bright.SetValue(brightness)
}

func (light *LightBulb) propertyValueUpdate(propName string, newValue interface{}) {

}

func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}
