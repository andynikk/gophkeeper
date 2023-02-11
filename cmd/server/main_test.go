package main

import (
	"context"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/cryptography"
	"gophkeeper/internal/handlers"
	"gophkeeper/internal/postgresql"
	"gophkeeper/internal/postgresql/model"
	"gophkeeper/internal/tests"
	"gophkeeper/internal/token"
	"testing"
)

func TestFuncClient(t *testing.T) {

	srv := &handlers.Server{}
	t.Run("Checking init config", func(t *testing.T) {
		srv.InitConfig()
		if srv.ServerConfig.Address == "" {
			t.Errorf("Error init config")
		}
	})

	t.Run("Checking init DB", func(t *testing.T) {
		srv.InitDataBase()
		if srv.DatabaseDsn == "" {
			t.Errorf("Error init DB")
		}
	})

	t.Run("Checking init Router", func(t *testing.T) {
		srv.InitRouters()
		if srv.Router == nil {
			t.Errorf("Error init Router")
		}
	})

	srv.InListUserData = handlers.InListUserData{}

	t.Run("Checking init server", func(t *testing.T) {
		if srv.ServerConfig.Address == "" {
			t.Errorf("Error checking init server")
		}
	})

	t.Run("Checking init router", func(t *testing.T) {
		if srv.Router == nil {
			t.Errorf("Error checking init router")
		}
	})

	t.Run("Checking init config", func(t *testing.T) {
		if srv.ServerConfig.Address == "" {
			t.Errorf("Error checking init config")
		}
	})

	tc := token.NewClaims("test")
	strToken, _ := tc.GenerateJWT()
	ck := "test crypto key"
	//claims, _ := token.ExtractClaims(strToken)

	t.Run("Checking DB", func(t *testing.T) {
		t.Run("Checking create DB table", func(t *testing.T) {
			if srv.DatabaseDsn != "" {
				t.Run("Checking ping DB", func(t *testing.T) {
					err := postgresql.CreateModeLDB(srv.Pool)
					if err != nil {
						t.Errorf("Error handlers ping DB")
					}
					t.Run("Checking create user DB user", func(t *testing.T) {
						user := tests.CreateUser("")
						err = srv.DBConnector.Update(&user)
						if err != nil {
							t.Errorf("Error create user DB user")
						}

						err = srv.DBConnector.CheckAccount(&user)
						if err != nil {
							t.Errorf("Error create user DB user")
						}
					})
					t.Run("Checking re-creation user DB user", func(t *testing.T) {
						user := tests.CreateUser("")
						err = srv.DBConnector.NewAccount(&user)
						if err == nil {
							t.Errorf("Error re-creation user DB user")
						}
					})
					t.Run("Checking delete DB user", func(t *testing.T) {
						user := tests.CreateUser("")
						user.HashPassword = cryptography.HashSHA256(user.Password, srv.Key)
						err = srv.DBConnector.DelAccount(&user)
						if err != nil {
							t.Errorf("Error delete ping DB")
						}

						err = srv.DBConnector.CheckAccount(&user)
						if err == nil {
							t.Errorf("Error delete user DB user")
						}
					})

					t.Run("Checking Pairs login/password DB", func(t *testing.T) {
						plp := tests.CreatePairLoginPassword(strToken, "", ck)
						t.Run("Checking update Pairs login/password DB", func(t *testing.T) {
							err = srv.DBConnector.Update(&plp)
							if err != nil {
								t.Errorf("Error Pairs login/password DB")
							}
						})
						t.Run("Checking select Pairs login/password DB", func(t *testing.T) {
							ctx := context.Background()
							ctxWV := context.WithValue(ctx, model.KeyContext("user"), strToken)
							arrPlp, err := srv.DBConnector.Select(ctxWV, constants.TypePairLoginPassword.String())
							if err != nil || len(arrPlp) == 0 {
								t.Errorf("Error select Pairs login/password DB")
							}
						})
						t.Run("Checking delete Pairs login/password DB", func(t *testing.T) {
							err := srv.DBConnector.Delete(&plp)
							if err != nil {
								t.Errorf("Error delete Pairs login/password DB")
							}
						})
					})

					t.Run("Checking Text data DB", func(t *testing.T) {
						td := tests.CreateTextData(strToken, "", ck)
						t.Run("Checking update Text data DB", func(t *testing.T) {
							err = srv.DBConnector.Update(&td)
							if err != nil {
								t.Errorf("Error update Text data DB")
							}
						})
						t.Run("Checking select Text data DB", func(t *testing.T) {
							ctx := context.Background()
							ctxWV := context.WithValue(ctx, model.KeyContext("user"), strToken)
							arrPlp, err := srv.DBConnector.Select(ctxWV, constants.TypeTextData.String())
							if err != nil || len(arrPlp) == 0 {
								t.Errorf("Error select Text data DB")
							}
						})
						t.Run("Checking delete Text data DB", func(t *testing.T) {
							err := srv.DBConnector.Delete(&td)
							if err != nil {
								t.Errorf("Error delete Text data DB")
							}
						})
					})

					t.Run("Checking Binary data DB", func(t *testing.T) {
						bd := tests.CreateBinaryData(strToken, "")
						t.Run("Checking update Binary data DB", func(t *testing.T) {
							err = srv.DBConnector.Update(&bd)
							if err != nil {
								t.Errorf("Error handlers ping DB")
							}
						})
						t.Run("Checking select Binary data DB", func(t *testing.T) {
							ctx := context.Background()
							ctxWV := context.WithValue(ctx, model.KeyContext("user"), strToken)
							arrPlp, err := srv.DBConnector.Select(ctxWV, constants.TypeBinaryData.String())
							if err != nil || len(arrPlp) == 0 {
								t.Errorf("Error select text data DB")
							}
						})
						t.Run("Checking delete Binary data DB", func(t *testing.T) {
							err := srv.DBConnector.Delete(&bd)
							if err != nil {
								t.Errorf("Error delete text data DB")
							}
						})
					})

					t.Run("Checking Bank data DB", func(t *testing.T) {
						bd := tests.CreateBankCard(strToken, "", ck)
						t.Run("Checking update Bank data DB", func(t *testing.T) {
							err = srv.DBConnector.Update(&bd)
							if err != nil {
								t.Errorf("Error update Bank data DB")
							}
						})
						t.Run("Checking select Bank data DB", func(t *testing.T) {
							ctx := context.Background()
							ctxWV := context.WithValue(ctx, model.KeyContext("user"), strToken)
							arrPlp, err := srv.DBConnector.Select(ctxWV, constants.TypeBankCardData.String())
							if err != nil || len(arrPlp) == 0 {
								t.Errorf("Error select Bank data DB")
							}
						})
						t.Run("Checking delete Bank data DB", func(t *testing.T) {
							err := srv.DBConnector.Delete(&bd)
							if err != nil {
								t.Errorf("Error delete Bank data DB")
							}
						})
					})

				})
			}
		})
	})

	ctx := context.Background()
	conn, err := srv.Pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()
	pc := postgresql.PgxpoolConn{
		Conn: conn,
	}

	t.Run("Checking methods Pair login/password", func(t *testing.T) {
		plp := tests.CreatePairLoginPassword(strToken, "", ck)
		t.Run("Checking method 'CheckExistence' Pair login/password", func(t *testing.T) {

			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &plp)
			_, err = pc.CheckExistence(ctxVW)
			if err != nil {
				t.Errorf("Error method 'CheckExistence' Bank card DB")
			}
		})

		t.Run("Checking method 'Insert' Pair login/password", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &plp)
			err = pc.Insert(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Pair login/password DB")
			}
		})

		t.Run("Checking method 'Update' Pair login/password", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &plp)
			err = pc.Update(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Pair login/password DB")
			}
		})

		t.Run("Checking delete Pairs login/password DB", func(t *testing.T) {
			err := srv.DBConnector.Delete(&plp)
			if err != nil {
				t.Errorf("Error delete Pairs login/password DB")
			}
		})

		t.Run("Checking method 'GetType' Pair login/password", func(t *testing.T) {
			tp := plp.GetType()
			if tp == "" {
				t.Errorf("Error method 'GetType' Pair login/password DB")
			}
		})

		t.Run("Checking method 'GetMainText' Pair login/password", func(t *testing.T) {
			mt := plp.GetMainText()
			if mt == "" {
				t.Errorf("Error method 'GetMainText' Pair login/password DB")
			}
		})

		t.Run("Checking method 'GetSecondaryText' Pair login/password", func(t *testing.T) {
			st := plp.GetSecondaryText(ck)
			if st == "" {
				t.Errorf("Error method 'GetSecondaryText' Pair login/password DB")
			}
		})

		t.Run("Checking method 'SetFromInListUserData' Pair login/password", func(t *testing.T) {
			plpInListUserData, ok := srv.InListUserData[constants.TypePairLoginPassword.String()]
			if !ok {
				plpInListUserData = model.Appender{}
			}
			plp.SetFromInListUserData(plpInListUserData)

			if plpInListUserData[plp.Uid] == nil {
				t.Errorf("Error method 'SetFromInListUserData' Pair login/password DB")
			}
		})
	})

	t.Run("Checking methods Text data", func(t *testing.T) {
		td := tests.CreateTextData(strToken, "", ck)
		t.Run("Checking method 'CheckExistence' Text data", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &td)
			_, err = pc.CheckExistence(ctxVW)
			if err != nil {
				t.Errorf("Error method 'CheckExistence' Bank card DB")
			}
		})

		t.Run("Checking method 'Insert' Text data", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &td)
			err = pc.Insert(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Text data DB")
			}
		})

		t.Run("Checking method 'Update' Text data", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &td)
			err = pc.Update(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Text data DB")
			}
		})

		t.Run("Checking delete Text data DB", func(t *testing.T) {
			err := srv.DBConnector.Delete(&td)
			if err != nil {
				t.Errorf("Error delete Text data DB")
			}
		})

		t.Run("Checking method 'GetType' Text data", func(t *testing.T) {
			tp := td.GetType()
			if tp == "" {
				t.Errorf("Error method 'GetType' Text data DB")
			}
		})

		t.Run("Checking method 'GetMainText' Text data", func(t *testing.T) {
			mt := td.GetMainText()
			if mt == "" {
				t.Errorf("Error method 'GetMainText' Text data")
			}
		})

		t.Run("Checking method 'GetSecondaryText' Text data", func(t *testing.T) {
			st := td.GetSecondaryText(ck)
			if st == "" {
				t.Errorf("Error method 'GetSecondaryText' Text data")
			}
		})

		t.Run("Checking method 'SetFromInListUserData' Text data", func(t *testing.T) {
			plpInListUserData, ok := srv.InListUserData[constants.TypeTextData.String()]
			if !ok {
				plpInListUserData = model.Appender{}
			}
			td.SetFromInListUserData(plpInListUserData)

			if plpInListUserData[td.Uid] == nil {
				t.Errorf("Error method 'SetFromInListUserData' Text data DB")
			}
		})
	})

	t.Run("Checking methods Binary data", func(t *testing.T) {
		bd := tests.CreateBinaryData(strToken, "")
		t.Run("Checking method 'CheckExistence' Binary data", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &bd)
			_, err = pc.CheckExistence(ctxVW)
			if err != nil {
				t.Errorf("Error method 'CheckExistence' Bank card DB")
			}
		})

		t.Run("Checking method 'Insert' Binary data", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &bd)
			err = pc.Insert(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Binary data DB")
			}
		})

		t.Run("Checking method 'Update' Binary data", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &bd)
			err = pc.Update(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Binary data DB")
			}
		})

		t.Run("Checking delete Binary data DB", func(t *testing.T) {
			err := srv.DBConnector.Delete(&bd)
			if err != nil {
				t.Errorf("Error delete Binary data DB")
			}
		})

		t.Run("Checking method 'GetType' Binary data", func(t *testing.T) {
			tp := bd.GetType()
			if tp == "" {
				t.Errorf("Error method 'GetType' Binary data DB")
			}
		})

		t.Run("Checking method 'GetMainText' Binary data", func(t *testing.T) {
			mt := bd.GetMainText()
			if mt == "" {
				t.Errorf("Error method 'GetMainText' Binary data")
			}
		})

		t.Run("Checking method 'GetSecondaryText' Binary data", func(t *testing.T) {
			st := bd.GetSecondaryText(ck)
			if st == "" {
				t.Errorf("Error method 'GetSecondaryText' Binary data")
			}
		})

		t.Run("Checking method 'SetFromInListUserData' Binary data", func(t *testing.T) {
			plpInListUserData, ok := srv.InListUserData[constants.TypeBinaryData.String()]
			if !ok {
				plpInListUserData = model.Appender{}
			}
			bd.SetFromInListUserData(plpInListUserData)

			if plpInListUserData[bd.Uid] == nil {
				t.Errorf("Error method 'SetFromInListUserData' Binary data DB")
			}
		})
	})

	t.Run("Checking methods Bank card", func(t *testing.T) {
		bc := tests.CreateBankCard(strToken, "", ck)
		t.Run("Checking method 'CheckExistence' Bank card", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &bc)
			_, err = pc.CheckExistence(ctxVW)
			if err != nil {
				t.Errorf("Error method 'CheckExistence' Bank card DB")
			}
		})

		t.Run("Checking method 'Insert' Bank card", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &bc)
			err = pc.Insert(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Bank card DB")
			}
		})

		t.Run("Checking method 'Update' Bank card", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &bc)
			err = pc.Update(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' Bank card DB")
			}
		})

		t.Run("Checking delete Bank card DB", func(t *testing.T) {
			err := srv.DBConnector.Delete(&bc)
			if err != nil {
				t.Errorf("Error delete Bank card DB")
			}
		})

		t.Run("Checking method 'GetType' Bank card", func(t *testing.T) {
			tp := bc.GetType()
			if tp == "" {
				t.Errorf("Error method 'GetType' Bank card DB")
			}
		})

		t.Run("Checking method 'GetMainText' Binary data", func(t *testing.T) {
			mt := bc.GetMainText()
			if mt == "" {
				t.Errorf("Error method 'GetMainText' Bank card")
			}
		})

		t.Run("Checking method 'GetSecondaryText' Bank card", func(t *testing.T) {
			st := bc.GetSecondaryText(ck)
			if st == "" {
				t.Errorf("Error method 'GetSecondaryText' Bank card")
			}
		})

		t.Run("Checking method 'SetFromInListUserData' Bank card", func(t *testing.T) {
			plpInListUserData, ok := srv.InListUserData[constants.TypeBankCardData.String()]
			if !ok {
				plpInListUserData = model.Appender{}
			}
			bc.SetFromInListUserData(plpInListUserData)

			if plpInListUserData[bc.Uid] == nil {
				t.Errorf("Error method 'SetFromInListUserData' Bank card DB")
			}
		})
	})

	t.Run("Checking methods User", func(t *testing.T) {
		usr := tests.CreateUser("")
		t.Run("Checking method 'Insert' User", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &usr)
			err = pc.Insert(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Insert' User DB")
			}
		})

		t.Run("Checking method 'CheckExistence' User", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &usr)
			_, err = pc.CheckExistence(ctxVW)
			if err != nil {
				t.Errorf("Error method 'CheckExistence' Bank card DB")
			}
		})

		t.Run("Checking method 'Delete' User", func(t *testing.T) {
			ctxVW := context.WithValue(ctx, model.KeyContext("data"), &usr)
			err = pc.Delete(ctxVW)
			if err != nil {
				t.Errorf("Error method 'Delete' User DB")
			}
		})

		t.Run("Checking method 'SetFromInListUserData' User", func(t *testing.T) {
			plpInListUserData, ok := srv.InListUserData[constants.TypeUserData.String()]
			if !ok {
				plpInListUserData = model.Appender{}
			}
			usr.SetFromInListUserData(plpInListUserData)

			if plpInListUserData[usr.Name] == nil {
				t.Errorf("Error method 'SetFromInListUserData' User DB")
			}
		})
	})

}
