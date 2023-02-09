package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/theplant/luhn"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

// openRegisterForms отображает окно входа пользователя в систему
func (f *Forms) openLoginForm(c *Client) {

	user := postgresql.User{}
	f.Form.AddInputField("name", "", 20, nil, func(name string) {
		user.Name = name
	})
	f.Form.AddPasswordField("password", "", 20, ' ', func(password string) {
		user.Password = password
	})

	f.Form.AddButton("Login", func() {
		err := c.inputLoginUser(user)
		if err != nil {
			f.Form.AddTextView("", err.Error(), 100, 1, true, false)

			constants.Logger.ErrorLog(err)
			return
		}

		f.Pages.SwitchToPage(constants.NameMainPage)
	})

	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}

// openRegisterForms отображает окно регистрации нового пользователя
func (f *Forms) openRegisterForms(c *Client) {

	user := postgresql.User{}
	f.Form.AddInputField("name", "", 20, nil, func(name string) {
		user.Name = name
	})
	f.Form.AddPasswordField("password", "", 20, ' ', func(password string) {
		user.Password = password
	})

	f.Form.AddButton("Register new user", func() {
		err := c.registerNewUser(user)
		if err != nil {
			f.Form.AddTextView("", err.Error(), 100, 1, true, false)

			constants.Logger.ErrorLog(err)
			return
		}
		f.TextView.SetText(c.Name + "\n\n" + f.TextDefault)
		f.Pages.SwitchToPage(constants.NameMainPage)
	})

	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}

// openRegisterForms отображает окно всех данных текущего пользователя
func (f *Forms) openListForms(c *Client) {

	f.List.Clear()
	for k, v := range c.DataList {
		for _, val := range v {
			f.List.AddItem(k+":::"+val.MainText, val.SecondaryText, '*', nil).
				SetSelectedFunc(func(count int, mainText string, secondaryText string, rune rune) {
					arrMainText := strings.Split(mainText, ":::")
					switch arrMainText[0] {
					case constants.TypePairLoginPassword.String():

						arrSecondaryText := strings.Split(secondaryText, ":::")
						plp := postgresql.PairLoginPassword{
							Uid:      arrMainText[1],
							TypePair: arrSecondaryText[0],
							Name:     arrSecondaryText[1],
							Password: arrSecondaryText[2],
						}
						f.openPairLoginPasswordForms(c, plp)
						f.Pages.SwitchToPage("PairLoginPassword")

					case constants.TypeTextData.String():
						td := postgresql.TextData{
							Uid:  arrMainText[1],
							Text: secondaryText,
						}
						f.openTextDataForms(c, td)
						f.Pages.SwitchToPage("TextData")
					case constants.TypeBinaryData.String():
						arrSecondaryText := strings.Split(secondaryText, ":::")
						bd := postgresql.BinaryData{
							Uid:       arrMainText[1],
							Name:      arrSecondaryText[0],
							Expansion: arrSecondaryText[1],
							Size:      arrSecondaryText[2],
							Patch:     arrSecondaryText[3],
						}
						f.openBinaryDataForms(c, bd)
						f.Pages.SwitchToPage("BinaryData")
					case constants.TypeBankCardData.String():
						arrSecondaryText := strings.Split(secondaryText, ":::")
						bd := postgresql.BankCard{
							Uid:    arrMainText[1],
							Number: arrSecondaryText[0],
							Cvc:    arrSecondaryText[1],
						}
						f.openBankCardForms(c, bd)
						f.Pages.SwitchToPage("BinaryData")
					default:
						return
					}
				})
		}
	}
}

