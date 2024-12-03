package db

type User struct {
  ID        string   `json:"id"`
  Username  string   `json:"username"`
  UserID    string   `json:"user_id"`
  Status    Status   `json:"status"`
}
