package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestLogConversation(t *testing.T) {
	// Crear un archivo temporal
	tmpFile, err := os.CreateTemp("", "conversations_test_*.json")
	if err != nil {
		t.Fatalf("No se pudo crear el archivo temporal: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Sobrescribir la variable 'file' para usar el archivo temporal
	originalFile := "conversations.json"
	file := tmpFile.Name()
	defer func() { file = originalFile }()

	// Crear un mensaje de prueba
	message := &tgbotapi.Message{
		From: &tgbotapi.User{
			ID:        12345,
			FirstName: "TestUser",
		},
		Text: "Hola bot",
	}

	update := tgbotapi.Update{
		Message: message,
	}

	// Parámetros de prueba
	botResp := "Hola TestUser! Respuesta generada"
	arrivalTime := time.Now().Format(time.RFC3339)
	responseTime := time.Now().Format(time.RFC3339)

	// Ejecutar la función
	err = logConversation(update, botResp, arrivalTime, responseTime)
	if err != nil {
		t.Errorf("logConversation retornó un error: %v", err)
	}

	// Leer el contenido del archivo temporal
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("No se pudo leer el archivo temporal: %v", err)
	}

	// Verificar el contenido
	var users map[int64]UserConversations
	err = json.Unmarshal(data, &users)
	if err != nil {
		t.Fatalf("Error al deserializar JSON: %v", err)
	}

	userConv, exists := users[12345]
	if !exists {
		t.Fatalf("Usuario con ID 12345 no encontrado en las conversaciones")
	}

	if userConv.FirstName != "TestUser" {
		t.Errorf("FirstName esperado 'TestUser', obtenido '%s'", userConv.FirstName)
	}

	if len(userConv.Conversations) != 1 {
		t.Errorf("Se esperaban 1 conversación, obtenidas %d", len(userConv.Conversations))
	}

	conv := userConv.Conversations[0]
	if conv.UserMessage != "Hola bot" {
		t.Errorf("UserMessage esperado 'Hola bot', obtenido '%s'", conv.UserMessage)
	}
	if conv.BotResponse != botResp {
		t.Errorf("BotResponse esperado '%s', obtenido '%s'", botResp, conv.BotResponse)
	}
	if conv.ArrivalTime != arrivalTime {
		t.Errorf("ArrivalTime esperado '%s', obtenido '%s'", arrivalTime, conv.ArrivalTime)
	}
	if conv.ResponseTime != responseTime {
		t.Errorf("ResponseTime esperado '%s', obtenido '%s'", responseTime, conv.ResponseTime)
	}
}

func TestLogConversation_ErrorReadingFile(t *testing.T) {
	// Establecer un archivo que no existe
	file = "/ruta/inexistente/conversations.json"

	message := &tgbotapi.Message{
		From: &tgbotapi.User{
			ID:        67890,
			FirstName: "ErrorUser",
		},
		Text: "Este es un mensaje",
	}

	update := tgbotapi.Update{
		Message: message,
	}

	err := logConversation(update, "Respuesta de error", "2024-01-01T00:00:00Z", "2024-01-01T00:00:01Z")
	if err == nil {
		t.Errorf("Se esperaba un error al leer el archivo, pero no se obtuvo ninguno")
	}
}

func TestGetPhoneNumber(t *testing.T) {
	contactMsg := &tgbotapi.Message{
		Contact: &tgbotapi.Contact{
			PhoneNumber: "123456789",
		},
	}

	phone := getPhoneNumber(contactMsg)
	if phone != "123456789" {
		t.Errorf("Se esperaba '123456789', obtenido '%s'", phone)
	}

	nonContactMsg := &tgbotapi.Message{}
	phone = getPhoneNumber(nonContactMsg)
	if phone != "" {
		t.Errorf("Se esperaba una cadena vacía, obtenida '%s'", phone)
	}
}
