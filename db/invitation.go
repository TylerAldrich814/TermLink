package db

import (
	"encoding/json"
	"errors"
	"time"
)

type InvitationStatus uint
const (
  InvitePending InvitationStatus = iota+1
  InviteAccepted
  InviteRejected
)
var inviteStatusToString = map[InvitationStatus]string{
  InvitePending  : "pending",
  InviteAccepted : "accepted",
  InviteRejected : "rejected",
}
var inviteStatusFromString = map[string]InvitationStatus{
  "pending"  : InvitePending,
  "accepted" : InviteAccepted,
  "rejected" : InviteRejected,
}

func(s InvitationStatus) String()string {
  if str, ok := inviteStatusToString[s]; ok {
    return str
  }
  return "unknown"
}
func(s InvitationStatus) MarshalJSON()( []byte, error ){
  str := s.String()
  if str == "unknown"{
    return nil, errors.New("Invalid InvitationStatus Value")
  }
  return json.Marshal(str)
}
func(s *InvitationStatus) UnmarshalJSON(data []byte) error {
  var inviteStatusStr string
  if err := json.Unmarshal(data, &inviteStatusStr); err != nil {
    return err
  }
  if status, ok := inviteStatusFromString[inviteStatusStr]; ok {
    *s = status
    return nil
  }
  return errors.New("Invalid InvitationStatus Value")
}

type Invitations struct {
  ID              string           `json:"id"`
  ConversationsID string           `json:"conversations_id"`
  InviterID       string           `json:"inviter_id"`
  InviteeID       string           `json:"invitee_id"`
  Status          InvitationStatus `json:"invitation_status"`
  CreatedAt       time.Time        `json:"created_at"`
  UpdatedAt       time.Time        `json:"updated_at"`
}
