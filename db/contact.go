package db

type Contact struct {
  ID        string `json:"id"`
  Username  string `json:"username"`
  Status    Status `json:"status"`
}

//TODO: Create Database Access methods for getting updated Contact information
