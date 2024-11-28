package tui

import (
	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TermLinkPage interface {
  GetPageKind()     utils.Page
  GenerateUI()      tview.Primitive
  StartFocus()      tview.Primitive
  RefreshFocus()    tview.Primitive
  ShiftFocus()
  HandleInput(event *tcell.EventKey)
}
