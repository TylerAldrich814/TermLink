package utils


type Page uint

const (
  Auth Page = iota+1
  Inbox
  Contacts
  Settings
)

func(p Page)String() string {
  return [...]string{
    "Auth",
    "Inbox",
    "Contacts",
    "Settings",
  }[p-1]
}
func(p Page)Idx() uint {
  return uint(p)
}
