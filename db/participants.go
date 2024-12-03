package db

import (
	"encoding/json"
	"errors"
	"time"
)

type Role uint
const (
  RoleAdmin Role = iota+1
  RoleMember
)
var roleFromString = map[string]Role {
  "admin"  : RoleAdmin,
  "member" : RoleMember,
}
var roleToString = map[Role]string {
  RoleAdmin  : "admin",
  RoleMember : "member",
}
func(r Role) String()string {
  if str, ok := roleToString[r]; ok {
    return str
  }
  return "unknown"
}
func(r Role) MarshalJSON()( []byte,error ){
  str := r.String()
  if str == "unknown" {
    return nil, errors.New("Invalid Role Value")
  }
  return json.Marshal(str)
}
func(r *Role) UnmarshallJSON(data []byte) error {
  var roleStr string
  if err := json.Unmarshal(data, &roleStr); err != nil {
    return err
  }
  if role, ok := roleFromString[roleStr]; ok {
    *r = role
    return nil
  }
  return errors.New("Invalid Role Value")
}

type Participants struct {
  ID               string    `json:"id"`
  ConversationID   string    `json:"conversation_id"`
  UserID           string    `json:"user_id"`
  Role             Role      `json:"role"`
  JoinedAt         time.Time `json:"joined_at"`
}
