package db

import (
	"encoding/json"
	"errors"
)

type Status uint

const (
  StatusActive Status = iota+1
  StatusAway
  StatusInactive
)

var statusFromString = map[string]Status{
  "active"   : StatusActive,
  "away"     : StatusAway,
  "inactive" : StatusInactive,
}
var statusToString = map[Status]string {
  StatusActive   : "active",
  StatusAway     : "away",
  StatusInactive : "inactive",
}

func(s Status)String()string {
  if str, ok := statusToString[s]; ok { 
    return str 
  }
  return "unknown"
}
func(s Status) MarshalJSON()( []byte,error ){
  str := s.String()
  if str == "unknown" {
    return nil, errors.New("Invalid Status Value")
  }
  return json.Marshal(str)
}

func(s *Status) UnmarshalJSON(data []byte) error {
  var statusStr string
  if err := json.Unmarshal(data, &statusStr); err != nil {
    return err
  }
  if status, ok := statusFromString[statusStr]; ok {
    *s = status
    return nil
  }
  return errors.New("Invalid Status Value")
}
