package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/theplant/luhn"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

func (c *Client) openLoginForm() {

	user := postgresql.User{}
	form.AddInputField("name", "", 20, nil, func(name string) {
		user.Name = name
	})
	form.AddPasswordField("password", "", 20, ' ', func(password string) {
		user.Password = password
	})

	form.AddButton("Login", func() {
		err := c.inputLoginUser(user)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}

		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})

	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}

func (c *Client) openRegisterForms() {

	user := postgresql.User{}
	form.AddInputField("name", "", 20, nil, func(name string) {
		user.Name = name
	})
	form.AddPasswordField("password", "", 20, ' ', func(password string) {
		user.Password = password
	})

	form.AddButton("Register new user", func() {
		err := c.registerNewUser(user)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.Name + "\n\n" + textDefault)
		pages.SwitchToPage("Menu")
	})

	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}

func (c *Client) openListForms(list *tview.List) {

	list.Clear()
	for k, v := range c.DataList {
		for _, val := range v {
			list.AddItem(k+":::"+val.MainText, val.SecondaryText, '*', nil).
				SetSelectedFunc(func(count int, mainText string, secondaryText string, rune rune) {
					arrMainText := strings.Split(mainText, ":::")
					switch arrMainText[0] {
					case constants.TypePairsLoginPassword.String():

						arrSecondaryText := strings.Split(secondaryText, ":::")
						plp := postgresql.PairsLoginPassword{
							Uid:       arrMainText[1],
							TypePairs: arrSecondaryText[0],
							Name:      arrSecondaryText[1],
							Password:  arrSecondaryText[2],
						}
						c.openPairsLoginPasswordForms(plp)
						pages.SwitchToPage("PairsLoginPassword")

					case constants.TypeTextData.String():
						td := postgresql.TextData{
							Uid:  arrMainText[1],
							Text: secondaryText,
						}
						c.openTextDataForms(td)
						pages.SwitchToPage("TextData")
					case constants.TypeBinaryData.String():
						arrSecondaryText := strings.Split(secondaryText, ":::")
						bd := postgresql.BinaryData{
							Uid:       arrMainText[1],
							Name:      arrSecondaryText[0],
							Expansion: arrSecondaryText[1],
							Size:      arrSecondaryText[2],
							Patch:     arrSecondaryText[3],
						}
						c.openBinaryDataForms(bd)
						pages.SwitchToPage("BinaryData")
					case constants.TypeBankCardData.String():
						arrSecondaryText := strings.Split(secondaryText, ":::")
						bd := postgresql.BankCard{
							Uid:    arrMainText[1],
							Number: arrSecondaryText[0],
							Cvc:    arrSecondaryText[1],
						}
						c.openBankCardForms(bd)
						pages.SwitchToPage("BinaryData")
					default:
						return
					}
				})
		}
	}
}

