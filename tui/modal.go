package tui

import "github.com/rivo/tview"

func MessageModal(
  title   string,
  message string,
  onClose func(),
) *tview.Modal {
  modal := tview.NewModal().
    SetText(message).
    AddButtons([]string{"OK"}).
    SetDoneFunc(func(buttonIndex int, buttonLabel string) {
      onClose()
    })
  modal.
    SetTitle(title).
    SetTitleAlign(tview.AlignCenter)

  return modal
}
