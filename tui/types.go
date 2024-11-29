package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page uint
const (
  pAuth Page = iota+1
  pNewMessage
  pInbox
  pContacts
  pChatrooms
  pSettings
  pChatroom
  pContact
)

func(p Page) String() string {
  return [...]string{
    "pAuth",
    "pNewMessage",
    "pInbox",
    "pContacts",
    "pChatrooms",
    "pSettings",
    "pChatroom",
    "pContact",
  }[p-1]
}

type TermLinkPage interface {
  GetPageKind()     Page
  GenerateUI()      tview.Primitive
  StartFocus()      tview.Primitive
  RefreshFocus()    tview.Primitive
  ShiftFocus()
  HandleInput(event *tcell.EventKey)
}