func (c *Client) openPairsLoginPasswordForms(plp postgresql.PairsLoginPassword) {

	plp.User = c.Name

	if plp.Uid == "" {
		id := uuid.New()
		plp.Uid = id.String()
	}
	form.AddTextView("UID:", plp.Uid, 36, 1, true, false)
	form.AddInputField("type", plp.TypePairs, 30, nil, func(typePairs string) {
		plp.TypePairs = typePairs
	})
	form.AddInputField("name", plp.Name, 30, nil, func(name string) {
		plp.Name = name
	})
	form.AddPasswordField("password", plp.Password, 30, ' ', func(password string) {
		plp.Password = password
	})

	form.AddButton("Add/edit login/password pairs", func() {
		plp.TypePairs = encryption.EncryptString(plp.TypePairs, c.Config.CryptoKey)
		plp.Name = encryption.EncryptString(plp.Name, c.Config.CryptoKey)
		plp.Password = encryption.EncryptString(plp.Password, c.Config.CryptoKey)

		err := c.inputPairsLoginPassword(plp, "edit")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Delete login/password pairs", func() {
		plp.TypePairs = encryption.EncryptString(plp.TypePairs, c.Config.CryptoKey)
		plp.Name = encryption.EncryptString(plp.Name, c.Config.CryptoKey)
		plp.Password = encryption.EncryptString(plp.Password, c.Config.CryptoKey)

		err := c.inputPairsLoginPassword(plp, "del")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}

func (c *Client) openTextDataForms(td postgresql.TextData) {

	td.User = c.Name

	if td.Uid == "" {
		id := uuid.New()
		td.Uid = id.String()
	}
	form.AddTextView("UID:", td.Uid, 36, 1, true, false)
	form.AddTextArea("text", td.Text, 200, 10, 15000, func(text string) {
		td.Text = text
	})

	form.AddButton("Add/edit text", func() {
		err := c.inputTextData(td, "edit")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Delete text", func() {
		err := c.inputTextData(td, "del")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}

func (c *Client) openBinaryDataForms(bd postgresql.BinaryData) {

	bd.User = c.Name

	if bd.Uid == "" {
		id := uuid.New()
		bd.Uid = id.String()
	}

	form.AddTextView("UID:", bd.Uid, 36, 1, true, false)
	form.AddInputField("Patch:", bd.Patch, 200, nil, func(patch string) {
		bd.Patch = patch
	})
	form.AddInputField("Download patch:", bd.Patch, 200, nil, func(patch string) {
		bd.DownloadPatch = patch
	})
	form.AddTextView("Name:", bd.Name, 50, 1, true, false)
	form.AddTextView("Expansion:", bd.Expansion, 50, 1, true, false)
	form.AddTextView("Size:", fmt.Sprintf("%v", bd.Size), 50, 1, true, false)

	form.AddButton("Upload binary", func() {
		if bd.Patch == "" {
			fmt.Println("не указан путь к файлу")
			return
		}
		fileInfo, err := os.Stat(bd.Patch)
		if fileInfo == nil || err != nil {
			fmt.Println("по указанному пути, файл не найден")
			return
		}
		bd.Name = fileInfo.Name()
		bd.Expansion = "pdf"
		bd.Size = fmt.Sprintf("%d", fileInfo.Size())

		err = c.inputBinaryData(bd, "edit")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Download binary", func() {
		if bd.DownloadPatch == "" {
			fmt.Println("не указан путь к файлу")
			return
		}

		err := c.downloadBinaryData(bd)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Delete binary", func() {
		err := c.inputBinaryData(bd, "del")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}

func (c *Client) openBankCardForms(bc postgresql.BankCard) {

	bc.User = c.Name

	if bc.Uid == "" {
		id := uuid.New()
		bc.Uid = id.String()
	}

	form.AddTextView("UID:", bc.Uid, 36, 1, true, false)
	form.AddInputField("Number:", bc.Number, 30, nil, func(number string) {
		bc.Number = number
	})
	form.AddInputField("CVC", bc.Cvc, 30, nil, func(cvc string) {
		bc.Cvc = cvc
	})

	form.AddButton("Add/edit bank card", func() {
		numCard, err := strconv.Atoi(bc.Number)
		if !luhn.Valid(numCard) {
			constants.Logger.ErrorLog(err)
			return
		}

		err = c.inputBankCard(bc, "edit")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Delete binary", func() {
		err := c.inputBankCard(bc, "del")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}

func (c *Client) openEncryptionKeyForms(k encryption.KeyRSA) {

	k.User = c.Name

	form.AddTextArea("Key:", k.Key, 200, 10, 15000, func(key string) {
		k.Key = key
	})
	form.AddInputField("Patch:", k.Patch, 100, nil, func(patch string) {
		k.Patch = patch
	})

	form.AddButton("Create key", func() {
		if k.Patch == "" {
			fmt.Println("Не указан путь к файлу")
			return
		}
		err := c.creteEncryptionKey(k)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Cancel", func() {
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
}
