package tui

import (
	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HeaderDropdown struct {
  app        *TermLinkTUI
  dropdown   *tview.DropDown
  Title      string

  newMessage *utils.BulletItem
  inbox      *utils.BulletItem
  contacts   *utils.BulletItem
  chatrooms  *utils.BulletItem
  settings   *utils.BulletItem
}

func(h *HeaderDropdown) handleDropdownSelection(
  text  string,
  index int,
) {
  switch text {
  case h.newMessage.Text(): 
    h.app.ChangePage(pNewMessage)
  case h.inbox.Text():
    h.app.ChangePage(pInbox)
  case h.chatrooms.Text():
    h.app.ChangePage(pChatrooms)
  case h.settings.Text():
    h.app.ChangePage(pSettings)
  }
}

func GetHeaderDropdown(
  app  *TermLinkTUI,
) *HeaderDropdown {
  title := "▼ TermLink"

  newMsgText := utils.Bullet(
    " New Message  ",
    "green",
    "black",
  )
  inboxText := utils.Bullet(
    " Inbox        ",
    "green",
    "black",
  )
  contactsText := utils.Bullet(
    " Contacts     ",
    "green",
    "black",
  )
  chatroomsText := utils.Bullet(
    " Chatrooms    ",
    "green",
    "black",
  )
  settingsText := utils.Bullet(
    " Settings     ",
    "green",
    "black",
  )

  dropdown := &HeaderDropdown{
    app        : app,
    newMessage : newMsgText,
    inbox      : inboxText,
    contacts   : contactsText,
    chatrooms  : chatroomsText,
    settings   : settingsText,
  }

  dropdown.dropdown = tview.NewDropDown().
    SetLabel(title).
    SetLabelColor(tcell.ColorTeal).
    SetOptions(
      []string{
        newMsgText.View(),
        inboxText.View(),
        contactsText.View(),
        chatroomsText.View(),
        settingsText.View(),
      },
      dropdown.handleDropdownSelection,
    )

  return dropdown
}

type Header struct {
  app         *TermLinkTUI
  header      *tview.Flex
  bannerTitle string
  
  dropDown  *HeaderDropdown
  helpBtn   *tview.Button
  closeBtn  *tview.Button
  logoutBtn *tview.Button

  kill chan struct{}
}

func(h *Header) UpdateBanner(){
  var bannerText string

  bannerText = "[green]Welcome to Termlink!"

  h.bannerTitle = bannerText
}

func(h *Header) Header() *tview.Flex {
  return h.header
}

func GetHeader(
  app   *TermLinkTUI,
  kill  chan struct{},
) *Header {
  header := &Header{
    app  : app,
    kill : kill,
  }

  btnStyle := tcell.StyleDefault.
    Background(tcell.ColorDarkGreen).
    Foreground(tcell.ColorDarkGray).
    Bold(true)
  btnSpacer := tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignCenter).
    SetText(" | ")


  helpText := utils.Bullet(
    "    Help    ",
    "green",
    "black",
  )
  closeText := utils.Bullet(
    " Close ",
    "green",
    "black",
  )
  logoutText := utils.Bullet(
    "Logout",
    "green",
    "black",
  )

  helpBtn := tview.NewButton(helpText.View()).
    SetSelectedFunc(func(){
      utils.Log("[yellow] -> Help")
    }).
    SetStyle(btnStyle)

  closeBtn := tview.NewButton(closeText.View()).
    SetSelectedFunc(func(){
      header.app.Stop()
    }).
    SetStyle(btnStyle)

  logoutBtn := tview.NewButton(logoutText.View()).
    SetSelectedFunc(func(){
      if err := header.app.db.Signout(); err != nil {
        header.app.SetRoot(MessageModal(
          "Logout Failure",
          "Something unexpected happened during user sign out..\nPlease try again.",
          func(){
            header.app.ResetMainRoot()
          },
        ))
      }
    }).
    SetStyle(btnStyle)

  header.UpdateBanner()
  banner := tview.NewTextView().
    SetDynamicColors(true).
    SetText(header.bannerTitle).
    SetTextAlign(tview.AlignLeft)

  dropdown := GetHeaderDropdown(app)

  headerView := tview.NewFlex().
    AddItem(dropdown.dropdown, 10, 1, false).
    AddItem(tview.NewBox(), 2, 0, false).
    AddItem(banner, len(header.bannerTitle), 1, false).
    AddItem(
      tview.NewFlex().
        AddItem(tview.NewBox(), 0, 1, false).
        AddItem(helpBtn, helpText.Length(), 0 ,false).
        AddItem(btnSpacer, 3, 0, false).
        AddItem(closeBtn, closeText.Length(), 0, false).
        AddItem(btnSpacer, 3, 0, false).
        AddItem(logoutBtn, logoutText.Length(), 0, false).
        AddItem(tview.NewBox(), 2, 0, false),
      0, 1, false,
    )
  headerView.SetBorder(true)
  headerView.SetBorderPadding(1, 1, 1, 1)

  header.dropDown  = dropdown
  header.helpBtn   = helpBtn
  header.closeBtn  = closeBtn
  header.logoutBtn = logoutBtn
  header.header = headerView

  return header
}
