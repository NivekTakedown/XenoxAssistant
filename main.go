package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NivekTakedown/XenoxAssistants/llm_handler/text"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// User details structure - simplified
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
}

// UserConversations with minimal user info
type UserConversations struct {
	UserID        int64          `json:"user_id"`
	FirstName     string         `json:"first_name"`
	PhoneNumber   string         `json:"phone_number,omitempty"`
	Conversations []Conversation `json:"conversations"`
}

// Conversation without duplicating user info
type Conversation struct {
	UserMessage  string `json:"user_message"`
	BotResponse  string `json:"bot_response"`
	ArrivalTime  string `json:"arrival_time"`
	ResponseTime string `json:"response_time"`
}

// Modified function signature
func logConversation(update tgbotapi.Update, botResp string, arrivalTime, responseTime string) error {
	user := User{
		ID:        update.Message.From.ID,
		FirstName: update.Message.From.FirstName,
	}

	conv := Conversation{
		UserMessage:  update.Message.Text,
		BotResponse:  botResp,
		ArrivalTime:  arrivalTime,
		ResponseTime: responseTime,
	}

	var users map[int64]UserConversations
	file := "conversations.json"

	data, err := os.ReadFile(file)
	if err == nil {
		if err := json.Unmarshal(data, &users); err != nil {
			return err
		}
	} else {
		users = make(map[int64]UserConversations)
	}

	userConv, exists := users[user.ID]
	if !exists {
		userConv = UserConversations{
			UserID:      user.ID,
			FirstName:   user.FirstName,
			PhoneNumber: getPhoneNumber(update.Message),
		}
	}

	userConv.Conversations = append(userConv.Conversations, conv)

	users[user.ID] = userConv

	newData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, newData, 0644)
}

// Add this helper function
func getPhoneNumber(message *tgbotapi.Message) string {
	if message != nil && message.Contact != nil {
		return message.Contact.PhoneNumber
	}
	return ""
}

func main() {
	// Obtén el token del bot desde una variable de entorno.
	// Esto es más seguro que hardcodearlo en el código.
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("La variable de entorno TELEGRAM_BOT_TOKEN no está definida.")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true // Activa el modo de depuración para ver más información en la consola

	log.Printf("Autorizado en la cuenta %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Teclado para solicitar teléfono
	replyContactRequest := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("Enviar teléfono"),
		),
	)

	for update := range updates {
		if update.Message == nil { // Ignora updates no relacionados con mensajes
			continue
		}

		arrivalTime := time.Now().Format(time.RFC3339)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		fmt.Printf("Mensaje recibido: [%s] %s\n", update.Message.From.UserName, update.Message.Text)

		var respuesta string
		if update.Message.Text == "/start" {
			respuesta = "¡Hola! Soy un asistente ético. ¿En qué puedo ayudarte? Por favor, comparte tu número telefónico."
		} else {
			// Use text handler from submodule
			genAIResp, err := text.HandleText(update.Message.Text)
			if err != nil {
				log.Printf("Error generating response: %v", err)
				respuesta = "Lo siento, hubo un error al procesar tu solicitud."
			} else {
				respuesta = genAIResp
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, respuesta)
		if update.Message.Text == "/start" {
			msg.ReplyMarkup = replyContactRequest
		}

		responseTime := time.Now().Format(time.RFC3339)
		err = logConversation(update, respuesta, arrivalTime, responseTime)
		if err != nil {
			log.Printf("Error al guardar conversación: %v", err)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error al enviar mensaje: %v", err)
		}
	}
}
