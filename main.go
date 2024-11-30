package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Создаем нового бота, используя токен, который вы получили от BotFather
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}

	// ID чата, в который будет отправляться сообщение
	chatID := int64(-1) // Замените на ID вашего чата

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
