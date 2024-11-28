package db

import (
	"errors"
	"time"

  // "github.com/supabase-community/gotrue-go/types"
  "github.com/supabase-community/supabase-go"
)

const (
  fakeLogin = false
  fakeUsers = false
  fakseMsgs = false
)

type Supabase struct {
  client       *supabase.Client
  url          string
  anonKey      string
  userID       string
  session      *UserSession
  killRefresh  chan struct{}
  authChannel  chan struct{}
}

func InitSupbase(
  url     string,
  anonKey string,
)( *Supabase, error ){
  if len(url) == 0 {
    return nil, errors.New("Supabse URL failed to laod")
  }
  if len(anonKey) == 0 {
    return nil, errors.New("Supabse anonkey failed to laod")
  }
  opts := supabase.ClientOptions {}
  client, err := supabase.NewClient(
    url, 
    anonKey,
    &opts,
  )
  if err != nil {
    return nil, errors.New("Failed to create Supabase Client")
  }

  db := &Supabase{
    client      : client,
    url         : url,
    anonKey     : anonKey,
    userID      : "",
    killRefresh : make(chan struct{}),
    authChannel : make(chan struct{}),
  }

  session, err := LoadUserSession()
  if err == nil {
    if fakeLogin {
      go func(){
        time.Sleep(100*time.Millisecond)
        db.authChannel <- struct{}{}
      }()
    } else {
      session.StartAutoTokenRefresh(
        db,
        5 * time.Minute,
      )
      go func(){
        // TODO: Make Token Refresh Call to Supbase

      }()
    }
  }

  db.session = session

  return db, nil
}

func(client *Supabase) RefreshTokens() error {

  return nil
}