// openPairLoginPasswordForms отображает окно для ввода и действий данных типа "пары логин/пароль"
func (f *Forms) openPairLoginPasswordForms(c *Client, plp postgresql.PairLoginPassword) {

	if plp.Uid == "" {
		id := uuid.New()
		plp.Uid = id.String()
	}

	f.Form.AddTextView("UID:", plp.Uid, 36, 1, true, false)
	f.Form.AddInputField("type", plp.TypePair, 30, nil, func(TypePair string) {
		plp.TypePair = TypePair
	})
	f.Form.AddInputField("name", plp.Name, 30, nil, func(name string) {
		plp.Name = name
	})
	f.Form.AddPasswordField("password", plp.Password, 30, ' ', func(password string) {
		plp.Password = password
	})

	f.Form.AddButton("Add/edit login/password pairs", func() {
		plp.TypePair = encryption.EncryptString(plp.TypePair, c.Config.CryptoKey)
		plp.Name = encryption.EncryptString(plp.Name, c.Config.CryptoKey)
		plp.Password = encryption.EncryptString(plp.Password, c.Config.CryptoKey)
		plp.Event = constants.EventAddEdit.String()

		err := c.inputPairLoginPassword(plp)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Delete login/password pairs", func() {
		plp.TypePair = encryption.EncryptString(plp.TypePair, c.Config.CryptoKey)
		plp.Name = encryption.EncryptString(plp.Name, c.Config.CryptoKey)
		plp.Password = encryption.EncryptString(plp.Password, c.Config.CryptoKey)
		plp.Event = constants.EventDel.String()

		err := c.inputPairLoginPassword(plp)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}

// openTextDataForms отображает окно для ввода и действий данных типа "произвольные текстовые данные"
func (f *Forms) openTextDataForms(c *Client, td postgresql.TextData) {

	td.User = c.Token

	if td.Uid == "" {
		id := uuid.New()
		td.Uid = id.String()
	}

	f.Form.AddTextView("UID:", td.Uid, 36, 1, true, false)
	f.Form.AddTextArea("text", td.Text, 200, 10, 15000, func(text string) {
		td.Text = text
	})

	f.Form.AddButton("Add/edit text", func() {
		td.Event = constants.EventAddEdit.String()

		err := c.inputTextData(td)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Delete text", func() {
		td.Event = constants.EventDel.String()

		err := c.inputTextData(td)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}

// openBinaryDataForms отображает окно для ввода и действий данных типа "произвольные бинарные данные"
func (f *Forms) openBinaryDataForms(c *Client, bd postgresql.BinaryData) {

	bd.User = c.Token

	if bd.Uid == "" {
		id := uuid.New()
		bd.Uid = id.String()
	}

	f.Form.AddTextView("UID:", bd.Uid, 36, 1, true, false)
	f.Form.AddInputField("Patch:", bd.Patch, 200, nil, func(patch string) {
		bd.Patch = patch
	})
	f.Form.AddInputField("Download patch:", bd.Patch, 200, nil, func(patch string) {
		bd.DownloadPatch = patch
	})
	f.Form.AddTextView("Name:", bd.Name, 50, 1, true, false)
	f.Form.AddTextView("Expansion:", bd.Expansion, 50, 1, true, false)
	f.Form.AddTextView("Size:", fmt.Sprintf("%v", bd.Size), 50, 1, true, false)

	f.Form.AddButton("Upload binary", func() {
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
		bd.Event = constants.EventAddEdit.String()

		err = c.inputBinaryData(bd)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Download binary", func() {
		if bd.DownloadPatch == "" {
			fmt.Println("не указан путь к файлу")
			return
		}

		err := c.downloadBinaryData(bd)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Delete binary", func() {
		bd.Event = constants.EventDel.String()

		err := c.inputBinaryData(bd)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}

// openBankCardForms отображает окно для ввода и действий данных типа "данные банковских карт"
func (f *Forms) openBankCardForms(c *Client, bc postgresql.BankCard) {

	bc.User = c.Token

	if bc.Uid == "" {
		id := uuid.New()
		bc.Uid = id.String()
	}

	f.Form.AddTextView("UID:", bc.Uid, 36, 1, true, false)
	f.Form.AddInputField("Number:", bc.Number, 30, nil, func(number string) {
		bc.Number = number
	})
	f.Form.AddInputField("CVC", bc.Cvc, 30, nil, func(cvc string) {
		bc.Cvc = cvc
	})

	f.Form.AddButton("Add/edit bank card", func() {
		numCard, err := strconv.Atoi(bc.Number)
		if !luhn.Valid(numCard) {
			constants.Logger.ErrorLog(err)
			return
		}

		bc.Event = constants.EventAddEdit.String()
		err = c.inputBankCard(bc)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Delete binary", func() {
		bc.Event = constants.EventDel.String()
		err := c.inputBankCard(bc)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}

// openEncryptionKeyForms отображает окно для ввода информации для формирования ключа для шифровки данных
func (f *Forms) openEncryptionKeyForms(c *Client, k encryption.KeyRSA) {

	k.User = c.Token

	f.Form.AddTextArea("Key:", k.Key, 200, 10, 15000, func(key string) {
		k.Key = key
	})
	f.Form.AddInputField("Patch:", k.Patch, 100, nil, func(patch string) {
		k.Patch = patch
	})

	f.Form.AddButton("Create key", func() {
		if k.Patch == "" {
			fmt.Println("Не указан путь к файлу")
			return
		}
		err := c.createEncryptionKey(k)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
	f.Form.AddButton("Cancel", func() {
		f.Pages.SwitchToPage(constants.NameMainPage)
	})
}
