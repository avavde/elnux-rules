package drivers

import (
	"log"
	"time"

	"github.com/goburrow/modbus"
)

// ModbusRTUDriver представляет драйвер для взаимодействия с устройствами по протоколу Modbus RTU.
type ModbusRTUDriver struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
}

// NewModbusRTUDriver создает новый экземпляр ModbusRTUDriver для заданного последовательного порта.
func NewModbusRTUDriver(serialPort string, baudRate int, slaveID byte) *ModbusRTUDriver {
	handler := modbus.NewRTUClientHandler(serialPort)
	handler.BaudRate = baudRate
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = slaveID
	handler.Timeout = 10 * time.Second

	return &ModbusRTUDriver{
		handler: handler,
		client:  modbus.NewClient(handler),
	}
}

// Connect открывает соединение с последовательным портом для Modbus RTU.
func (d *ModbusRTUDriver) Connect() error {
	if err := d.handler.Connect(); err != nil {
		return err
	}
	log.Println("Modbus RTU соединение установлено")
	return nil
}

// Disconnect закрывает соединение с последовательным портом.
func (d *ModbusRTUDriver) Disconnect() error {
	d.handler.Close()
	log.Println("Modbus RTU соединение закрыто")
	return nil
}

// ReadHoldingRegisters читает holding регистры начиная с `address` длиной `quantity`.
func (d *ModbusRTUDriver) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	results, err := d.client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// WriteSingleRegister записывает значение в один holding регистр по адресу `address`.
func (d *ModbusRTUDriver) WriteSingleRegister(address, value uint16) error {
	_, err := d.client.WriteSingleRegister(address, value)
	return err
}

// ReadDiscreteInputs позволяет читать состояние дискретных входов (Discrete Inputs).
func (d *ModbusRTUDriver) ReadDiscreteInputs(address, quantity uint16) ([]bool, error) {
	results, err := d.client.ReadDiscreteInputs(address, quantity)
	if err != nil {
		return nil, err
	}
	boolResults := make([]bool, len(results)*8)
	for i, byteVal := range results {
		for bit := 0; bit < 8; bit++ {
			boolResults[i*8+bit] = byteVal&(1<<bit) != 0
		}
	}
	return boolResults, nil
}

// Позволяет читать значения из input регистров.
func (d *ModbusRTUDriver) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	return d.client.ReadInputRegisters(address, quantity)
}

func (d *ModbusRTUDriver) WriteMultipleRegisters(address uint16, values []byte) error {
	_, err := d.client.WriteMultipleRegisters(address, uint16(len(values)/2), values)
	return err
}

// Позволяет читать состояние coils.
func (d *ModbusRTUDriver) ReadCoils(address, quantity uint16) ([]bool, error) {
	results, err := d.client.ReadCoils(address, quantity)
	if err != nil {
		return nil, err
	}
	boolResults := make([]bool, len(results)*8)
	for i, byteVal := range results {
		for bit := 0; bit < 8; bit++ {
			boolResults[i*8+bit] = byteVal&(1<<bit) != 0
		}
	}
	return boolResults, nil
}

// изменить состояние одного coil
func (d *ModbusRTUDriver) WriteSingleCoil(address uint16, value bool) error {
	var val uint16
	if value {
		val = 0xFF00
	}
	_, err := d.client.WriteSingleCoil(address, val)
	return err
}
