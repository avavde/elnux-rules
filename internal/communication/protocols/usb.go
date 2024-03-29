package communication

import (
	"github.com/google/gousb"
)

// USBProtocol описывает интерфейс для работы с USB.
type USBProtocol interface {
	FindDevices(vid, pid uint16) ([]*gousb.Device, error)                              // Поиск устройств по VID и PID
	Read(device *gousb.Device, endpoint uint8, buf []byte, timeout int) (int, error)   // Чтение данных из устройства
	Write(device *gousb.Device, endpoint uint8, data []byte, timeout int) (int, error) // Запись данных в устройство
}

func FindUSBDevices(vid, pid uint16) ([]*gousb.Device, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Находим все USB устройства, которые соответствуют заданным VID и PID
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == gousb.ID(vid) && desc.Product == gousb.ID(pid)
	})
	if err != nil {
		return nil, err
	}

	return devices, nil
}
