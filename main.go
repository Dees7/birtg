package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Settings struct {
	Token  string `json:"token"`
	ChatID int64  `json:"chat_id"`
	Time   string `json:"time"`
}
type Person struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Day   int    `json:"day"`
	Month int    `json:"month"`
}

func getBirthdays() []Person {
	var people []Person

	jsonFile, err := os.Open("birthdays.json")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %s", err)
	}
	defer jsonFile.Close()

	byteValue, err := os.ReadFile("birthdays.json")
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %s", err)
	}

	err = json.Unmarshal(byteValue, &people)
	if err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %s", err)
	}
	return people
}

func main() {
	// Загружаем настройки
	var settings Settings
	var people []Person

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
	people = getBirthdays()

	// Разделяем строку по символу ":"
	parts := strings.Split(settings.Time, ":")

	// Преобразуем части строки в числа
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])

	scheduleMessage(bot, chatID, people, hours, minutes)

	select {} // Бесконечный цикл, чтобы программа не завершилась
}

// Функция для отправки сообщения по расписанию
func scheduleMessage(bot *tgbotapi.BotAPI, chatID int64, people []Person, hours int, minutes int) {
	go func() {
		for {
			// Задаем время для отправки сообщения (например, 9:00 утра каждый день)
			currentTime := time.Now()
			if currentTime.Hour() == hours && currentTime.Minute() == minutes {
				for _, person := range people {
					if person.Day == currentTime.Day() && person.Month == int(currentTime.Month()) {
						// Отправляем сообщение
						msg := tgbotapi.NewMessage(chatID, "С днем рождения, "+string(person.Name)+"! "+string(person.Login))
						if _, err := bot.Send(msg); err != nil {
							log.Printf("Ошибка при отправке сообщения: %v", err)
						}
					}
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()
}
