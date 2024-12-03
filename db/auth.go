package db

import (
	"errors"

	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/supabase-community/gotrue-go/types"
)

func(db *Database) Signup(
  email, password string,
) error {
  var creds types.SignupRequest 
  creds.Email = email
  creds.Password = password

  _, err := db.supabase.Auth.Signup(creds)
  if err != nil {
    return err
  }

  return nil
}

// When an access Token is recovered. We can call this method to attempt a signin with the AccessToken.
// If successful, we update our local accesstoken and continue
func(db *Database) TokenSignin() error {
  if db.session == nil {
    return errors.New("User Session is not available")
  }
  if fakeLogin {
    utils.Warn(" ->> BYPASSING SESSION LOGIN")
    return nil
  }
  res, err := db.supabase.Auth.RefreshToken(
    db.session.RefreshToken,
  )

  if err != nil {
    return err
  }

  if err := db.session.UpdateCurrentSession(res); err != nil {
    return err
  }
  db.userID = res.User.ID.String()

  db.supabase.UpdateAuthSession(res.Session)
  db.supabase.EnableTokenAutoRefresh(res.Session)
  db.authChannel <- struct{}{}

  return nil
}

func(db *Database) Login(email, password string) error {
  res, err := db.supabase.Auth.SignInWithEmailPassword(email, password)
  if err != nil {
    return err
  }

  if err := db.session.UpdateCurrentSession(res); err != nil {
    return err
  }
  db.supabase.UpdateAuthSession(res.Session)
  db.supabase.EnableTokenAutoRefresh(res.Session)
  db.authChannel <- struct{}{}

  return nil
}

func(db *Database) Signout() error {
  if fakeLogin {
    return nil
  }

  if err := db.supabase.Auth.Logout(); err != nil {
    utils.Error("Failed to log user out: %w", err)
    return errors.New("Logout Failure")
  }
  if err := db.session.StoreLoggedoutSession(); err != nil {
    utils.Error("Failed to store Loggedout Session: %w", err)
    return errors.New("Store Logged Out Session Error")
  }

  return nil
}
