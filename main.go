package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Settings struct {
	Token  string `json:"token"`
	ChatID int64  `json:"chat_id"`
}

func main() {
	// Загружаем настройки
	var settings Settings
	//var people []Person

	// Открываем JSON-файл
	jsonFile, err := os.Open("setting.json")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %s", err)
	}
	defer jsonFile.Close()

	// Получаем информацию о файле
	fileInfo, err := jsonFile.Stat()
	if err != nil {
		log.Fatalf("Ошибка при получении информации о файле: %s", err)
	}

	// Создаем срез байтов для хранения содержимого файла
	byteValue := make([]byte, fileInfo.Size())

	// Читаем содержимое файла
	_, err = jsonFile.Read(byteValue)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %s", err)
	}

	// Парсим JSON данные
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %s", err)
	}
	// Создаем нового бота, используя токен, который вы получили от BotFather
	bot, err := tgbotapi.NewBotAPI(settings.Token)
	if err != nil {
		log.Panic(err)
	}

	// ID чата, в который будет отправляться сообщение
	chatID := int64(settings.ChatID) // Замените на ID вашего чата

	// Запускаем задачу, которая будет отправлять сообщение каждый день в определенное время
	scheduleMessage(bot, chatID)

	select {} // Бесконечный цикл, чтобы программа не завершилась
}

// Функция для отправки сообщения по расписанию
func scheduleMessage(bot *tgbotapi.BotAPI, chatID int64) {
	go func() {
		for {
			// Задаем время для отправки сообщения (например, 9:00 утра каждый день)
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 21, 34, 0, 0, now.Location())
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}

			// Ждем до нужного времени
			duration := next.Sub(now)
			fmt.Printf("Следующее сообщение будет отправлено через: %v\n", duration)
			time.Sleep(duration)

			// Отправляем сообщение
			msg := tgbotapi.NewMessage(chatID, "Доброе утро! Это сообщение по расписанию.")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
			}
		}
	}()
}
