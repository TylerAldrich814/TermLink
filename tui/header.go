package tui


type Header struct {
  app  *TermLinkTUI

  kill chan struct{}
}
