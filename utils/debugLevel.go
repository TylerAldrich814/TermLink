package utils

import "github.com/gdamore/tcell/v2"


type DebugLevel uint


const (
  DebugLog DebugLevel = iota+1
  DebugWarn
  DebugError
)

func(l DebugLevel)String() string {
  return [...]string{
    "DebugLog",
    "DebugWarn",
    "DebugError",
  }[l-1]
}

func(l DebugLevel) Tag() string {
  switch l.String(){
  case "DebugLog":
    return "[green]"
  case "DebugWarn":
    return "[yellow]"
  case "DebugError":
    return "[red]"
  default:
    return "[darkred]"
  }
}

func(l DebugLevel) Idx() uint {
  return uint(l)
}

func(l DebugLevel) Color() tcell.Color {
  switch l.String() {
  case "DebugLog":
    return tcell.ColorGreen
  case "DebugWarn":
    return tcell.ColorYellow
  case "DebugError":
    return tcell.ColorRed
  default:
    return tcell.ColorDarkRed
  }
}
