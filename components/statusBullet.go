package components

import "github.com/TylerAldrich814/TermLink/db"


func StatusView(status *db.Status) string {
  switch status.String() {
  case "active": 
    return "[lightgreen]Active"
  case "away":
    return "[yellow]Away"
  case "inactive":
    return "[darkorange]Inactive"
  default:
    return "[red]unknown"
  }
}

func StatusBullet(status *db.Status) string {
  symbol := "⬤ "
  switch status.String() {
  case "active": 
    return "[lightgreen]"+symbol
  case "away":
    return "[yellow]"+symbol
  case "inactive":
    return "[darkorange]"+symbol
  default:
    return "[red]"+symbol
  }
}
