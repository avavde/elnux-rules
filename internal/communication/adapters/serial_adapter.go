package adapters

import (
	"github.com/tarm/serial"
)

// SerialAdapter представляет адаптер для работы с последовательным портом.
type SerialAdapter struct {
	config *serial.Config
	port   *serial.Port
}

// NewSerialAdapter создает новый экземпляр SerialAdapter.
func NewSerialAdapter(portName string, baudRate int) *SerialAdapter {
	return &SerialAdapter{
		config: &serial.Config{Name: portName, Baud: baudRate},
	}
}

// Connect устанавливает соединение с последовательным портом.
func (adapter *SerialAdapter) Connect() error {
	var err error
	adapter.port, err = serial.OpenPort(adapter.config)
	return err
}

// Disconnect закрывает соединение с последовательным портом.
func (adapter *SerialAdapter) Disconnect() error {
	return adapter.port.Close()
}

// Read читает данные из последовательного порта.
func (adapter *SerialAdapter) Read() ([]byte, error) {
	buf := make([]byte, 128) // Размер буфера
	n, err := adapter.port.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// Write записывает данные в последовательный порт.
func (adapter *SerialAdapter) Write(data []byte) error {
	_, err := adapter.port.Write(data)
	return err
}
