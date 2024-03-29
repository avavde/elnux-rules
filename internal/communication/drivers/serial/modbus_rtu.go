package serial

import (
	"log"
	"time"

	"github.com/goburrow/modbus"
)

// ModbusRTUSerialDriver представляет собой драйвер для Modbus RTU через последовательный порт.
type ModbusRTUSerialDriver struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

// NewModbusRTUSerialDriver создает новый экземпляр ModbusRTUSerialDriver с указанными настройками порта.
func NewModbusRTUSerialDriver(device string, baudrate int, databits int, parity string, stopbits int) *ModbusRTUSerialDriver {
	handler := modbus.NewRTUClientHandler(device)
	handler.BaudRate = baudrate
	handler.DataBits = databits
	handler.Parity = parity
	handler.StopBits = stopbits
	handler.Timeout = 10 * time.Second

	return &ModbusRTUSerialDriver{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

// Connect подключается к Modbus RTU устройству через последовательный порт.
func (driver *ModbusRTUSerialDriver) Connect() error {
	log.Println("Connecting to Modbus RTU device via serial port")
	return driver.handler.Connect()
}

// Disconnect отключается от Modbus RTU устройства.
func (driver *ModbusRTUSerialDriver) Disconnect() {
	log.Println("Disconnecting from Modbus RTU device")
	driver.handler.Close()
}

// ReadHoldingRegisters читает holding регистры с Modbus RTU устройства.
func (driver *ModbusRTUSerialDriver) ReadHoldingRegisters(slaveID byte, address uint16, quantity uint16) ([]byte, error) {
	driver.handler.SlaveId = slaveID
	results, err := driver.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		log.Printf("Error reading holding registers: %v", err)
		return nil, err
	}
	return results, nil
}

// WriteSingleRegister записывает значение в один holding регистр на Modbus RTU устройстве.
func (driver *ModbusRTUSerialDriver) WriteSingleRegister(slaveID byte, address uint16, value uint16) error {
	driver.handler.SlaveId = slaveID
	results, err := driver.client.WriteSingleRegister(address, value)
	if err != nil {
		log.Printf("Error writing single register: %v", err)
		return err
	}
	log.Printf("Write single register response: %v", results)
	return nil
}
