package communication

import "net"

// EthernetProtocol описывает интерфейс для работы с Ethernet-соединениями.
type EthernetProtocol interface {
	Connect(address string) (net.Conn, error)     // Установление соединения с устройством
	Disconnect(conn net.Conn) error               // Закрытие соединения с устройством
	Send(conn net.Conn, data []byte) (int, error) // Отправка данных устройству
	Receive(conn net.Conn) ([]byte, error)        // Получение данных от устройства
}
