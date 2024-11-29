package utils

import "fmt"

const (
  LeftCircle  = "◖"
  RightCircle = "◗"
)

type BulletItem struct {
  item  string
  text  string
  fg    string
  bg    string
} 
func(b *BulletItem) Item() string {
  return b.item
}
func(b *BulletItem) Text() string {
  return b.text
}
func(b *BulletItem) Length() int {
  return len(b.text)+2
}
func(b *BulletItem) Build() {
  left  := fmt.Sprintf(
    "[%s:%s]◖[-]", b.fg, b.bg)
  right := fmt.Sprintf("[%s:%s]◗[-]", b.fg, b.bg)

  content :=  fmt.Sprintf(
    "[%s:%s]%s[-]", 
    b.bg,
    b.fg,
    b.text,
  )
  item := fmt.Sprintf(
    "%s%s%s[-]", 
    left, 
    content,
    right,
  )
  b.item = item
}
func(b *BulletItem) UpdateColors(
  fg      string,
  bg      string,
) {
  if len(fg) != 0 {
    b.fg = fg
  }
  if len(bg) != 0 {
    b.bg = bg
  }
}
func(b *BulletItem) UpdateContent(
  content string,
) {
  b.text = content
}

func Bullet(
  content    string,
  foreground string,
  background string,
) * BulletItem {
  b := &BulletItem{
    text : content,
    fg   : foreground,
    bg   : background,
    item : "",
  }
  b.Build()

  return b
}
