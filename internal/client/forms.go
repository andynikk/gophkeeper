package client

import (
	"gophkeeper/internal/encryption"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rivo/tview"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/postgresql"
)

func (c *Client) loginForm() {

	user := postgresql.User{}
	form.AddInputField("name", "", 20, nil, func(name string) {
		user.Name = name
	})
	form.AddPasswordField("password", "", 20, ' ', func(password string) {
		user.Password = password
	})

	form.AddButton("Login", func() {
		err := c.loginUser(user)
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

func (c *Client) registerForms() {

	user := postgresql.User{}
	form.AddInputField("name", "", 20, nil, func(name string) {
		user.Name = name
	})
	form.AddPasswordField("password", "", 20, ' ', func(password string) {
		user.Password = password
	})

	form.AddButton("Register new user", func() {
		err := c.registerUser(user)
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

func (c *Client) listForms(list *tview.List) {

	res, err := os.ReadFile("e:\\Bases\\key\\gophkeeper.xor")
	if err != nil {
		constants.Logger.ErrorLog(err)
		return
	}

	list.Clear()
	for k, v := range c.DataList {
		for _, val := range v {
			secondaryText0 := ""
			secondaryText1 := ""
			secondaryText, err := encryption.DecryptString(val.SecondaryText, string(res))
			if err != nil {
				constants.Logger.ErrorLog(err)
				secondaryText = val.SecondaryText
			}
			if strings.Contains(val.SecondaryText, ":::") {
				arrSecondaryText := strings.Split(val.SecondaryText, ":::")
				secondaryText0, err = encryption.DecryptString(arrSecondaryText[0], string(res))
				if err != nil {
					constants.Logger.ErrorLog(err)
					secondaryText0 = arrSecondaryText[0]
				}
				secondaryText1, err = encryption.DecryptString(arrSecondaryText[1], string(res))
				if err != nil {
					constants.Logger.ErrorLog(err)
					secondaryText1 = arrSecondaryText[1]
				}
				secondaryText = secondaryText0 + ":::" + secondaryText1
			}

			list.AddItem(k+":::"+val.MainText, secondaryText, '*', nil).
				SetSelectedFunc(func(count int, mainText string, secondaryText string, rune rune) {
					arrMainText := strings.Split(mainText, ":::")
					switch arrMainText[0] {
					case constants.TypePairsLoginPassword.String():

						arrSecondaryText := strings.Split(secondaryText, ":::")

						name := arrSecondaryText[0]
						typePairs := arrMainText[1]
						password := arrSecondaryText[1]

						plp := postgresql.PairsLoginPassword{
							Name:      name,
							TypePairs: typePairs,
							Password:  password,
						}
						c.pairsLoginPasswordForms(plp)
						pages.SwitchToPage("PairsLoginPassword")

					case constants.TypeTextData.String():
						td := postgresql.TextData{
							Uid:  arrMainText[1],
							Text: secondaryText,
						}
						c.textDataForms(td)
						pages.SwitchToPage("TextData")
					default:
						return
					}
				})
		}
	}
}

func (c *Client) pairsLoginPasswordForms(plp postgresql.PairsLoginPassword) {

	plp.User = c.Name

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
		err := c.pairsLoginPassword(plp, "edit")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Delete login/password pairs", func() {
		err := c.pairsLoginPassword(plp, "del")
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

func (c *Client) textDataForms(td postgresql.TextData) {

	td.User = c.Name

	if td.Uid == "" {
		id := uuid.New()
		td.Uid = id.String()
	}
	form.AddInputField("UID", td.Uid, 36, nil, func(Uid string) {
		td.Uid = Uid
	})
	form.AddTextArea("text", td.Text, 200, 10, 15000, func(text string) {
		td.Text = text
	})

	form.AddButton("Add/edit text", func() {
		err := c.textData(td, "edit")
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		text.SetText(c.setMainText())
		pages.SwitchToPage("Menu")
	})
	form.AddButton("Delete text", func() {
		err := c.textData(td, "del")
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

func (c *Client) keyRSAForms(k encryption.KeyRSA) {

	k.User = c.Name
	if k.Patch == "" {
		k.Patch = c.Config.CryptoKey
	}
	form.AddInputField("Patch", k.Patch, 100, nil, func(patch string) {
		k.Patch = patch
	})
	form.AddInputField("Number sert.", string(k.NumSert), 20, nil, func(numSert string) {
		res, err := strconv.ParseInt(numSert, 10, 64)
		if err != nil {
			k.NumSert = 0
		} else {
			k.NumSert = res
		}
	})
	form.AddInputField("Subject key ID", k.SubjectKeyID, 20, nil, func(subjectKeyID string) {
		k.SubjectKeyID = subjectKeyID
	})
	form.AddInputField("Len key byte.", string(k.LenKeyByte), 20, nil, func(lenKeyByte string) {
		res, err := strconv.ParseInt(lenKeyByte, 0, 64)
		if err != nil {
			k.LenKeyByte = 0
		} else {
			k.LenKeyByte = int(res)
		}
	})

	form.AddButton("Create file", func() {
		err := c.keyRSA(k)
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
