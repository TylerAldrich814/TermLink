package tui

import (
	"github.com/TylerAldrich814/TermLink/components"
	"github.com/TylerAldrich814/TermLink/db"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


type ContactsPage struct {
  app    *TermLinkTUI
  view   *tview.Flex
  header *Header


  list     *components.ScrollView
  contacts []*db.Contact

  kill  chan struct{}
}
func(c *ContactsPage) GetPageKind() Page {
  return pContact
}
func(c *ContactsPage) GenerateUI() tview.Primitive {
  c.view.SetDirection(tview.FlexRow).
    AddItem(c.header.header, 5, 0, false)

  return c.view
}
func(c *ContactsPage) StartFocus() tview.Primitive {
  return nil
}
func(c *ContactsPage) RefreshFocus() tview.Primitive {
  return nil
}

func(c *ContactsPage) ShiftFocus() {

}

func(c *ContactsPage) HandleInput(event *tcell.EventKey) {
}

func GetContactsPage(
  app   *TermLinkTUI,
  header *Header,
  kill  chan struct{},
) *ContactsPage {
  contacts := &ContactsPage{
    app    : app,
    header : header,
    view   : tview.NewFlex(),
    kill   : kill,
  }
  contacts.contacts = []*db.Contact {
    {
      ID   : "0",
      Username : "SomeUser",
      Status   : db.StatusActive,
    },
    {
      ID   : "1",
      Username : "SomeOtherUser",
      Status   : db.StatusAway,
    },
    {
      ID   : "2",
      Username : "NeverOnline",
      Status   : db.StatusInactive,
    },
    {
      ID   : "1",
      Username : "SomeOtherUser",
      Status   : db.StatusAway,
    },
    {
      ID   : "2",
      Username : "NeverOnline",
      Status   : db.StatusInactive,
    },
    {
      ID   : "1",
      Username : "SomeOtherUser",
      Status   : db.StatusAway,
    },
    {
      ID   : "2",
      Username : "NeverOnline",
      Status   : db.StatusInactive,
    },
    {
      ID   : "1",
      Username : "SomeOtherUser",
      Status   : db.StatusAway,
    },
  }


  return contacts
}
