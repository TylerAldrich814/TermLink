package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type ConversationType uint

const (
  ConversationDirect ConversationType = iota+1
  ConversationGroup
)
var typeFromString = map[string]ConversationType{
  "direct" : ConversationDirect,
  "group"  : ConversationGroup,
}
var typeToString   = map[ConversationType]string {
  ConversationDirect : "direct",
  ConversationGroup  : "group",
}

func(c ConversationType)String() string {
  if str, ok := typeToString[c]; ok {
    return str
  }
  return "unknown"
}
func(c ConversationType) MarshalJSON()( []byte,error ){
  str := c.String()
  if str == "unknown" {
    return nil, errors.New("Invalid ConversationType Value")
  }
  return json.Marshal(str)
}
func(c *ConversationType) UnmarshalJSON(data []byte) error {
  var convoTypeStr string
  if err := json.Unmarshal(data, &convoTypeStr); err != nil {
    return err
  }
  if status, ok := typeFromString[convoTypeStr]; ok {
    *c = status
    return nil
  }
  return errors.New("Invalid ConversationType Value")
}

type Conversations struct {
  ID             string           `json:"id"`
  CreatedAt      time.Time        `json:"created_at"`
  Type           ConversationType `json:"type"`
  UserID         string           `json:"user_id"`
  LastSentUserID string           `json:"last_sent_user_id"`
  OtherUserID    string           `json:"other_user_id"`
}

// For Creating a new Direct; user to user Conversation. 
// On the backend. This RPC call will push the data into
// our 'create_direct_conversation' database function.
// Which will
//  - Create a new conversation row item, with the type set to 'direct'
//  - Create a new participant row item for the conversation creator.
//  - Create a new participant row item for the invitee.
// and finally, returns the conversation_id if successful.
func(db *Database) CreateNewDirectConversation(
  invitee_id string,
)( string,error ){
  if invitee_id == "" {
    return "", fmt.Errorf("InviteeID Cannot be empty")
  }

  res := db.supabase.Rpc(
    "create_direct_conversation",
    "",
    struct{
      CreatorID  string `json:"creator_id"`
      InviteeID  string `json:"invitee_id"`
    }{
      CreatorID: db.userID,
      InviteeID: invitee_id,
    },
  )
  if res == "" {
    return "", fmt.Errorf("Unexpected empty response from RPC")
  }

  return res, nil
}
