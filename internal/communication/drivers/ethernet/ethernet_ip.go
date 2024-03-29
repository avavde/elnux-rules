package communication

import (
	"log"

	"github.com/loki-os/go-ethernet-ip/client"
)

type EthernetIPDriver struct {
	c *client.Client
}

func NewEthernetIPDriver(ip string) *EthernetIPDriver {
	// Создание нового клиента EtherNet/IP
	c := client.NewClient(ip)

	return &EthernetIPDriver{
		c: c,
	}
}

// Connect устанавливает соединение с устройством EtherNet/IP
func (driver *EthernetIPDriver) Connect() error {
	if err := driver.c.RegisterSession(); err != nil {
		log.Printf("Failed to connect to EtherNet/IP device: %s\n", err)
		return err
	}

	log.Println("Connected to EtherNet/IP device")
	return nil
}

// Disconnect закрывает соединение с устройством
func (driver *EthernetIPDriver) Disconnect() {
	if err := driver.c.UnRegisterSession(); err != nil {
		log.Printf("Failed to disconnect from EtherNet/IP device: %s\n", err)
	}
	log.Println("Disconnected from EtherNet/IP device")
}

// ReadTag читает значение тега с устройства
func (driver *EthernetIPDriver) ReadTag(tagName string) ([]byte, error) {
	response, err := driver.c.ReadTagService(tagName)
	if err != nil {
		log.Printf("Failed to read tag '%s': %s\n", tagName, err)
		return nil, err
	}

	log.Printf("Read tag '%s': %v\n", tagName, response)
	return response, nil
}

// WriteTag записывает значение в тег устройства
func (driver *EthernetIPDriver) WriteTag(tagName string, data []byte) error {
	if err := driver.c.WriteTagService(tagName, data); err != nil {
		log.Printf("Failed to write tag '%s': %s\n", tagName, err)
		return err
	}

	log.Printf("Wrote tag '%s'\n", tagName)
	return nil
}
