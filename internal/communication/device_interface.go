package communication

// DeviceInterface определяет общий интерфейс для всех устройств.
type DeviceInterface interface {
	Connect() error          // Устанавливает соединение с устройством.
	Disconnect() error       // Закрывает соединение с устройством.
	Read() ([]byte, error)   // Читает данные с устройства.
	Write(data []byte) error // Записывает данные на устройство.
}
