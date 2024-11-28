package tui

import (
	"fmt"

	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AuthKind uint
const (
  Signup AuthKind = iota+1
  Login
)
func(auth AuthKind)String()string{
  return [...]string{
    "Signup",
    "Login",
  }[auth-1]
}

const (
  signupToLogin = "Already have an account?"
  loginToSignup = "Don't have an account?"
)

type AuthPage struct {
  app           *TermLinkTUI
  kind          AuthKind
  form          *tview.Flex
  emailForm     *tview.InputField
  passwordForm  *tview.InputField
  submitText    *string
  submitButton  *tview.Button
  quitButton    *tview.Button
  switchButton  *tview.Button
  focusId       int
}

func( auth *AuthPage) authCallback(){
  email := auth.emailForm.GetText()
  passw := auth.passwordForm.GetText()

  validateEmail := utils.ValidateEmail(email)
  validatePassw := utils.ValidatePassword(passw)

  if !validateEmail {
    auth.app.SetRoot(MessageModal(
      auth.kind.String() + " Error",
      "The Email provided is not formatted correctly",
      func(){
        auth.app.ResetMainRoot()
      },
    ))
    return
  }
  if !validatePassw {
    auth.app.SetRoot(MessageModal(
      auth.kind.String() + " Error",
      "The Password Provided isn't strong enough",
      func(){
        auth.app.ResetMainRoot()
      },
    ))
    return
  }

  if auth.kind == Signup {
    if err := auth.app.db.Signup(email, passw); err != nil {
      Error("Failed to Sign up: %v", err)
      auth.app.SetRoot(MessageModal(
        auth.kind.String() + " Error",
        fmt.Sprintf("An Error Occurred: %s", err.Error()),
        func(){
          auth.app.ResetMainRoot()
        },
      ))
      return
    }
    Log("Successfully Signed up!")
  } else { // Login
    Log("UserLogin: %s - %s", email, passw)
    if err := auth.app.db.Login(email, passw); err != nil {
      Error("Failed to Log into account: %v", err)
      auth.app.SetRoot(MessageModal(
        auth.kind.String() + " Error",
        fmt.Sprintf("An Error Occurred: %s", err.Error()),
        func(){
          auth.app.ResetMainRoot()
        },
      ))
      return
    }
    Log("Successfully Logged back in!")
  }
}

func( auth *AuthPage) switchKind(){
  if auth.kind == Signup {
    auth.kind = Login
    auth.form.SetTitle(Login.String())
    auth.submitButton.SetLabel(Login.String())
    auth.switchButton.SetLabel(signupToLogin)
    auth.switchButton.SetLabel(loginToSignup)
  } else {
    auth.kind = Signup
    auth.form.SetTitle(Signup.String())
    auth.submitButton.SetLabel(Signup.String())
    auth.switchButton.SetLabel(signupToLogin)
  }
}

func(auth *AuthPage) GetPageKind() utils.Page {
  return utils.Auth
}

func(auth *AuthPage) GenerateUI() tview.Primitive{
  form := tview.NewForm().
    AddFormItem(auth.emailForm).
    AddFormItem(auth.passwordForm)

  page := tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(form, 6, 1, false).
    AddItem(
      tview.NewFlex().
        AddItem(tview.NewBox(), 14, 0, false).
        AddItem(auth.submitButton, 0, 1, false).
        AddItem(tview.NewBox(), 1, 0, false).
        AddItem(auth.quitButton, 0, 1, false).
        AddItem(tview.NewBox(), 8, 0, false),
      1, 0, false,
    ).
    AddItem(tview.NewBox(), 1, 0, false).
    AddItem(
      tview.NewFlex().
        AddItem(tview.NewBox(), 10, 0, false).
        AddItem(auth.switchButton, 0, 1, false).
        AddItem(tview.NewBox(), 4, 0, false),
      1, 0, false,
    )
  page.SetBorder(true).
    SetTitle(auth.kind.String()).
    SetTitleAlign(tview.AlignCenter).
    SetBorderPadding(0, 2, 2, 2)

  auth.form = page

  return utils.CenterUIComponent(
    page,
    50,
  ) 
}

func(auth *AuthPage) StartFocus() tview.Primitive {
  return auth.emailForm
}

func(auth *AuthPage) RefreshFocus() tview.Primitive {
  switch auth.focusId {
  case 0:
    return auth.emailForm
  case 1:
    return auth.passwordForm
  case 2:
    return auth.submitButton
  case 3:
    return auth.quitButton
  case 4:
    return auth.switchButton
  default:
    auth.focusId = 0
    return auth.emailForm
  }
}

func(auth *AuthPage) ShiftFocus() {
  switch auth.focusId {
  case 0:
    auth.app.SetFocus(auth.emailForm)
  case 1:
    auth.app.SetFocus(auth.passwordForm)
  case 2:
    auth.app.SetFocus(auth.submitButton)
  case 3:
    auth.app.SetFocus(auth.quitButton)
  case 4:
    auth.app.SetFocus(auth.switchButton)
  }
}

func(auth *AuthPage) HandleInput(event *tcell.EventKey) {
  switch event.Key() {
  case tcell.KeyBacktab:
    if auth.focusId == 0 {
      auth.focusId = 4
    } else {
      auth.focusId--
    }
    auth.ShiftFocus()
  case tcell.KeyTab:
    if auth.focusId == 4 {
      auth.focusId = 0
    } else {
      auth.focusId++
    }
    auth.ShiftFocus()
  }
}

func GetAuthPage(
  tui   *TermLinkTUI,
  kind  AuthKind,
) *AuthPage {
  auth := &AuthPage{
    app  : tui,
    kind : kind,
  }
  auth.emailForm = tview.NewInputField().
    SetLabel("Email").
    SetFieldWidth(30)
  auth.passwordForm = tview.NewInputField().
    SetLabel("Password").
    SetFieldWidth(30).
    SetMaskCharacter('*')

  auth.submitButton = tview.NewButton(kind.String()).
    SetSelectedFunc(auth.authCallback)

  auth.quitButton = tview.NewButton("Quit").
    SetSelectedFunc(func(){ tui.Stop() })

  if kind == Signup {
    auth.switchButton = tview.NewButton(signupToLogin).
      SetSelectedFunc(auth.switchKind)
  } else {
    auth.switchButton = tview.NewButton(loginToSignup).
      SetSelectedFunc(auth.switchKind)
  }

  return auth
}
