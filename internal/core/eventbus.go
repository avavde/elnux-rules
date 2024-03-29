package core

import (
	"fmt"
	"sync"
)

// FilterFunc определяет тип функции-фильтра для событий.
type FilterFunc func(Event) bool

// Subscription структура для хранения подписки, включая обработчики и фильтры.
type Subscription struct {
	Handler EventHandler
	Filters []FilterFunc
}

// EventBus структура для управления событиями и подписчиками.
type EventBus struct {
	subscriptions map[string][]Subscription
	lock          sync.Mutex
}

// NewEventBus создаёт новый экземпляр EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscriptions: make(map[string][]Subscription),
	}
}

// Subscribe добавляет новый обработчик для заданного типа события с опциональными фильтрами.
func (eb *EventBus) Subscribe(eventType string, handler EventHandler, filters ...FilterFunc) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	subscription := Subscription{
		Handler: handler,
		Filters: filters,
	}
	eb.subscriptions[eventType] = append(eb.subscriptions[eventType], subscription)
}

// Publish публикует событие, вызывая все соответствующие обработчики, которые проходят фильтрацию.
func (eb *EventBus) Publish(event Event) {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	// Вызов обработчиков, подписанных на конкретные типы событий
	for _, sub := range eb.subscriptions[event.Type] {
		if passFilters(event, sub.Filters) {
			fmt.Printf("Публикация события: %s\n", event.Type)
			fmt.Println("Применяемые фильтры:")
			for _, filter := range sub.Filters {
				fmt.Printf("- %v\n", filter)
			}
			go sub.Handler(event)
		}
	}

	// Универсальная подписка для обработчиков, использующих шаблоны или префиксы
	if universalSubs, ok := eb.subscriptions["*"]; ok {
		for _, sub := range universalSubs {
			if passFilters(event, sub.Filters) {
				fmt.Printf("Публикация события: %s\n", event.Type)
				fmt.Println("Применяемые фильтры:")
				for _, filter := range sub.Filters {
					fmt.Printf("- %v\n", filter)
				}
				go sub.Handler(event)
			}
		}
	}
}

// passFilters проверяет, проходит ли событие через все фильтры подписки.
func passFilters(event Event, filters []FilterFunc) bool {
	for _, filter := range filters {
		if !filter(event) {
			fmt.Printf("Событие: %s не прошло фильтр\n", event.Type)
			return false // Событие не прошло один из фильтров
		}
	}
	fmt.Printf("Событие: %s прошло все фильтры\n", event.Type)
	return true // Событие прошло все фильтры
}
