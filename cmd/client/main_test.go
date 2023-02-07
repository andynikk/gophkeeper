package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/client"
	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/cryptography"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/postgresql"
	"gophkeeper/internal/tests"
	"gophkeeper/internal/token"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestFuncClient(t *testing.T) {
	config := environment.ClientConfig{}
	t.Run("Checking init config", func(t *testing.T) {
		config.InitConfigAgentENV()
		if config.Address == "" {
			t.Errorf("Error checking init config")
		}
	})

	c := &client.Client{}
	t.Run("Checking init client", func(t *testing.T) {
		c = client.NewClient()
		if c.Config.Address == "" {
			t.Errorf("Error checking init client")
		}
	})
	c.Config.CryptoKey = "test key for data encryption/decryption"

	arrTests := []struct {
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

		for _, tt := range arrTests {
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
		for _, tt := range arrTests {
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
		for _, tt := range arrTests {
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
		})
		t.Run("Checking token create", func(t *testing.T) {
			constants.Logger.InfoLog("test info")
		})
	})

	t.Run("Checking post pair login password", func(t *testing.T) {
		events := [2]string{constants.EventAddEdit.String(), constants.EventDel.String()}
		for i := 0; i < 2; i++ {
			event := events[i]
			t.Run("Checking post pair login password",
				func(t *testing.T) {

					testStruct := tests.CreatePairLoginPassword(c.Token, event, c.Config.CryptoKey)

					body, err := json.Marshal(testStruct)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypePairLoginPassword.String(), event))
					}
					body, err = compression.Compress(body)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypePairLoginPassword.String(), event))
					}

					resp := httptest.NewRecorder()
					req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/resource/pairs", c.Config.Address),
						strings.NewReader(string(body)))
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypePairLoginPassword.String(), event))
					}
					http.DefaultServeMux.ServeHTTP(resp, req)
					if p, err := io.ReadAll(resp.Body); err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypePairLoginPassword.String(), event))
					} else {
						if string(p) != "" {
							t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypePairLoginPassword.String(), event))
						}
					}

				})
		}
	})

	t.Run(fmt.Sprintf("Checking post Text"), func(t *testing.T) {
		events := []string{constants.EventAddEdit.String(), constants.EventDel.String()}
		for i := 0; i < 2; i++ {
			event := events[i]
			t.Run("Checking post Text",
				func(t *testing.T) {

					testStruct := tests.CreateTextData(c.Token, event, c.Config.CryptoKey)

					body, err := json.Marshal(testStruct)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeTextData.String(), event))
					}
					body, err = compression.Compress(body)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeTextData.String(), event))
					}

					resp := httptest.NewRecorder()
					req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/resource/pairs", c.Config.Address),
						strings.NewReader(string(body)))
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeTextData.String(), event))
					}
					http.DefaultServeMux.ServeHTTP(resp, req)
					if p, err := io.ReadAll(resp.Body); err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeTextData.String(), event))
					} else {
						if string(p) != "" {
							t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeTextData.String(), event))
						}
					}

				})
		}
	})

	t.Run("Checking post Binary", func(t *testing.T) {
		events := [2]string{constants.EventAddEdit.String(), constants.EventDel.String()}
		for i := 0; i < 2; i++ {
			event := events[i]
			t.Run("Checking post Binary",
				func(t *testing.T) {

					testStruct := tests.CreateBinaryData(c.Token, event)

					body, err := json.Marshal(testStruct)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBinaryData.String(), event))
					}
					body, err = compression.Compress(body)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBinaryData.String(), event))
					}

					resp := httptest.NewRecorder()
					req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/resource/pairs", c.Config.Address),
						strings.NewReader(string(body)))
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBinaryData.String(), event))
					}
					http.DefaultServeMux.ServeHTTP(resp, req)
					if p, err := io.ReadAll(resp.Body); err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBinaryData.String(), event))
					} else {
						if string(p) != "" {
							t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBinaryData.String(), event))
						}
					}

				})
		}
	})

	t.Run("Checking post Bank card", func(t *testing.T) {
		events := []string{constants.EventAddEdit.String(), constants.EventDel.String()}
		for i := 0; i < 2; i++ {
			event := events[i]
			t.Run("Checking post Bank card",
				func(t *testing.T) {

					testStruct := tests.CreateBankCard(c.Token, event, c.Config.CryptoKey)

					body, err := json.Marshal(testStruct)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBankCardData.String(), event))
					}
					body, err = compression.Compress(body)
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBankCardData.String(), event))
					}

					resp := httptest.NewRecorder()
					req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/resource/pairs", c.Config.Address),
						strings.NewReader(string(body)))
					if err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBankCardData.String(), event))
					}
					http.DefaultServeMux.ServeHTTP(resp, req)
					if p, err := io.ReadAll(resp.Body); err != nil {
						t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBankCardData.String(), event))
					} else {
						if string(p) != "" {
							t.Errorf(fmt.Sprintf("Error checking post %s %s", constants.TypeBankCardData.String(), event))
						}
					}

				})
		}
	})

	t.Run("Checking login User", func(t *testing.T) {

		testStruct := postgresql.User{
			Name:         "user1",
			Password:     "password",
			HashPassword: "",
			Event:        "edit",
		}

		body, err := json.Marshal(testStruct)
		if err != nil {
			t.Errorf("Error checking login User")
		}
		body, err = compression.Compress(body)
		if err != nil {
			t.Errorf("Error checking login User")
		}

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/user/login", c.Config.Address),
			strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error checking login User")
		}
		http.DefaultServeMux.ServeHTTP(resp, req)
		if p, err := io.ReadAll(resp.Body); err != nil {
			t.Errorf("Error checking login User")
		} else {
			if string(p) != "" {
				t.Errorf("Error checking login User")
			}
		}
	})

	t.Run("Checking register User", func(t *testing.T) {

		testStruct := postgresql.User{
			Name:         "user1",
			Password:     "password",
			HashPassword: "",
			Event:        "edit",
		}

		body, err := json.Marshal(testStruct)
		if err != nil {
			t.Errorf("Error checking register User")
		}
		body, err = compression.Compress(body)
		if err != nil {
			t.Errorf("Error checking register User")
		}

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/user/register", c.Config.Address),
			strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error checking register User")
		}
		http.DefaultServeMux.ServeHTTP(resp, req)
		if p, err := io.ReadAll(resp.Body); err != nil {
			t.Errorf("Error checking register User")
		} else {
			if string(p) != "" {
				t.Errorf("Error checking register User")
			}
		}
	})
}

func BenchmarkPlp(b *testing.B) {

	config := environment.ClientConfig{}
	config.InitConfigAgentENV()

	c := client.NewClient()
	c.Config.CryptoKey = "test key for data encryption/decryption"
	addressPost := fmt.Sprintf("http://%s/api/resource/pairs", c.Config.Address) //a.cfg.Address)

	wg := sync.WaitGroup{}
	plp := tests.CreatePairLoginPassword(c.Token, "edit", c.Config.CryptoKey)
	plpJSON, err := json.MarshalIndent(plp, "", " ")

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err != nil {
				return
			}
			_, err = client.ExecuteAPI(plpJSON, addressPost, c.Token)
			return
		}()
	}
	wg.Wait()

	plp.Event = constants.EventDel.String()
	plpJSON, err = json.MarshalIndent(plp, "", " ")
	if err != nil {
		return
	}
	_, err = client.ExecuteAPI(plpJSON, addressPost, c.Token)
}
