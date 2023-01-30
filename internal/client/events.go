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

func (c *Client) creteEncryptionKey(k encryption.KeyRSA) error {
	if err := os.WriteFile(k.Patch, []byte(k.Key), 0664); err != nil {
		constants.Logger.ErrorLog(err)
	}
	c.Config.CryptoKey = k.Key
	return nil
}

func (c *Client) downloadBinaryData(bd postgresql.BinaryData) error {

	ctx := context.Background()
	ctxWV := context.WithValue(ctx, postgresql.KeyContext("additionalBinaryParameters"), additionalBinaryParameters{
		patch: bd.DownloadPatch,
		uid:   bd.Uid,
	})
	go c.wsDownloadBinaryData(ctxWV)
	return nil
}

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

func (c *Client) inputPairsLoginPassword(plp postgresql.PairsLoginPassword, event string) error {
	addressPost := fmt.Sprintf("http://%s/api/user/pairs/%s", c.Config.Address, event) //a.cfg.Address)

	arrJSON, err := json.MarshalIndent(plp, "", " ")
	if err != nil {
		return err
	}

	compressJSON, err := compression.Compress(arrJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	body := bytes.NewReader(compressJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.ErrInvalidLoginPassword
	}

	return nil
}

func (c *Client) inputTextData(td postgresql.TextData, event string) error {
	addressPost := fmt.Sprintf("http://%s/api/user/text/%s", c.Config.Address, event) //a.cfg.Address)

	encryptionTd := postgresql.TextData{}
	encryptionTd.User = c.Name
	encryptionTd.Uid = td.Uid
	encryptionTd.Text = encryption.EncryptString(td.Text, c.Config.CryptoKey)

	arrJSON, err := json.MarshalIndent(encryptionTd, "", " ")
	if err != nil {
		return err
	}

	compressJSON, err := compression.Compress(arrJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	body := bytes.NewReader(compressJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.ErrInvalidLoginPassword
	}

	return nil
}

func (c *Client) inputBinaryData(bd postgresql.BinaryData, event string) error {
	addressPost := fmt.Sprintf("http://%s/api/user/binary/%s", c.Config.Address, event) //a.cfg.Address)

	arrJSON, err := json.MarshalIndent(bd, "", " ")
	if err != nil {
		return err
	}

	compressJSON, err := compression.Compress(arrJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	body := bytes.NewReader(compressJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.ErrInvalidLoginPassword
	}

	if event != constants.EventDel.String() {
		ctx := context.Background()
		ctxWV := context.WithValue(ctx, postgresql.KeyContext("additionalBinaryParameters"), additionalBinaryParameters{
			patch: bd.Patch,
			uid:   bd.Uid,
		})
		go c.wsBinaryData(ctxWV)
	}
	return nil
}

func (c *Client) inputBankCard(bc postgresql.BankCard, event string) error {
	addressPost := fmt.Sprintf("http://%s/api/user/card/%s", c.Config.Address, event) //a.cfg.Address)

	bc.Number = encryption.EncryptString(bc.Number, c.Config.CryptoKey)
	bc.Cvc = encryption.EncryptString(bc.Cvc, c.Config.CryptoKey)

	arrJSON, err := json.MarshalIndent(bc, "", " ")
	if err != nil {
		return err
	}

	compressJSON, err := compression.Compress(arrJSON)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	body := bytes.NewReader(compressJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errs.ErrInvalidLoginPassword
	}

	return nil
}

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

	return nil
}
