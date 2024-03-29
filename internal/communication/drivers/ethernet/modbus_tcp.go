package communication

import (
	"log"
	"time"

	"github.com/goburrow/modbus"
)

// ModbusTCPDriver представляет собой драйвер для Modbus TCP.
type ModbusTCPDriver struct {
	handler *modbus.TCPClientHandler
	client  modbus.Client
}

// NewModbusTCPDriver создает и инициализирует новый экземпляр ModbusTCPDriver.
func NewModbusTCPDriver(address string) *ModbusTCPDriver {
	handler := modbus.NewTCPClientHandler(address)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1

	return &ModbusTCPDriver{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

// Connect открывает соединение с Modbus TCP сервером.
func (d *ModbusTCPDriver) Connect() error {
	if err := d.handler.Connect(); err != nil {
		log.Printf("Не удалось подключиться к Modbus TCP серверу: %v", err)
		return err
	}
	log.Println("Подключение к Modbus TCP серверу успешно установлено.")
	return nil
}

// Disconnect закрывает соединение с Modbus TCP сервером.
func (d *ModbusTCPDriver) Disconnect() {
	d.handler.Close()
	log.Println("Соединение с Modbus TCP сервером закрыто.")
}

// ReadHoldingRegisters читает holding регистры с сервера Modbus.
func (d *ModbusTCPDriver) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	results, err := d.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		log.Printf("Ошибка чтения holding регистров: %v", err)
		return nil, err
	}
	return results, nil
}

// WriteSingleRegister записывает значение в один holding регистр на сервере Modbus.
func (d *ModbusTCPDriver) WriteSingleRegister(address, value uint16) error {
	_, err := d.client.WriteSingleRegister(address, value)
	if err != nil {
		log.Printf("Ошибка записи в holding регистр: %v", err)
		return err
	}
	return nil
}
