package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/TylerAldrich814/TermLink/utils"
	sb "github.com/supabase-community/supabase-go"
)

const (
  fakeLogin = false
  fakeUsers = false
  fakseMsgs = false
)

type Database  struct {
  sqlite       *sql.DB
  supabase     *sb.Client
  url          string
  anonKey      string
  userID       string
  session      *UserSession
  killRefresh  chan struct{}
  authChannel  chan struct{}
}

func InitDatabase(
  url     string,
  anonKey string,
)( *Database, error ){
  if len(url) == 0 {
    return nil, errors.New("Supabse URL failed to laod")
  }
  if len(anonKey) == 0 {
    return nil, errors.New("Supabse anonkey failed to laod")
  }
  opts := sb.ClientOptions {}
  client, err := sb.NewClient(
    url, 
    anonKey,
    &opts,
  )
  if err != nil {
    return nil, errors.New("Failed to create Database Client")
  }

  sql, err := InitSQLite()
  if err != nil {
    return nil, err
  }

  db := &Database{
    sqlite      : sql,
    supabase    : client,
    url         : url,
    anonKey     : anonKey,
    userID      : "",
    killRefresh : make(chan struct{}),
    authChannel : make(chan struct{}),
  }

  return db, nil
}

func(db *Database) TrySession() error {
  session, err := LoadUserSession()
  if err == nil {
    utils.Warn("Session Found: %s", session.RefreshToken)
    db.session = session

    err = db.TokenSignin()
    if err == nil {
      session.StartAutoTokenRefresh(
        db,
        5 * time.Minute,
      )
      go func(){
        db.authChannel <-struct{}{}
      }()
    } else {
      utils.Error("AUTH ERROR : %w", err)
    }
  }

  db.session = session

  return nil
}

func(supabase *Database) RefreshTokens() error {

  return nil
}

func(db *Database) GetAuthChannel() chan struct{} {
  return db.authChannel
}

func(db *Database) CloseAuthChannel()  {
  if db.authChannel != nil {
    utils.Warn("Supbase.AuthChannel is now closed")
    close(db.authChannel)
  }
}
