package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


type ScrollViewItem struct {
  view    tview.Primitive
  id      int
  update  func(focused bool)
  focused bool
}

func(item *ScrollViewItem) SetFocus(focus bool){
  if item.focused == focus { return }
  item.focused = focus
  item.update(focus)
}

type ScrollView struct {
  *tview.Box
  items      []*ScrollViewItem
  offsetY    int
  itemSize   int
  deactivated bool
  selected   int
}

func NewScrollView() *ScrollView {
  return &ScrollView{
    Box         : tview.NewBox(),
    items       : []*ScrollViewItem{},
    offsetY     : 0,
    itemSize    : 0,
    deactivated : false,
    selected    : -1,
  }
}

func(sv *ScrollView) SetItemSize(size int) *ScrollView {
  sv.itemSize = size
  return sv
}

func(sv *ScrollView) GetOffset() int {
  return sv.offsetY
}

func(sv *ScrollView) AddItem(item tview.Primitive) {
  sv.items = append(
    sv.items,
    &ScrollViewItem {
      id      : len(sv.items),
      view    : item,
      update  : func(_ bool){},
      focused : false,
    },
  )
}

func(sv *ScrollView) AddItemWithSelectedFunc(
  item tview.Primitive,
  fn   func(focused bool),
){
  sv.items = append(
    sv.items,
    &ScrollViewItem {
      id      : len(sv.items),
      view    : item,
      update  : fn,
      focused : false,
    },
  )
}

func(sv *ScrollView) Deactivate() {
  sv.deactivated = true
}
func(sv *ScrollView) Activate(){
  sv.deactivated = false
}

func(sv *ScrollView) Draw(screen tcell.Screen) {
  sv.Box.DrawForSubclass(screen, sv)
  x, y, width, height := sv.GetRect()

  for i, item := range sv.items {
    itemY := y + (i - sv.offsetY) * sv.itemSize
    if itemY >= y && itemY < height + y - sv.itemSize {
      item.view.SetRect(x, itemY, width, sv.itemSize)
      item.view.Draw(screen)
    }
  }
}

func(sv *ScrollView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
  return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
    switch event.Key(){
    case tcell.KeyDown:
      if sv.offsetY < len(sv.items)-1 {
        sv.offsetY++
      }
    case tcell.KeyUp:
      if sv.offsetY > 0 {
        sv.offsetY--
      }
    }
  }
}

func(sv *ScrollView) MouseHandler(
) func(
  action tview.MouseAction,
  event *tcell.EventMouse, 
  setFocus func(p tview.Primitive),
)( bool, tview.Primitive ){
  return func(
    action tview.MouseAction,
    event *tcell.EventMouse, 
    setFocus func(p tview.Primitive),
  )(bool, tview.Primitive){
    switch action {
    case tview.MouseScrollDown:
      _, _, _, height := sv.GetRect()
      offset := sv.offsetY + height/sv.itemSize - sv.itemSize
      if offset > len(sv.items) - sv.itemSize || sv.deactivated {
        return false, nil
      }
      sv.offsetY++
      return true, nil
    case tview.MouseScrollUp:
      if sv.offsetY-1 < 0 || sv.deactivated {
        return false, nil
      }
      sv.offsetY--
      return true, nil
    case tview.MouseLeftDown:
      sv.deactivated = true

      setFocus(sv)
      _, mouseY := event.Position()
      _, viewY, _, _ := sv.GetInnerRect()

      index := (mouseY - viewY) / sv.itemSize
      if index >= 0 && index < len(sv.items) - sv.offsetY {
        prevSelect := -1
        trueIndex  := index + sv.offsetY

        for i, item := range sv.items {
          if item.focused {
            prevSelect = i
            break
          }
        }
        if prevSelect == trueIndex {
          return true, nil
        } else if prevSelect == -1 {
          sv.items[trueIndex].SetFocus(true)
        } else {
          sv.items[prevSelect].SetFocus(false)
          sv.items[trueIndex].SetFocus(true)
        }
      }
      return true, nil
    case tview.MouseLeftUp:
      sv.deactivated = false
      return true, nil
    default:
      return false, nil
    }
  }
}
