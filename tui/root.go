package tui

import (
	"fmt"
	"sync"

	"github.com/TylerAldrich814/TermLink/db"
	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TermLinkTUI struct {
  app         *tview.Application
  db          *db.Supabase
  pages       *tview.Pages
  rootWindow  *tview.Flex
  debugWindow *utils.DebugWindow

  header      *Header
  termPages   map[Page]TermLinkPage
  currentPage Page
  kill        chan struct{}
  wg          sync.WaitGroup
}

func(tui *TermLinkTUI) GenerateAuthPage() *TermLinkTUI {
  if tui.app == nil {
    panic("TermLinkTUI hasn't been initialized yet")
  }

  authPage := GetAuthPage(tui, Signup)

  tui.termPages[pAuth] = authPage
  tui.pages.AddPage(pAuth.String(), authPage.GenerateUI(), true, false)


  tui.pages.SwitchToPage(tui.currentPage.String())
  return tui
}

func(tui *TermLinkTUI) HandleInput() *TermLinkTUI {
  tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    tui.termPages[tui.currentPage].HandleInput(event)
    return event
  })
  return tui
}

// A Go Routine that's attached to Supabase.AuthChannel - When a user successfully signs in.
// This Channel will be triggered, which will tell the rest of the app that the User is now
// Authenticated. Where we can then load the rest of the applicaiton.
func(tui *TermLinkTUI) AwaitForAuthentication() *TermLinkTUI{
  go func(){
    utils.Warn("Awaiting for User Authentication...")
    authChannel := tui.db.GetAuthChannel()
    for {
      select {
      case <-authChannel:
        // When authChannel is received from Supabase. We then update the structure of TermLink
        // By creating&adding our HEader, then reattaching pages into a dedicated Flex view that
        // will take up the rest of the page
        utils.Warn("User Successfully Logged in!")
        tui.currentPage = pContact

        header := GetHeader(tui, tui.kill)
        tui.header = header

        contactsPage := GetContactsPage(tui, header, tui.kill)

        tui.termPages[pContact] = contactsPage
        tui.pages.AddPage(
            pContact.String(), 
            contactsPage.GenerateUI(), 
            true, false,
          )

        tui.SwitchToPage(pContact)
        return 
      }
    }
  }()

  return tui
}

func(tui *TermLinkTUI) SwitchToPage(page Page) {
  tui.app.QueueUpdateDraw(func(){
    tui.pages.SwitchToPage(pContact.String())
    tui.SetFocus(tui.termPages[pContact].StartFocus())
  })
}

func(tui *TermLinkTUI) Start() {
  err := tui.app.
    SetRoot(tui.rootWindow, true).
		SetFocus(tui.termPages[tui.currentPage].StartFocus()).
    Run()
  if err != nil {
    panic(err)
  }
}

func(tui *TermLinkTUI) Stop(){
  close(tui.kill)
  tui.wg.Wait()

  tui.app.Stop()
  if tui.debugWindow != nil {
    tui.debugWindow.Stop()
  }
}

func( tui *TermLinkTUI) SetFocus(focus tview.Primitive) {
  tui.app.SetFocus(focus)
}

func( tui *TermLinkTUI) SetRoot(root tview.Primitive) {
  tui.app.SetRoot(root, true)
}

// Used for when we need to use a Temporary Root View, i.e., A Modal.
// This function simply pushes the original Root View back into our 
// tview Application and reset the focus back to where it was before.
func( tui *TermLinkTUI) ResetMainRoot() {
  tui.app.
    SetRoot(tui.rootWindow, true).
    SetFocus(tui.termPages[tui.currentPage].RefreshFocus())
}

// Changes the Root View to the provided Page.
func(tui *TermLinkTUI) ChangePage(page Page) error {
  termPage := tui.termPages[page]
  if termPage == nil {
    return fmt.Errorf("The Page %s doesn't exist yet", page.String())
  }
  tui.pages.SwitchToPage(page.String())
  tui.app.SetFocus(termPage.StartFocus())
  tui.currentPage = page
  return nil
}

func GetTermLinkTUI(
  mode utils.Build,
  db   *db.Supabase,
) *TermLinkTUI {
  app := tview.NewApplication().
    EnableMouse(true).
    EnablePaste(true)
  utils.InitializeDebugWindow(app, mode)

  tl := &TermLinkTUI{ 
    app         : app,
    db          : db,
    pages       : tview.NewPages(),
    termPages   : map[Page]TermLinkPage{},
    currentPage : pAuth,
    rootWindow  : tview.NewFlex(),
    kill        : make(chan struct{}),
    debugWindow : utils.GetInstance(),
  }

  tl.rootWindow.
    AddItem(tl.pages, 0, 2, true)

  if tl.debugWindow != nil {
    debug := tview.NewFlex().
      AddItem(tl.debugWindow.View, 0, 1, false)
    debug.
      SetTitle("Debug Messages").
      SetTitleAlign(tview.AlignCenter).
      SetBorder(true).
      SetBorderPadding(1, 1, 2, 2)

    tl.rootWindow.
      AddItem(debug, 0, 1, true)
  }

  return tl
}

