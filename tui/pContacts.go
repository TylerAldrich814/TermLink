package tui

import (
	"fmt"
	"strings"

	"github.com/TylerAldrich814/TermLink/components"
	"github.com/TylerAldrich814/TermLink/db"
	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
  maxTagLength       = 30
  defaultBorderColor = int32(0x857048)
  focusedBorderColor = int32(0x4d2929)
)

type ContactItem struct {
  app        *TermLinkTUI
  contact    *db.Contact
  view       *tview.Flex

  tagView    *tview.TextView
  spacer     *tview.Box
  textStatus *tview.TextView

  username   string
  bullet     *utils.BulletItem
  contactCB  func(contact *db.Contact)
  kill       chan struct{}
}

func(c *ContactItem) CreateTag(focused bool) {
  statusBullet, offset := components.StatusBullet(&c.contact.Status)
  username := fmt.Sprintf(
    "%s  %s",
    statusBullet,
    c.contact.Username,
  )

  ulen := len(username)-offset

  if ulen < maxTagLength {
    pad := strings.Repeat(" ", maxTagLength -ulen)
    username = username + pad
  }

  c.username = username
  c.bullet = utils.BuildBullet(
    utils.LeftBullet,
    username,
    fmt.Sprintf(
      "#%x",
      func()int32{
        if focused { return focusedBorderColor}
        return defaultBorderColor
      }(),
    ),
    "black",
  )
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

  c.CreateTag(false)

  c.tagView = tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignLeft).
    SetText(c.bullet.Item())

  c.spacer = tview.NewBox().
    SetBackgroundColor(
      tcell.NewHexColor(defaultBorderColor),
    )

  c.textStatus = tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignCenter).
    SetText(fmt.Sprintf(
      "[teal][[:]%s[:][teal]] ",
      components.StatusView(&contact.Status),
    ))
  c.textStatus.SetBackgroundColor(tcell.NewHexColor(defaultBorderColor))

  c.view = tview.NewFlex().
    AddItem(c.tagView, maxTagLength, 0, false).
    AddItem(c.spacer, 0, 1, false).
    AddItem(c.textStatus, len(contact.Status.String())+3, 1, false)

  c.view.SetBorder(true)
  c.view.SetBorderColor(tcell.NewHexColor(defaultBorderColor))


  return c
}

func(c *ContactItem) selectedCallback(focused bool) {
  if focused {
    c.view.SetBorderColor(
      tcell.NewHexColor(focusedBorderColor),
    )

    c.contactCB(c.contact)
    c.CreateTag(true)
    c.tagView.SetText(c.bullet.Item())
    c.spacer.SetBackgroundColor(tcell.NewHexColor(focusedBorderColor))
    c.textStatus.SetBackgroundColor(tcell.NewHexColor(focusedBorderColor))
  } else {
    c.view.SetBorderColor(tcell.NewHexColor(defaultBorderColor))
    c.CreateTag(false)
    c.tagView.SetText(c.bullet.Item())
    c.spacer.SetBackgroundColor(tcell.NewHexColor(defaultBorderColor))
    c.textStatus.SetBackgroundColor(tcell.NewHexColor(defaultBorderColor))
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

// Updates the Contact Body for when a Contact List item is selected in the contact List panel
func(c *ContactsPage) updateBody(contact *db.Contact) {
  c.visibleContact = contact
  if contact == nil {
    c.contactBody.Clear().
      SetTitle("Contacts")
  } else {
    topPad := 2
    botPad := 3

    info := tview.NewTextView().
      SetTextAlign(tview.AlignCenter).
      SetDynamicColors(true).
      SetText(contact.Username)

    info.SetBorder(true)
    infoContainer := tview.NewFlex().
      SetDirection(tview.FlexRow).
      AddItem(tview.NewBox(), topPad, 0, false).
      AddItem(
        info,
        0, 1, false,
      ).
      AddItem(tview.NewBox(), botPad, 0, false)

    msgs := tview.NewFlex()

    msgs.SetBorder(true)
    msgContainer := tview.NewFlex().
      SetDirection(tview.FlexRow).
      AddItem(tview.NewBox(), topPad, 0, false).
      AddItem(
        msgs,
        0, 1, false,
      ).
      AddItem(tview.NewBox(), botPad, 0, false)

    c.contactBody.Clear().
      AddItem(
        tview.NewFlex().
          AddItem(infoContainer, 0, 2, false).
          AddItem(msgContainer, 0, 3, false),
        0, 1, false,
      )
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
