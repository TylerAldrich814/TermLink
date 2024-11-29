package tui

import (
	"fmt"

	"github.com/TylerAldrich814/TermLink/components"
	"github.com/TylerAldrich814/TermLink/db"
	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ContactItem struct {
  app        *TermLinkTUI
  contact    *db.Contact
  view       *tview.Flex
  dotStatus  *tview.TextView
  username   *tview.TextView
  spacer     *tview.Box
  textStatus *tview.TextView

  contactCB  func(contact *db.Contact)
  kill       chan struct{}
}

func NewContactItem(
  app       *TermLinkTUI,
  contact   *db.Contact,
  contactCB func(contact *db.Contact),
  kill      chan struct{},
) *ContactItem {
  c := &ContactItem{
    app     : app,
    contact : contact,
    kill    : kill,
    contactCB: contactCB,
  }
  c.dotStatus = tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignCenter).
    SetText(" " + components.StatusBullet(&contact.Status))
  c.username = tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignCenter).
    SetText(fmt.Sprintf("[green]%s[:]", contact.Username))
  c.spacer = tview.NewBox()
  c.textStatus = tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignCenter).
    SetText(fmt.Sprintf(
      "[teal][[:]%s[:][teal]]",
      components.StatusView(&contact.Status),
    ))
  c.view = tview.NewFlex().
    AddItem(c.dotStatus, 3, 0, false).
    AddItem(c.username, len(contact.Username)+4, 0, false).
    AddItem(c.spacer, 0, 1, false).
    AddItem(c.textStatus, len(contact.Status.String())+2, 1, false)

  c.view.SetBorder(true)

  return c
}

func(c *ContactItem) selectedCallback(focused bool) {
  if focused {
    c.username.
      SetText(fmt.Sprintf("[orange]%s[:]", c.contact.Username))
    focusedColor := tcell.ColorDarkGreen

    c.dotStatus.SetBackgroundColor(focusedColor)
    c.username.SetBackgroundColor(focusedColor)
    c.textStatus.SetBackgroundColor(focusedColor)
    c.spacer.SetBackgroundColor(focusedColor)
    c.view.SetBackgroundColor(focusedColor)

    c.contactCB(c.contact)
  } else {
    c.username.
      SetText(fmt.Sprintf("[green]%s[:]", c.contact.Username))

    c.dotStatus.SetBackgroundColor(tcell.ColorNone)
    c.username.SetBackgroundColor(tcell.ColorNone)
    c.textStatus.SetBackgroundColor(tcell.ColorNone)
    c.spacer.SetBackgroundColor(tcell.ColorNone)
    c.view.SetBackgroundColor(tcell.ColorNone)
  }
}

type ContactsPage struct {
  app            *TermLinkTUI
  view           *tview.Flex
  header         *Header
  contactList    *components.ScrollView
  contacts       []*ContactItem

  contactBody    *tview.Flex

  visibleContact *db.Contact

  kill  chan struct{}
}
func(c *ContactsPage) GetPageKind() Page {
  return pContact
}

func(c *ContactsPage) contactsPanel() *tview.Flex {
  c.contactList = components.NewScrollView().
    SetItemSize(3)

  for _, contact := range c.contacts {
    utils.Warn("Creating Contact: %s", contact.contact.Username)
    c.contactList.AddItemWithSelectedFunc(
      contact.view,
      contact.selectedCallback,
    )
  }

  c.view.SetDirection(tview.FlexRow).
    AddItem(c.header.header, 5, 0, false)


  contacts := tview.NewFlex().
    AddItem(tview.NewBox(), 2, 0, false).
    AddItem(c.contactList, 0, 1, false).
    AddItem(tview.NewBox(), 2, 0, false)

  contacts.SetBorder(true)

  return contacts
}

func(c *ContactsPage) updateBody(contact *db.Contact) {
  c.visibleContact = contact
  if contact == nil {
    c.contactBody.Clear().
      SetTitle("Contacts")
  } else {
    c.contactBody.Clear().
      SetTitle(c.visibleContact.Username)
  }
}

func(c *ContactsPage) GenerateUI() tview.Primitive {
  contacts := c.contactsPanel()

  body := tview.NewFlex().
    AddItem(
      tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(contacts, 0, 4, false).
        AddItem(tview.NewBox(), 0, 1, false),
      0, 1, false,
    ).
    AddItem(
      c.contactBody,
      0, 2, false,
    )

  body.SetBorder(true)
  body.SetTitle(" Contacts ")

  wrapped := tview.NewFlex().
    AddItem(tview.NewBox(), 4, 0, false).
    AddItem(
      tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(tview.NewBox(), 2, 0, false).
        AddItem(tview.NewBox(), 1, 0, false).
        AddItem(body, 0, 1, false).
        AddItem(tview.NewBox(), 4, 0, false),
      0, 1, false,
    ).
    AddItem(tview.NewBox(), 2, 0, false)
  
  c.view.AddItem(wrapped, 0, 1, false)

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
    app         : app,
    header      : header,
    view        : tview.NewFlex(),
    contactBody : tview.NewFlex(),
    kill        : kill,
  }
  contacts.contactBody.
    SetDirection(tview.FlexRow).
    SetTitle("Contacts")
  contacts.contactBody.SetBorder(true)

  // TODO: TEMP -- We need to load contacts in dynmaically. 
  //    - Create a ContactsListUpdater Channel
  //    - Initially load contacts from DB / Local storage
  //    - Attach Updater Channel to a database function


  contacts.contacts = []*ContactItem {
    NewContactItem(
      app, 
      &db.Contact{
        ID   : "0",
        Username : "SomeUser",
        Status   : db.StatusActive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app, 
      &db.Contact{
        ID   : "0",
        Username : "SomeUser",
        Status   : db.StatusActive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app, 
      &db.Contact{
        ID   : "0",
        Username : "SomeUser",
        Status   : db.StatusActive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "SomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "NeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "LastSomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "2",
        Username : "LastNeverOnline",
        Status   : db.StatusInactive,
      },
      contacts.updateBody,
      kill,
    ),
    NewContactItem(
      app,
      &db.Contact{
        ID   : "1",
        Username : "LastSomeOtherUser",
        Status   : db.StatusAway,
      },
      contacts.updateBody,
      kill,
    ),
  }

  return contacts
}
