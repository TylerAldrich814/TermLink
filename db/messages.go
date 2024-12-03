package db

import "time"

type Messages struct {
  ID              string    `json:"id"`
  ConversationID  string    `json:"conversation_id"`
  SenderID        string    `json:"sender_id"`
  Content         string    `json:"content"`
  CreatedAT       time.Time `json:"created_at"`
  Read            bool      `json:"read"`
}

