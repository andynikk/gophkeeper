package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"gophkeeper/internal/constants"
)

type KeyRSA struct {
	User  string
	Patch string
	Key   string
}

// DecryptString дешифрует строку. Использует текстовый ключ. Если возникает ошибка, то возвращает изночальную строку
func DecryptString(cryptoText string, keyString string) string {

	if keyString == "" {
		return cryptoText
	}

	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return cryptoText
	}
	if len(encrypted) < aes.BlockSize {
		constants.Logger.ErrorLog(err)
		return cryptoText
	}

	decrypted, err := decryptAES(hashTo32Bytes(keyString), encrypted)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return cryptoText
	}

	return string(decrypted)
}

// EncryptString шифрует строку. Использует текстовый ключ. Если возникает ошибка, то возвращает изночальную строку
func EncryptString(plainText string, keyString string) string {

	if keyString == "" {
		return plainText
	}
	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return plainText
	}

	return base64.URLEncoding.EncodeToString(encrypted)
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}

func encryptAES(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func hashTo32Bytes(input string) []byte {

	data := sha256.Sum256([]byte(input))
	return data[0:]

}
