package utils

import "github.com/gdamore/tcell/v2"


type DebugLevel uint


const (
  Log DebugLevel = iota+1
  Warn
  Error
)

func(l DebugLevel)String() string {
  return [...]string{
    "Log",
    "Warn",
    "Error",
  }[l-1]
}

func(l DebugLevel) Tag() string {
  switch l.String(){
  case "Log":
    return "[green]"
  case "Warn":
    return "[yellow]"
  case "Error":
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
  case "Log":
    return tcell.ColorGreen
  case "Warn":
    return tcell.ColorYellow
  case "Error":
    return tcell.ColorRed
  default:
    return tcell.ColorDarkRed
  }
}
