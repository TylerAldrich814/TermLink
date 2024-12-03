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
  title := "[green]▼ [teal]TermLink[:]"

  newMsgBullet := utils.GetBullet(
    " New Message  ",
    "green",
    "black",
  )
  inboxBullet := utils.GetBullet(
    " Inbox        ",
    "teal",
    "black",
  )
  contactsBullet := utils.GetBullet(
    " Contacts     ",
    "green",
    "black",
  )
  chatroomsBullet := utils.GetBullet(
    " Chatrooms    ",
    "teal",
    "black",
  )
  settingsBullet := utils.GetBullet(
    " Settings     ",
    "green",
    "black",
  )

  curPage := app.currentPage
  utils.Warn("Current Page: %s", curPage.String())
  switch curPage.String() {
  case pNewMessage.String():
    newMsgBullet.UpdateColors(
      "red",
      "black",
    )
    newMsgBullet.UpdateContent(
      " - Contacts   ",
    )
    newMsgBullet.Build()
  case pInbox.String():
    inboxBullet.UpdateColors(
      "red",
      "black",
    )
    inboxBullet.UpdateContent(
      " - Inbox     ",
    )
    inboxBullet.Build()
  case pContact.String():
    contactsBullet.UpdateColors(
      "red",
      "black",
    )
    contactsBullet.UpdateContent(
      " - Contacts  ",
    )
    contactsBullet.Build()
  case pChatrooms.String():
    chatroomsBullet.UpdateColors(
      "red",
      "black",
    )
    chatroomsBullet.UpdateContent(
      " - Chatrooms ",
    )
    chatroomsBullet.Build()
  case pSettings.String():
    settingsBullet.UpdateColors(
      "red",
      "black",
    )
    settingsBullet.UpdateContent(
      " - Settings  ",
    )
    settingsBullet.Build()
  }

  dropdown := &HeaderDropdown{
    app        : app,
    newMessage : newMsgBullet,
    inbox      : inboxBullet,
    contacts   : contactsBullet,
    chatrooms  : chatroomsBullet,
    settings   : settingsBullet,
  }

  dropdown.dropdown = tview.NewDropDown().
    SetLabel(title).
    SetLabelColor(tcell.ColorTeal).
    SetOptions(
      []string{
        newMsgBullet.Item(),
        inboxBullet.Item(),
        contactsBullet.Item(),
        chatroomsBullet.Item(),
        settingsBullet.Item(),
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
    Foreground(tcell.ColorDarkGray)
  btnSpacer := tview.NewTextView().
    SetDynamicColors(true).
    SetTextAlign(tview.AlignCenter).
    SetText("[teal]┃[:]")
    // SetText("[teal]｜[:]")

  helpText := utils.BuildBullet(
    utils.LeftBullet,
    "  Help  ",
    "green",
    "black",
  )
  closeText := utils.GetBlock(
    " Close ",
    "green",
    "black",
  )
  logoutText := utils.BuildBullet(
    utils.RightBullet,
    " Logout ",
    "green",
    "black",
  )

  helpBtn := tview.NewButton(helpText.Item()).
    SetSelectedFunc(func(){
      utils.Log("[yellow] -> Help")
    }).
    SetStyle(btnStyle)

  closeBtn := tview.NewButton(closeText.Item()).
    SetSelectedFunc(func(){
      header.app.Stop()
    }).
    SetStyle(btnStyle)

  logoutBtn := tview.NewButton(logoutText.Item()).
    SetSelectedFunc(func(){
      if err := header.app.db.Signout(); err != nil {
        header.app.SetRoot(MessageModal(
          "Logout Failure",
          "Something unexpected happened during user sign out..\nPlease try again.",
          func(){
            header.app.ResetMainRoot()
          },
        ))
      } else {
        header.app.ChangePage(pAuth)
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
        AddItem(btnSpacer, 1, 0, false).
        AddItem(closeBtn, closeText.Length(), 0, false).
        AddItem(btnSpacer, 1, 0, false).
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
