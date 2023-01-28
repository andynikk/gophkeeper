package encryption

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"gophkeeper/internal/constants"
)

type KeyRSA struct {
	User         string
	Patch        string
	NumSert      int64
	SubjectKeyID string
	LenKeyByte   int
}

type KeyEncryption struct {
	TypeEncryption string
	PublicKey      *rsa.PublicKey
	PrivateKey     *rsa.PrivateKey
}

func (key *KeyEncryption) RsaEncrypt(msg []byte) ([]byte, error) {
	if key == nil {
		return msg, nil
	}
	encryptedBytes, err := rsa.EncryptOAEP(sha512.New512_256(), rand.Reader, key.PublicKey, msg, nil)
	return encryptedBytes, err
}

func (key *KeyEncryption) RsaDecrypt(msgByte []byte) ([]byte, error) {
	if key == nil {
		return msgByte, nil
	}
	msgByte, err := key.PrivateKey.Decrypt(nil, msgByte, &rsa.OAEPOptions{Hash: crypto.SHA512_256})
	return msgByte, err
}

func (k *KeyRSA) CreateCert() ([]bytes.Buffer, error) {

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(k.NumSert),
		Subject: pkix.Name{
			Organization: []string{"AdvancedMetrics"},
			Country:      []string{"RU"},
		},
		NotBefore: time.Now(),
		NotAfter: time.Now().AddDate(constants.TimeLivingCertificateYaer, constants.TimeLivingCertificateMounth,
			constants.TimeLivingCertificateDay),
		SubjectKeyId: []byte(k.SubjectKeyID),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, k.LenKeyByte)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	var certPEM bytes.Buffer
	_ = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	_ = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return []bytes.Buffer{certPEM, privateKeyPEM}, nil
}

func SaveKeyInFile(key *bytes.Buffer, pathFile string) {

	file, err := os.Create(pathFile)
	if err != nil {
		return
	}
	_, err = file.WriteString(key.String())
	if err != nil {
		return
	}
}

func InitPrivateKey(cryptoKeyPath string) (*KeyEncryption, error) {

	if cryptoKeyPath == "" {
		return nil, errors.New("путь к приватному ключу не указан")
	}
	pvkData, _ := os.ReadFile(cryptoKeyPath)
	pvkBlock, _ := pem.Decode(pvkData)
	pvk, err := x509.ParsePKCS1PrivateKey(pvkBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &KeyEncryption{TypeEncryption: constants.TypeEncryption, PrivateKey: pvk, PublicKey: &pvk.PublicKey}, nil
}

func InitPublicKey(cryptoKeyPath string) (*KeyEncryption, error) {
	if cryptoKeyPath == "" {
		return nil, errors.New("не указан путь к публичному ключу")
	}
	certData, _ := os.ReadFile(cryptoKeyPath)
	certBlock, _ := pem.Decode(certData)
	cert, _ := x509.ParseCertificate(certBlock.Bytes)
	certPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return &KeyEncryption{TypeEncryption: constants.TypeEncryption, PublicKey: certPublicKey}, nil
}

func DecryptString(cryptoText string, keyString string) (plainTextString string, err error) {

	newKeyString, err := hashTo32Bytes(keyString)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return cryptoText, err
	}

	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(newKeyString))
	if err != nil {
		constants.Logger.ErrorLog(err)
		return cryptoText, err
	}

	if len(cipherText) < aes.BlockSize {
		constants.Logger.ErrorLog(err)
		return cryptoText, errors.New("это не шифрованный текст")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}

func EncryptString(plainText string, keyString string) (cipherTextString string, err error) {

	newKeyString, err := hashTo32Bytes(keyString)

	if err != nil {
		return "", err
	}

	key := []byte(newKeyString)
	value := []byte(plainText)

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(value))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], value)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func hashTo32Bytes(input string) (output string, err error) {

	if len(input) == 0 {
		return "", errors.New("No input supplied")
	}

	hasher := sha256.New()
	hasher.Write([]byte(input))

	stringToSHA256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return stringToSHA256[:32], nil
}

func main() {

	argumentsCount := len(os.Args)
	if argumentsCount != 6 {
		fmt.Printf("Usage:\n-e to encrypt, -d to decrypt.\n")
		fmt.Printf("--key \"I am a key\" to load the key.\n")
		fmt.Printf("--value \"I am a text to be encrypted or decrypted\".\n")
		return
	}

	encrypt := false
	decrypt := false
	key := false
	expectKeyString := 0
	keyString := false
	value := false
	expectValueString := 0
	valueString := false

	encryptionFlag := ""
	stringToEncrypt := ""
	encryptionKey := ""

	for index, element := range os.Args {

		if element == "-e" {
			if decrypt == true {
				fmt.Printf("Can't set both -e and -d.\nBye!\n")
				return
			}
			encrypt = true
			encryptionFlag = "-e"

		} else if element == "-d" {
			if encrypt == true {
				fmt.Printf("Can't set both -e and -d.\nBye!\n")
				return
			}
			decrypt = true
			encryptionFlag = "-d"

		} else if element == "--key" {
			key = true
			expectKeyString++

		} else if element == "--value" {
			value = true
			expectValueString++

		} else if expectKeyString == 1 {
			encryptionKey = os.Args[index]
			keyString = true
			expectKeyString = 0

		} else if expectValueString == 1 {
			stringToEncrypt = os.Args[index]
			valueString = true
			expectValueString = 0
		}

		if expectKeyString >= 2 {
			fmt.Printf("Something went wrong, too many keys entered.\bBye!\n")
			return

		} else if expectValueString >= 2 {
			fmt.Printf("Something went wrong, too many keys entered.\bBye!\n")
			return
		}
	}

	if !(encrypt == true || decrypt == true) || key == false || keyString == false || value == false || valueString == false {
		fmt.Printf("Incorrect usage!\n")
		fmt.Printf("---------\n")
		fmt.Printf("-e or -d -> %v\n", (encrypt == true || decrypt == true))
		fmt.Printf("--key -> %v\n", key)
		fmt.Printf("Key string? -> %v\n", keyString)
		fmt.Printf("--value -> %v\n", value)
		fmt.Printf("Value string? -> %v\n", valueString)
		fmt.Printf("---------")
		fmt.Printf("\nUsage:\n-e to encrypt, -d to decrypt.\n")
		fmt.Printf("--key \"I am a key\" to load the key.\n")
		fmt.Printf("--value \"I am a text to be encrypted or decrypted\".\n")
		return
	}

	if false == (encryptionFlag == "-e" || encryptionFlag == "-d") {
		fmt.Println("Sorry but the first argument has to be either -e or -d")
		fmt.Println("for either encryption or decryption.")
		return
	}

	if encryptionFlag == "-e" {

		fmt.Printf("Encrypting '%s' with key '%s'\n", stringToEncrypt, encryptionKey)

		encryptedString, _ := EncryptString(stringToEncrypt, encryptionKey)

		fmt.Printf("Output: '%s'\n", encryptedString)

	} else if encryptionFlag == "-d" {
		// Decrypt!

		fmt.Printf("Decrypting '%s' with key '%s'\n", stringToEncrypt, encryptionKey)

		decryptedString, _ := DecryptString(stringToEncrypt, encryptionKey)

		fmt.Printf("Output: '%s'\n", decryptedString)

	}
}
