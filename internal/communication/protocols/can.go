package communication

// Frame представляет CAN-фрейм.
type Frame struct {
	ID     uint32 // Идентификатор сообщения
	Data   []byte // Данные сообщения
	Length int    // Длина данных
}

// CANProtocol описывает интерфейс для работы с CAN-шиной.
type CANProtocol interface {
	Connect() error                        // Подключение к шине CAN
	Disconnect() error                     // Отключение от шины CAN
	Send(frame Frame) error                // Отправка сообщения в шину CAN
	Receive() (<-chan Frame, error)        // Получение сообщений из шины CAN
	SetFilter(filterID, mask uint32) error // Установка фильтра для приема сообщений
}
