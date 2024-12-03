package utils

import "fmt"

type BulletKind uint
const (
  FullBlock BulletKind = iota+1
  FullBullet
  LeftBullet
  RightBullet
)
func(k BulletKind)String() string {
  return [...]string{
    "FullBlock",
    "FullBullet",
    "LeftBullet",
    "RightBullet",
  }[k-1]
}

type BulletItem struct {
  kind  BulletKind
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
  var (
    l string
    r string
  )
  switch b.kind.String() {
  case "FullBlock":
    l = Block
    r = Block
  case "FullBullet":
    l = LeftCircle
    r = RightCircle
  case "LeftBullet":
    l = LeftCircle
    r = Block
  case "RightBullet":
    l = Block
    r = RightCircle
  }

  left  := fmt.Sprintf(
    "[%s:%s]%s[-]", b.fg, b.bg, l)
  right := fmt.Sprintf("[%s:%s]%s[-]", b.fg, b.bg, r)

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

func bulletBuilder(
  kind       BulletKind,
  content    string,
  foreground string,
  background string,
) *BulletItem {
  b := &BulletItem{
    kind : kind,
    text : content,
    fg   : foreground,
    bg   : background,
    item : "",
  }
  b.Build()

  return b
}

func GetBullet(
  content    string,
  foreground string,
  background string,
) * BulletItem {
  return bulletBuilder(
    FullBullet,
    content,
    foreground,
    background,
  )
}

func GetBlock(
  content    string,
  foreground string,
  background string,
) *BulletItem {
  return bulletBuilder(
    FullBlock,
    content,
    foreground,
    background,
  )
}

func BuildBullet(
  kind       BulletKind,
  content    string,
  foreground string,
  background string,
) *BulletItem {
  b := &BulletItem{
    kind : kind,
    text : content,
    fg   : foreground,
    bg   : background,
    item : "",
  }
  b.Build()

  return b
}
