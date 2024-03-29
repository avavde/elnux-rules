package adapters

import (
	"fmt"
	"log"
	"time"

	"github.com/goburrow/modbus"
)

// EthernetAdapter представляет адаптер для работы через Ethernet.
type EthernetAdapter struct {
	client  modbus.Client
	handler *modbus.TCPClientHandler
}

// NewEthernetAdapter создает новый экземпляр EthernetAdapter для заданного адреса и порта.
func NewEthernetAdapter(address string, port int) *EthernetAdapter {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", address, port))
	handler.Timeout = 10 * time.Second // Установка таймаута соединения
	handler.SlaveId = 1                // Установка идентификатора устройства

	return &EthernetAdapter{
		client:  modbus.NewClient(handler),
		handler: handler, // Сохраняем handler в адаптере
	}
}

// ReadRegisters читает регистры устройства, начиная с указанного адреса и количество.
func (adapter *EthernetAdapter) ReadRegisters(startAddr uint16, quantity uint16) ([]byte, error) {
	// Чтение holding регистров
	results, err := adapter.client.ReadHoldingRegisters(startAddr, quantity)
	if err != nil {
		log.Printf("Ошибка чтения регистров: %v", err)
		return nil, err
	}
	return results, nil
}

// WriteRegister записывает значение в регистр устройства по указанному адресу.
func (adapter *EthernetAdapter) WriteRegister(addr uint16, value uint16) error {
	// Запись в один holding регистр
	results, err := adapter.client.WriteSingleRegister(addr, value)
	if err != nil {
		log.Printf("Ошибка записи в регистр: %v", err)
		return err
	}
	log.Printf("Регистр записан успешно: %v", results)
	return nil
}

// Disconnect закрывает соединение с устройством.
func (adapter *EthernetAdapter) Disconnect() {
	// В случае использования TCP, закрытие handler соединения
	if adapter.handler != nil {
		adapter.handler.Close()
		log.Println("Соединение через Ethernet закрыто")
	}
}
