package adapters

import (
	"log"

	"github.com/brutella/can"
)

// CANAdapter представляет адаптер для работы с CAN-шиной.
type CANAdapter struct {
	bus      *can.Bus
	stopChan chan struct{}
}

func (adapter *CANAdapter) HandleFrame(frame can.Frame) {
	log.Printf("Received CAN frame: ID=%X Data=%X\n", frame.ID, frame.Data)
}

// NewCANAdapter создает новый экземпляр CANAdapter с указанным интерфейсом.
func NewCANAdapter(interfaceName string) (*CANAdapter, error) {
	bus, err := can.NewBusForInterfaceWithName(interfaceName)
	if err != nil {
		return nil, err
	}

	adapter := &CANAdapter{
		bus:      bus,
		stopChan: make(chan struct{}),
	}

	// Подписка на шину CAN с использованием адаптера как обработчика
	bus.SubscribeFunc(adapter.HandleFrame)

	return adapter, nil
}

// Connect подключает адаптер к CAN-шине и начинает слушать сообщения.
func (adapter *CANAdapter) Connect() error {
	log.Println("Connecting to CAN bus")
	go adapter.listen()
	return nil
}

// listen запускает в горутине прослушивание канала сообщений CAN.
func (adapter *CANAdapter) listen() {
	for {
		select {
		case <-adapter.stopChan:
			return
		}
	}
}

// Disconnect отключает адаптер от CAN-шине и останавливает слушание сообщений.
func (adapter *CANAdapter) Disconnect() error {
	log.Println("Disconnecting from CAN bus")
	close(adapter.stopChan)
	return nil // Поскольку нет явного отключения, просто возвращаем nil
}

// Write отправляет данные на CAN-шину.
func (adapter *CANAdapter) Write(frame can.Frame) error {
	log.Printf("Writing to CAN bus: ID=%X Data=%X\n", frame.ID, frame.Data)
	return adapter.bus.Publish(frame)
}
