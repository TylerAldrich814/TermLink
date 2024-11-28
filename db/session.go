package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/supabase-community/gotrue-go/types"
)

var filePath = "session.json"

type UserSession struct {
  AccessToken  string    `json:"access_token"`
  RefreshToken string    `json:"refresh_token"`
  Email        string    `json:"email"`
  ExpiresAt    time.Time `json:"expires_at"`
}

// Attempts to Load and return the User's Session data.
func LoadUserSession()( *UserSession, error) {
  file, err := os.Open(filePath)
  if err != nil {
    return &UserSession{}, nil
  }
  defer file.Chdir()

  var sess UserSession
  decoder := json.NewDecoder(file)
  if err := decoder.Decode(&sess); err != nil {
    return &UserSession{}, err
  }
  return &sess, nil
}

func NewSession(
  email         string,
  tokenResponse *types.TokenResponse,
)( *UserSession, error ){
  session := &UserSession{
    Email        : email,
    AccessToken  : tokenResponse.AccessToken,
    RefreshToken : tokenResponse.RefreshToken,
		ExpiresAt    : time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
  }
  err := session.StoreUserSession()
  if err != nil {
    return nil, err
  }
  return session, nil
}

func(session *UserSession) UpdateCurrentSession(
  tokenResponse *types.TokenResponse,
) error {
  session.AccessToken = tokenResponse.AccessToken
  session.RefreshToken = tokenResponse.RefreshToken
  session.ExpiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresAt) * time.Second)

  if err := session.StoreUserSession(); err != nil {
    return err
  }
  return nil
}

// Attempts to extract a UserID from the JWT token provided by Supabase.
func(session *UserSession) ExtractUserID()( string,error ){
  token, _, err := jwt.NewParser().ParseUnverified(session.AccessToken, jwt.MapClaims{})
  if err != nil {
    return "", fmt.Errorf("Failed to parse token: %w", err)
  }

  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    return "", fmt.Errorf("Invalid Token Claims")
  }

  userId, ok := claims["sub"].(string)
  if !ok {
    return "", fmt.Errorf("user_id (sub) not found in token")
  }
  return userId, nil
}

// After the User is Authenticated. We take the UserSession data and store it
// on the users Machines
func(session *UserSession) StoreUserSession() error {
  file, err := os.Create(filePath)
  if err != nil {
    return fmt.Errorf("Failed to store UserSession: %v", err)
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  if err := encoder.Encode(session); err != nil {
    return fmt.Errorf("Failed to Write UserSession: %v", err)
  }
  return nil
}

// When the user logs out. We simply replace the UserSession file with their Email.
// This is used to autofill the Email field when the user loads the app back up.
func(session *UserSession) StoreLoggedoutSession() error {
  file, err := os.Create(filePath)
  if err != nil {
    return fmt.Errorf("Failed to store Logged Out User Session: %v", err)
  }
  defer file.Close()

  loggedOutSession := struct {
    Email string `json:"email"`
  }{
    Email: session.Email,
  }

  encoder := json.NewEncoder(file)
  if err := encoder.Encode(loggedOutSession); err != nil {
    return fmt.Errorf("Failed to Write Logged out UserSession: %v", err)
  }
  return nil
}

// Function to test if the users SessionToken is expired or will be expired soon.
// If so, attempts to refresh the users RefreshTokens
func(session *UserSession) EnsureValidSession(db *Supabase) error {
  if !time.Now().After(session.ExpiresAt.Add(-time.Minute * 5)){
    return nil
  }
  fmt.Println("Access Token expired or is about to expire. Refreshing...")
  err := db.RefreshTokens()
  if err != nil {
    return err
  }
  fmt.Println("Tokens refreshed")
  return  nil
}

// Go Routine for automatically refreshing user token in the background.
func(session *UserSession) StartAutoTokenRefresh(
  db              *Supabase,
  refreshInterval time.Duration,
){
  ticker := time.NewTicker(refreshInterval)
  go func(){
    for{
      select{
      case <-ticker.C:
        err := session.EnsureValidSession(db)
        if err != nil {
          log.Printf("Failed to Refresh Token: %v", err)
        }
      case <-db.killRefresh:
        ticker.Stop()
        return
      }
    }
  }()
}
