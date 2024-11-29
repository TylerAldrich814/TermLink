package utils

import "fmt"

type BulletItem struct {
  text  string
  left  string
  right string
  fg    string
  bg    string
}
func(b *BulletItem) Text() string {
  return b.text
}
func(b *BulletItem) Length() int {
  return len(b.text)+2
}
func(b *BulletItem) View() string {
  content :=  fmt.Sprintf(
    "[%s:%s]%s[-]", 
    b.bg,
    b.fg,
    b.text,
  )
  return fmt.Sprintf(
    "%s%s%s[-]", 
    b.left, 
    content,
    b.right,
  )
}

func Bullet(
  content    string,
  foreground string,
  background string,
) * BulletItem {
  return &BulletItem{
    text  : content,
    left  : fmt.Sprintf("[%s:%s]◖[-]", foreground, background),
    right : fmt.Sprintf("[%s:%s]◗[-]", foreground, background),
    fg    : foreground,
    bg    : background,
  }
}
