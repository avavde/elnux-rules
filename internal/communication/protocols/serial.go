package communication

import (
	"github.com/tarm/serial"
)

// SerialProtocol описывает интерфейс для работы с последовательными портами.
type SerialProtocol interface {
	Open(config *serial.Config) (*serial.Port, error)  // Открытие последовательного порта
	Close(port *serial.Port) error                     // Закрытие последовательного порта
	Read(port *serial.Port, buf []byte) (int, error)   // Чтение данных из порта
	Write(port *serial.Port, data []byte) (int, error) // Запись данных в порт
}
