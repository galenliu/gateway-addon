package devices

import (
	"addon"
	"addon/properties"
	"strconv"
	"strings"
)

type LightBulb struct {
	*addon.Device
	On *properties.OnOffProperty
}

func NewLightBulb(id, title string) *LightBulb {
	lightBulb := &LightBulb{}
	lightBulb.Device = addon.NewDevice(id, title)
	lightBulb.On = properties.NewOnOffProperty()

	lightBulb.AddProperty(lightBulb.On.Property)
	lightBulb.AddTypes(addon.Light, addon.OnOffSwitch)

	return lightBulb
}

func Color16ToRGB(colorStr string) (red, green, blue int, err error) {
	color64, err := strconv.ParseInt(strings.TrimPrefix(colorStr, "#"), 16, 32)
	if err != nil {
		return
	}
	colorInt := int(color64)
	return colorInt >> 16, (colorInt & 0x00FF00) >> 8, colorInt & 0x0000FF, nil
}
