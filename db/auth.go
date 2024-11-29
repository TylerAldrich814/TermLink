package db

import (
	"errors"

	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/supabase-community/gotrue-go/types"
)

func(db *Supabase) Signup(
  email, password string,
) error {
  var creds types.SignupRequest 
  creds.Email = email
  creds.Password = password

  _, err := db.client.Auth.Signup(creds)
  if err != nil {
    return err
  }

  return nil
}

// When an access Token is recovered. We can call this method to attempt a signin with the AccessToken.
// If successful, we update our local accesstoken and continue
func(db *Supabase) TokenSignin() error {
  res, err := db.client.Auth.RefreshToken(
    db.session.AccessToken,
  )
  if err != nil {
    return err
  }

  if err := db.session.UpdateCurrentSession(res); err != nil {
    return err
  }

  db.client.UpdateAuthSession(res.Session)
  db.client.EnableTokenAutoRefresh(res.Session)
  db.authChannel <- struct{}{}

  return nil
}

func(db *Supabase) Login(email, password string) error {
  res, err := db.client.Auth.SignInWithEmailPassword(email, password)
  if err != nil {
    return err
  }

  if err := db.session.UpdateCurrentSession(res); err != nil {
    return err
  }
  db.client.UpdateAuthSession(res.Session)
  db.client.EnableTokenAutoRefresh(res.Session)
  db.authChannel <- struct{}{}

  return nil
}

func(db *Supabase) Signout() error {
  if err := db.client.Auth.Logout(); err != nil {
    utils.Error("Failed to log user out: %v", err)
    return errors.New("Logout Failure")
  }
  if err := db.session.StoreLoggedoutSession(); err != nil {
    utils.Error("Failed to store Loggedout Session: %v", err)
    return errors.New("Store Logged Out Session Error")
  }

  return nil
}
