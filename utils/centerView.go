package utils

import "github.com/rivo/tview"

// Takes in a tview Primitive and wraps it in tview.Flexs 
// Each with FixedSize set to dynamic. This in turn will
// Center the provided Primitive both on the X and Y coordinate Plane
func CenterUIComponent(
  form  tview.Primitive,
  width int,
) *tview.Flex {
  return tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(tview.NewBox(), 0, 2, false).
    AddItem(
      tview.NewFlex().
        AddItem(tview.NewBox(), 0, 2, true).
        AddItem(form, width, 0, true).
        AddItem(tview.NewBox(), 0, 2, true),
      0, 1, false,
    ).
    AddItem(tview.NewBox(), 0, 2, false)
}
