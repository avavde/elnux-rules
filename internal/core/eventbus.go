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
	subscriptions, ok := eb.subscriptions[event.Type]
	eb.lock.Unlock()

	if !ok {
		fmt.Printf("Нет подписок на событие типа '%s'\n", event.Type)
		return
	}

	for _, sub := range subscriptions {
		if passFilters(event, sub.Filters) {
			fmt.Printf("Публикация события: %s\n", event.Type)
			go sub.Handler(event)
		}
	}
}

// passFilters проверяет, проходит ли событие через все фильтры подписки.
func passFilters(event Event, filters []FilterFunc) bool {
	for _, filter := range filters {
		if !filter(event) {
			return false // Событие не прошло один из фильтров
		}
	}
	return true // Событие прошло все фильтры
}
