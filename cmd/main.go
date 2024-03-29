package main

import (
	"elnux-rules/internal/communication/drivers/serial"
	"elnux-rules/internal/core"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %v", err)
	}

	// чтение переменных окружения из .env
	dbHost := os.Getenv("DB_HOST")
	fmt.Println("DB_HOST:", dbHost)
	// - - - - закончили чтение переменных окружения из .env

	eventBus := core.NewEventBus() // шина событий
	time.Sleep(1 * time.Second)
	javaScriptExecutor := core.NewJavaScriptExecutor() // executor.og
	time.Sleep(1 * time.Second)
	ruleEngine := core.NewRuleEngine(eventBus, javaScriptExecutor) // движок правил

	// Загрузка и применение всех правил из директории .rules

	if err := ruleEngine.LoadAndApplyRules("./rules"); err != nil {
		log.Fatalf("Не удалось загрузить и применить правила: %v", err)
	}
	log.Println("Правила успешно загружены и применены")
	// - - - - - ПРОХОДИМ ТЕСТЫ - - - - - - -
	// Регистрируем тестовые правила
	core.RegisterTestRules(eventBus, ruleEngine)

	testEvents := []core.Event{
		{
			Type: "testhumidityChange",
			Payload: map[string]interface{}{
				"humidity": 20, // должно запустить "Внимание: влажность ниже 30%!"
			},
		},
		{
			Type: "testmotionDetected",
			Payload: map[string]interface{}{
				"motion": true,
				"time":   "2024-03-27T12:00:00Z", // должно запустить "Движение обнаружено. Включение освещения."
			},
		},
		{
			Type: "testtemperatureChange",
			Payload: map[string]interface{}{
				"temperature": 30, // должно запустить "Внимание: температура превысила 25 градусов!"
			},
		},
		{
			Type: "test-ExampleEvent",
			Payload: map[string]interface{}{
				"message": "test", // должно запустить "Температура превысила 30 градусов!"
			},
		},
	}

	// опубликовать все тестовые события
	for _, event := range testEvents {
		eventBus.Publish(event)
	}

	// - - - - - ПРОШЛИ ТЕСТЫ  - - - - - - -
	time.Sleep(7 * time.Second) // Интервал ожидания перед следующим чтением
	// ИНИЦИАЛИЗИРУЕМ УСТРОЙСТВА
	modbusDriver := serial.NewModbusRTUSerialDriver("/dev/ttyUSB0", 9600, 8, "N", 1)

	// Настройка обработки сигналов ОС для корректного завершения программы
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Бесконечный цикл ожидания событий
	go func() {
		for {
			select {
			case <-stopChan:
				// Завершение работы
				log.Println("Завершение работы по сигналу от ОС")
				return
			default:
				// Чтение данных с устройства через Modbus RTU
				data, err := modbusDriver.ReadHoldingRegisters(1, 0, 10) // чтение 10 последних
				if err != nil {
					log.Printf("Ошибка при чтении данных: %v", err)
					continue
				}
				log.Printf("Прочитанные данные: %v", data)

				// генерация событий

				time.Sleep(1 * time.Second) // Интервал ожидания перед следующим чтением
			}
		}
	}()

	// Ожидаем сигнал для завершения
	<-stopChan
	log.Println("Движок правил остановлен")
}
