package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/rivo/tview"
)

type LogMessage struct {
  level   DebugLevel
  message string
}

type DebugMessage struct {
  time    time.Time
  message string
  level   DebugLevel
}

func( d *DebugMessage) String() string {
  return fmt.Sprintf(
    "%s %s %s", 
    d.time.Format("Jan 1 2024, 10:10:10:10"), 
    d.message, d.level.String(),
  )
}

type DebugWindow struct {
  app        *tview.Application
  View       *tview.TextView
  Messages   []*DebugMessage
  LogChannel chan LogMessage
  stopChan   chan struct{}
  stopOnce   sync.Once
}

var (
  instance  *DebugWindow
  startOnce sync.Once
)

func InitializeDebugWindow(
  app  *tview.Application,
  mode Build,
) {
  if mode != DevBuild {
    return
  }
  startOnce.Do(func(){
    view := tview.NewTextView().
      SetText("").
      SetTextAlign(tview.AlignTop).
      SetDynamicColors(true)

    instance = &DebugWindow{
      app        : app,
      View       : view,
      Messages   : []*DebugMessage{},
      LogChannel : make(chan LogMessage, 100),
      stopChan   : make(chan struct{}),
    }
    go instance.ProcessLogs()
  })
}

func GetInstance() *DebugWindow {
  return instance
}
func(d *DebugWindow) ProcessLogs(){
  for {
    select {
    case log := <-d.LogChannel:
      d.displayLog(log)
    case <-d.stopChan:
      close(d.LogChannel)
      return
    }
  }
}

func( d *DebugWindow) UpdateUI(){
  var Messages string = ""

  for i, msg := range d.Messages {
    Messages += fmt.Sprintf(
      "%s|%d| %s %s\n",
      msg.level.Tag(),
      i,
      msg.time.Format("15:04:05:00"),
      msg.message,
    )
  }

  d.app.QueueUpdateDraw(func(){
    d.View.SetText(Messages)
  })
}

func( d *DebugWindow ) displayLog(
  log LogMessage,
) {
  d.Messages = append(d.Messages, &DebugMessage {
    time    : time.Now(),
    message : log.message,
    level   : log.level,
  })
  d.UpdateUI()
}

func(d *DebugWindow) Stop() {
  if d.stopChan != nil {
    close(d.stopChan)
  }
}

func(d *DebugWindow) log(message LogMessage) {
  d.LogChannel <-message
}

// Logs a Regular Debug message by passing a LogMessage into
// The Global DebugWindow instance singleton. 
func Log(format string, f ...any){
  if instance == nil { return } 
  GetInstance().log(LogMessage{
    DebugLog,
    fmt.Sprintf(format, f...),
  })
}
// Logs a Warning Debug message by passing a LogMessage into
// The Global DebugWindow instance singleton. 
func Warn(format string, f ...any){
  if instance == nil { return } 
  GetInstance().log(LogMessage{
    DebugWarn,
    fmt.Sprintf(format, f...),
  })
}
// Logs an Error Debug message by passing a LogMessage into
// The Global DebugWindow instance singleton. 
func Error(format string, f ...any){
  if instance == nil { return } 
  GetInstance().log(LogMessage{
    DebugError,
    fmt.Sprintf(format, f...),
  })
}
