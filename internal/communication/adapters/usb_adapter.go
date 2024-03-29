package adapters

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

// USBAdapter представляет собой адаптер для работы с USB устройствами.
type USBAdapter struct {
	ctx    *gousb.Context
	device *gousb.Device
}

// NewUSBAdapter создает новый экземпляр USBAdapter.
func NewUSBAdapter() *USBAdapter {
	return &USBAdapter{
		ctx: gousb.NewContext(),
	}
}

// Connect подключается к USB устройству с заданным VID и PID.
func (adapter *USBAdapter) Connect(vid, pid gousb.ID) error {
	// Поиск устройства с заданным VID и PID.
	dev, err := adapter.ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		return fmt.Errorf("не удалось открыть USB устройство: %v", err)
	}
	if dev == nil {
		return fmt.Errorf("USB устройство не найдено")
	}

	adapter.device = dev
	log.Printf("Подключено к USB устройству: %v", dev)
	return nil
}

// Disconnect отключается от USB устройства.
func (adapter *USBAdapter) Disconnect() error {
	if adapter.device == nil {
		return fmt.Errorf("USB устройство не подключено")
	}

	if err := adapter.device.Close(); err != nil {
		return fmt.Errorf("не удалось закрыть USB устройство: %v", err)
	}

	adapter.device = nil
	log.Println("Отключение от USB устройства успешно")
	return nil
}

// Close освобождает ресурсы контекста USB.
func (adapter *USBAdapter) Close() {
	adapter.ctx.Close()
}
