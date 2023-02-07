package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

type additionalBinaryParameters struct {
	patch string
	uid   string
}

// createEncryptionKey событие формы, которое создает ключ для клиента
func (c *Client) createEncryptionKey(k encryption.KeyRSA) error {
	if err := os.WriteFile(k.Patch, []byte(k.Key), 0664); err != nil {
		constants.Logger.ErrorLog(err)
	}
	c.Config.CryptoKey = k.Key
	return nil
}

// downloadBinaryData событие формы, которое загружает файл с сервера и сохраняет на клиенте
func (c *Client) downloadBinaryData(bd postgresql.BinaryData) error {

	ctx := context.Background()
	ctxWV := context.WithValue(ctx, postgresql.KeyContext("additionalBinaryParameters"), additionalBinaryParameters{
		patch: bd.DownloadPatch,
		uid:   bd.Uid,
	})
	go c.wsDownloadBinaryData(ctxWV)
	return nil
}

// inputLoginUser событие формы, позволяет залогинится пользователю. Проверяется по имени и хешу пароля
func (c *Client) inputLoginUser(user postgresql.User) error {

	addressPost := fmt.Sprintf("http://%s/api/user/login", c.Config.Address) //a.cfg.Address)
	arrJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err
	}

	compressJSON, err := compression.Compress(arrJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	req, err := http.NewRequest("POST", addressPost, bytes.NewReader(compressJSON))
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер (1)")
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер (2)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.ErrInvalidLoginPassword
	}

	c.AuthorizedUser.User = user
	c.AuthorizedUser.Token = resp.Header.Get(constants.HeaderAuthorization)

	return nil
}

// inputPairLoginPassword событие формы, которое работает данными типа "пары логин/пароль"
func (c *Client) inputPairLoginPassword(plp postgresql.PairLoginPassword) error {
	addressPost := fmt.Sprintf("http://%s/api/resource/pairs", c.Config.Address) //a.cfg.Address)
	plpJSON, err := json.MarshalIndent(plp, "", " ")
	if err != nil {
		return err
	}
	_, err = ExecuteAPI(plpJSON, addressPost, c.Token)
	return err
}

// inputTextData событие формы, которое работают с данными типа "произвольные текстовые данные"
func (c *Client) inputTextData(td postgresql.TextData) error {
	addressPost := fmt.Sprintf("http://%s/api/resource/text", c.Config.Address) //a.cfg.Address)

	td.Text = encryption.EncryptString(td.Text, c.Config.CryptoKey)
	tdJSON, err := json.MarshalIndent(td, "", " ")
	if err != nil {
		return err
	}

	_, err = ExecuteAPI(tdJSON, addressPost, c.Token)
	return err
}

// inputBinaryData событие формы, которое работают с данными типа "произвольные бинарные данные"
func (c *Client) inputBinaryData(bd postgresql.BinaryData) error {
	addressPost := fmt.Sprintf("http://%s/api/resource/binary", c.Config.Address) //a.cfg.Address)

	bdJSON, err := json.MarshalIndent(bd, "", " ")
	if err != nil {
		return err
	}

	_, err = ExecuteAPI(bdJSON, addressPost, c.Token)

	if bd.Event != constants.EventDel.String() {
		ctx := context.Background()
		ctxWV := context.WithValue(ctx, postgresql.KeyContext("additionalBinaryParameters"), additionalBinaryParameters{
			patch: bd.Patch,
			uid:   bd.Uid,
		})
		go c.wsBinaryData(ctxWV)
	}

	return err
}

// inputBankCard событие формы, которое работают с данными типа "данные банковских карт"
func (c *Client) inputBankCard(bc postgresql.BankCard) error {
	addressPost := fmt.Sprintf("http://%s/api/resource/card", c.Config.Address)

	bc.Number = encryption.EncryptString(bc.Number, c.Config.CryptoKey)
	bc.Cvc = encryption.EncryptString(bc.Cvc, c.Config.CryptoKey)

	bcJSON, err := json.MarshalIndent(bc, "", " ")
	if err != nil {
		return err
	}

	_, err = ExecuteAPI(bcJSON, addressPost, c.Token)
	return err
}

// inputBankCard событие формы, которое работает с регистрацией нового пользователя
func (c *Client) registerNewUser(user postgresql.User) error {

	addressPost := fmt.Sprintf("http://%s/api/user/register", c.Config.Address) //a.cfg.Address)

	arrJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err
	}

	compressJSON, err := compression.Compress(arrJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	req, err := http.NewRequest("POST", addressPost, bytes.NewReader(compressJSON))
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер (1)")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер (2)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.ErrInvalidLoginPassword
	}

	c.AuthorizedUser.User = user
	c.AuthorizedUser.Token = resp.Header.Get(constants.HeaderAuthorization)

	return err
}

// ExecuteAPI общая фукция, которая сжимает в gzip, заполняет токены и отправляет на сервер данные,
// с которыми нужно произсести действия
func ExecuteAPI(bJSON []byte, addressPost, token string) (*http.Response, error) {
	compressJSON, err := compression.Compress(bJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return nil, err
	}

	body := bytes.NewReader(compressJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return nil, errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return nil, errors.New("-- ошибка отправки данных на сервер")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.ErrInvalidLoginPassword
	}

	return resp, nil
}
