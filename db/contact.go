package db

type Contact struct {
  ID        string `json:"id"`
  Username  string `json:"username"`
  Status    Status `json:"status"`
}

type Contacts struct {
  ID        string `json:"id"`
  UserID    string `json:"user_id"`
  Contacts  []User `json:"contacts"`
}

//TODO: Create Database Access methods for getting updated Contact information
