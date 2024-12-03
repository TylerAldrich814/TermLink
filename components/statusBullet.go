package components

import (
	"github.com/TylerAldrich814/TermLink/db"
	"github.com/TylerAldrich814/TermLink/utils"
)

func StatusColor(status *db.Status)(string, int){
  switch status.String() {
  case "active": 
    return "lightgreen", len("lightgreen")
  case "away":
    return "yellow",     len("yellow")
  case "inactive":
    return "darkorange", len("darkorange")
  default:
    return "red",        len("red")
  }
}

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

func StatusBullet(status *db.Status)(string, int){
  switch status.String() {
  case "active": 
    return "[lightgreen]"+utils.FullCircle, len("[lightgreen]")
  case "away":
    return "[yellow]"+utils.FullCircle, len("[yellow]")
  case "inactive":
    return "[darkorange]"+utils.FullCircle, len("[darkorange]")
  default:
    return "[red]"+utils.FullCircle, len("[red]")
  }
}
