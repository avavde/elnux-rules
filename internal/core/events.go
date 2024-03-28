package core

// Event тип представляет событие в системе.
type Event struct {
	Type    string                 // Тип события
	Payload map[string]interface{} // Полезная нагрузка события
}

// EventHandler функция, обрабатывающая события.
type EventHandler func(Event)
