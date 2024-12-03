package db

import "time"

type Contact struct {
  ID        string `json:"id"`
  Username  string `json:"username"`
  Status    Status `json:"status"`
}

type Contacts struct {
  ID        string    `json:"id"`
  UserID    string    `json:"user_id"`
  Contacts  []string  `json:"contacts"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

//TODO: Create Database Access methods for getting updated Contact information
