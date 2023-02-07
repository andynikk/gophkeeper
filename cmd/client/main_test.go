package main

import (
	"errors"
	"fmt"
	"gophkeeper/internal/client"
	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/cryptography"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/token"
	"testing"
)

func TestFuncClient(t *testing.T) {
	config := environment.ClientConfig{}
	config.InitConfigAgentENV()
	t.Run("Checking init config", func(t *testing.T) {
		if config.Address == "" {
			t.Errorf("Error checking init config")
		}
	})

	c := client.NewClient()
	t.Run("Checking init client", func(t *testing.T) {
		if c.TextDefault == "" {
			t.Errorf("Error checking init client")
		}
	})
	c.Config.CryptoKey = "test key for data encryption/decryption"

	tests := []struct {
		text        string
		encryptText string
		decryptText string
		gzipText    []byte
	}{
		{text: "Проверка кирилицы", encryptText: "", decryptText: "", gzipText: []byte("")},
		{text: "Checking the Latin alphabet", encryptText: "", decryptText: "", gzipText: []byte("")},
		{text: "检查汉字", encryptText: "", decryptText: "", gzipText: []byte("")},
		{text: "التحقق من الأبجدية العربية", encryptText: "", decryptText: "", gzipText: []byte("")},
		{text: "Checking the numbers 12345", encryptText: "", decryptText: "", gzipText: []byte("")},
	}

	t.Run("Checking crypt", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("Check encrypion %s", tt.text), func(t *testing.T) {
				tt.encryptText = encryption.EncryptString(tt.text, c.Config.CryptoKey)
				if tt.text == tt.encryptText || tt.encryptText == "" {
					t.Errorf(fmt.Sprintf("Check encrypion %s", tt.text))
				}
			})
			t.Run(fmt.Sprintf("Check decrypt %s", tt.text), func(t *testing.T) {
				tt.decryptText = encryption.DecryptString(tt.encryptText, c.Config.CryptoKey)
				if tt.text != tt.decryptText || tt.decryptText == "" || tt.decryptText == tt.encryptText {
					t.Errorf(fmt.Sprintf("Check decrypt %s", tt.text))
				}
			})
		}
	})

	t.Run("Checking gzip", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("Check compress gzip %s", tt.text), func(t *testing.T) {
				gzipText, err := compression.Compress([]byte(tt.text))
				tt.gzipText = gzipText
				if err != nil {
					t.Errorf(fmt.Sprintf("Check encrypion %s", tt.text))
				}
			})
			t.Run(fmt.Sprintf("Check decrypt %s", tt.text), func(t *testing.T) {
				text, err := compression.Decompress(tt.gzipText)
				if err != nil || string(text) != tt.text {
					t.Errorf(fmt.Sprintf("Check decrypt %s", tt.text))
				}
			})
		}
	})

	t.Run("Checking Hash SHA 256", func(t *testing.T) {
		configKey := "Test key hash SHA 256"
		for _, tt := range tests {
			hashData := cryptography.HashSHA256(tt.text, configKey)
			if hashData == "" || len(hashData) != 64 {
				t.Errorf("Error checking Hash SHA 256 (%s)", tt.text)
			}
		}
	})

	t.Run("Checking token", func(t *testing.T) {
		tokenString := ""
		userName := "Test username"
		t.Run("Checking token create", func(t *testing.T) {
			tc := token.NewClaims(userName)
			tokenString, _ = tc.GenerateJWT()
			if tokenString == "" {
				t.Errorf("Error checking token create (%s)", tokenString)
			}
		})
		t.Run("Checking token create", func(t *testing.T) {
			claims, ok := token.ExtractClaims(tokenString)
			if !ok || claims["user"] != userName {
				t.Errorf("Error checking token create (%s)", tokenString)
			}
		})
	})

	t.Run("Checking logger", func(t *testing.T) {
		t.Run("Checking error log", func(t *testing.T) {
			constants.Logger.ErrorLog(errors.New("test error"))
			//if constants.Logger.Log. = zerolog.DebugLevel {
			//
			//}
		})
		t.Run("Checking token create", func(t *testing.T) {
			constants.Logger.InfoLog("test info")
		})
	})

}
