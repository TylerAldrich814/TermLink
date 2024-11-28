package main

import (
	"fmt"
	"os"

	"github.com/TylerAldrich814/TermLink/db"
	"github.com/TylerAldrich814/TermLink/tui"
	"github.com/TylerAldrich814/TermLink/utils"
	"github.com/joho/godotenv"
)


func main(){
  err := godotenv.Load()
  if err != nil {
    panic(fmt.Sprintf("Error loading .env file: %v", err))
  }
  build := utils.Mode(os.Getenv("BUILD"))

  db, err := db.InitSupbase(
    os.Getenv("SUPABASE_URL"),
    os.Getenv("SUPABASE_ANON"),
  )
  if err != nil {
    panic(fmt.Sprintf("Failed to create Supabase Client: %v", err))
  }

  tui.GetTermLinkTUI(
    build, db,
  ).
    GeneratePages().
    HandleInput().
    Start()
}
