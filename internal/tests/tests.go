// Package tests: функции для создания данных для тестирования
package tests

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/cryptography"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql/model"
)

// CreateUser создание сущности User для тестов
func CreateUser(event string) model.User {
	return model.User{
		Name:         "test",
		Password:     "password",
		HashPassword: cryptography.HashSHA256("password", string(constants.HashKey)),
		Event:        event,
	}
}

// CreatePairLoginPassword создание сущности PairLoginPassword для тестов
func CreatePairLoginPassword(token, event, cryptoKey string) model.PairLoginPassword {
	return model.PairLoginPassword{
		User:     token,
		Uid:      "bf340769-687e-485e-968b-976cf12f7b64",
		TypePair: encryption.EncryptString(constants.TypePairLoginPassword.String(), cryptoKey),
		Name:     encryption.EncryptString("yandex.ru", cryptoKey),
		Password: encryption.EncryptString("test_password", cryptoKey),
		Event:    event,
	}
}

// CreateTextData создание сущности TextData для тестов
func CreateTextData(token, event, cryptoKey string) model.TextData {
	return model.TextData{
		User:  token,
		Uid:   "bf340769-687e-485e-968b-976cf12f7b64",
		Text:  encryption.EncryptString("Text test", cryptoKey),
		Event: event,
	}
}

// CreateBinaryData создание сущности BinaryData для тестов
func CreateBinaryData(token, event string) model.BinaryData {
	return model.BinaryData{
		User:          token,
		Uid:           "bf340769-687e-485e-968b-976cf12f7b64",
		Patch:         "./temp.tmp",
		DownloadPatch: "./temp_1.tmp",
		Name:          "temp_1",
		Expansion:     "tmp",
		Size:          "1",
		Event:         event,
	}
}

// CreateBankCard создание сущности BankCard для тестов
func CreateBankCard(token, event, cryptoKey string) model.BankCard {
	return model.BankCard{
		User:   token,
		Uid:    "bf340769-687e-485e-968b-976cf12f7b64",
		Number: encryption.EncryptString("4342 5654 5878 5475", cryptoKey),
		Cvc:    encryption.EncryptString("333", cryptoKey),
		Event:  event,
	}
}
