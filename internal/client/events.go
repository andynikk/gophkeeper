package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/encryption"
	"net/http"
	"os"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/postgresql"
)

func (c *Client) loginUser(user postgresql.User) error {
	addressPost := fmt.Sprintf("http://%s/api/user/login", "localhost:8080") //a.cfg.Address)

	arrJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", addressPost, bytes.NewReader(arrJSON))
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер (1)")
	}

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Encoding", "gzip")
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

func (c *Client) registerUser(user postgresql.User) error {
	addressPost := fmt.Sprintf("http://%s/api/user/register", "localhost:8080") //a.cfg.Address)

	arrJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", addressPost, bytes.NewReader(arrJSON))
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер (1)")
	}

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Encoding", "gzip")
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

func (c *Client) pairsLoginPassword(plp postgresql.PairsLoginPassword, event string) error {
	addressPost := fmt.Sprintf("http://%s/api/user/pairs/"+event, "localhost:8080") //a.cfg.Address)

	arrJSON, err := json.MarshalIndent(plp, "", " ")
	if err != nil {
		return err
	}

	body := bytes.NewReader(arrJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
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

func (c *Client) textData(td postgresql.TextData, event string) error {
	addressPost := fmt.Sprintf("http://%s/api/user/text/"+event, "localhost:8080") //a.cfg.Address)

	res, err := os.ReadFile("e:\\Bases\\key\\gophkeeper.xor")
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}

	encryptionTd := postgresql.TextData{}
	encryptionTd.User = c.Name
	encryptionTd.Uid = td.Uid
	encryptionTd.Text, _ = encryption.EncryptString(td.Text, string(res))

	arrJSON, err := json.MarshalIndent(encryptionTd, "", " ")
	if err != nil {
		return err
	}

	body := bytes.NewReader(arrJSON)
	req, err := http.NewRequest("POST", addressPost, body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return errors.New("-- ошибка отправки данных на сервер")
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
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

func (c *Client) keyRSA(k encryption.KeyRSA) error {
	arrCert, err := k.CreateCert()
	if err != nil {
		constants.Logger.ErrorLog(err)
		return err
	}
	encryption.SaveKeyInFile(&arrCert[1], k.Patch)
	return nil
}
