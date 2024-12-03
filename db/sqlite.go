package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/TylerAldrich814/TermLink/utils"
	_ "modernc.org/sqlite"
)

// const DBSYNC = `
//   CREATE TABLE IF NOT EXISTS db_sync (
//     key         TEXT PRIMARY KEY,
//     last_synced DATETIME,
//   );
// `
// const CONVERSATION = `
//   CREATE TABLE IF NOT EXISTS conversations (
//     id                TEXT PRIMARY KEY,
//     created_at        DATETIME,
//     type              TEXT NOT NULL CHECK (type in ('direct', 'group')),
//     user_id           TEXT,
//     last_sent_user_id TEXT,
//     other_user_id     TEXT,
//   );
// `
// const MESSAGES = `
//   CREATE TABLE IF NOT EXISTS messages (
//     id              TEXT     PRIMARY KEY,
//     conversation_id TEXT,
//     sender_id       TEXT,
//     content         TEXT,
//     created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
//   );
// `
// const CONTACT = `
//   CREATE TABLE IF NOT EXISTS contact (
//     id       TEXT PRIMARY KEY,
//     username TEXT,
//     status   TEXT NOT NULL CHECK (status in ('active', 'inactive', 'away'))
//   );
// `
const (
  schemaDir = "./schemas"
)

func InitSQLite()( *sql.DB, error ){ 
  db, err := sql.Open("sqlite", "local.db")
  if err != nil {
    return nil, fmt.Errorf("Failed to open SQLite Database: %w", err)
  }
  if err := loadSchemas(db); err != nil {
    return nil, err
  }

  return db, nil
}
func loadSchemas(sql *sql.DB) error {
  files, err := ioutil.ReadDir(schemaDir)
  if err != nil {
    utils.Error("SQLite::loadSchemas - Failed to load from ./schemas directory")
    return errors.New("Failed to load SQL Schemas")
  }

  for _, file := range files {
    if filepath.Ext(file.Name()) != ".sql" {
      continue
    }
    sqlPath := filepath.Join(schemaDir, file.Name())
    sqlContent, err := os.ReadFile(sqlPath)
    if err != nil {
      return fmt.Errorf("Failed to load SQL Schema %s\n -- %w", sqlPath, err)
    }
    _, err = sql.Exec(string(sqlContent))
    if err != nil {
      return fmt.Errorf("Failed to Execute SQL File: %s\n -- %w", sqlPath, err)
    }
  }

  return nil
}


func(db *Database) updateSyncedKey(key string) error {
  _, err := db.sqlite.Exec(`
    INSERT OR IGNORE INTO db_sync (key, last_synced)
    VALUES (?, ?)
  `, key, time.Now())
  if err != nil {
    utils.Error("SQLite::updateSyncedKey: Failed to Insert Sync Update: %w", err)
    return errors.New("Failed to Sync Database")
  }
  return nil
}

// Returns the time of the last Sync event for the provided Key
func(db *Database) getLastSynced(key string)( *time.Time, error ){
  rows, err := db.sqlite.Query(`
    SELECT *
    FROM db_sync
    WHERE key = ?
  `, key)
  if err != nil {
    utils.Error("SQLite::getLastSynced - Failed to Query from DB_SYNC: %w", err)
    return nil, errors.New("Failed to Fetch last Synced Event")
  }

  defer rows.Close()
  for rows.Next() {
    var lastSynced time.Time
    if err := rows.Scan(&lastSynced); err != nil {
      utils.Error("SQLite::getLastSynced - Failed to Scan Row: %w", err)
      return nil, errors.New("Failed to Fetch Last Synced Event")
    }
    return &lastSynced, nil
  }

  return nil, errors.New("Failed to find Last Synced Event")
}

// Fetch Converstations from Supabase and Sync them with SQLite
// - First, we query the last synced event time for conversations
// - Next, we fetch for any new converstations from Supbase based
//   on whether the created_at time is greater than lastSynced.
// - If any new conversations were found. We then add them to our
//   conversations SQlite Databse Table.
// - Finally, we update the last synced event time for conversations.
func(db *Database) fetchAndSyncConversations() error {
  lastSynced, err := db.getLastSynced("conversations")
  if err != nil { return err }

  res, _, err := db.supabase.From("converstations").
    Select("*", "", true).
    Eq("user_id", db.userID).
    Gt("created_at", lastSynced.Format(time.RFC3339)).
    Execute()
  if err != nil {
    utils.Error("Failed to fetch Conversations from Supabase: %w", err)
    return errors.New("Supabase Query Error")
  }

  var conversations []Conversations
  if err := json.Unmarshal(res, &conversations); err != nil {
    utils.Error("Failed to Marshal Supabase Conversations Query: %w")
    return errors.New("Supabase Query Unmarshal Error")
  }

  for _, c := range conversations {
    _, err := db.sqlite.Exec(`
      INSERT OR IGNORE INTO conversations (id, created_at, type, user_id, last_sent_user_id, other_user_id, synced)
      VALUES (?, ?, ?, ?, ?, 1)
    )
    `, c.ID, c.CreatedAt, c.Type, c.UserID, c.LastSentUserID, c.OtherUserID, 1)
    if err != nil {
      utils.Error("Failed to insert converstation into SQLite: %w", err)
      return errors.New("SQLite Sync Error - Conversation")
    }
  }
  if err := db.updateSyncedKey("conversations"); err != nil {
    return err
  }
  return nil
}

// Fetchees and Syncs all User Messages from Supabase.
// - First, we obtain all of the user's Conversation ID's
// - Next, we obtain the previous sync event time for messages
// - Next, we fetch all messages with a created_at time greater  
//   than lastSynced.
// - Next, if messages were found. We add them to our SQLite Messages Table..
// - Finally, we update the last synced event time for our messages table.
func(db *Database) fetchAndSyncMessages() error {
  // - Query all Synced Conversations from SQLite and Obtain all synced ConverstaionIDs
  rows, err := db.sqlite.Query(`
    SELECT conversation_id
    FROM converstations
    WHERE synced = 1
  `)
  if err != nil {
    utils.Error("SQLite::fetchAndSyncMessages - Failed to Query Conversation ID's: %w", err)
    return errors.New("Failed to Query Local Database")
  }
  defer rows.Close()

  var conversationIDs []string
  for rows.Next() {
    var conversationID string
    if err := rows.Scan(&conversationID); err != nil {
      utils.Error("SQLite::fetchAndSyncMessages - Failed to Scan row: %w", err)
      return errors.New("SQLite Row Scan Error")
    }
    conversationIDs = append(conversationIDs, conversationID)
  }

  // - Iterate over ConverstationIDs and fetchAndSync All Messages belonging to 
  //   each Conversation
  lastSynced, err := db.getLastSynced("messages")
  if err != nil {
    return err
  }
  for _, id := range conversationIDs {
    data, _, err := db.supabase.
      From("messages").
      Select("*", "", true).
      Eq("conversation_id", id).
      Gt("created_at", lastSynced.Format(time.RFC3339)).
      Execute()
    if err != nil {
      utils.Error("SQLite::fetchAndSyncMessages - Failed to fetch message from DB: %w", err)
      return errors.New("Failed to Fetch for Message")
    }
    var messages []Messages
    if err := json.Unmarshal(data, &messages); err != nil {
      utils.Error("SQLite::fetchAndSyncMessages - Failed to Mashal Conversation Messages: %w", err)
      return errors.New("Failed to fetch Messages")
    }
    
    // - If messages contains any data, we then store the messages in SQLite
    if len(messages) == 0 {
      continue;
    }
    for _, m := range messages {
      _, err := db.sqlite.Exec(`
        INSERT OR IGNORE INTO messages(id, conversation_id, sender_id, content, created_at, read)
        VALUES (?, ?, ?, ?, ?, ?)
      `, m.ID, m.ConversationID, m.SenderID, m.Content, m.CreatedAT, m.Read)
      if err != nil {
        utils.Error("SQLite::fetchAndSyncMessages - Failed to Store new Message: %w", err)
        return errors.New("Messages Failed to Sync")
      }
    }
  }
  // - Lastly, we'll store the time for our successful Database Sync Event
  utils.Log("SQLite::fetchAndSyncMessages - Successully Synced Messages")
  if err := db.updateSyncedKey("messages"); err != nil {
    return err
  }
  return nil
}

func(db *Database) syncAndUpdateContacts() error {
  lastSynced, err := db.getLastSynced("contacts")
  if err != nil { return err }

  // -- Query Supabase for the current user's Contacts 
  // -- Returns an array of user_ids
  res, _, err := db.supabase.
    From("contacts").
    Select("contacts", "", false).
    Eq("user_id", db.userID).
    Gt("updated_at", lastSynced.Format(time.RFC3339)).
    Execute()
  if err != nil {
    utils.Error("SQLite::syncAndUpdateContacts - Failed to fetch User Contacts")
    return errors.New("Failed to Query Contacts from Supabase")
  }

  var contactIds []struct{ user_id string  }
  if err := json.Unmarshal(res, &contactIds); err != nil {
    utils.Error("SQLite::syncAndUpdateContacts - Failed to Unmarshal Contacts - %w", err)
    return errors.New("Failed to Unmarshal Contacts")
  }

  if len(contactIds) == 0 {
    return nil
  }

  for _, contact := range contactIds {
    res, _, err := db.supabase.
      From("users").
      Select("username,user_id,status", "", false).
      Eq("user_id", contact.user_id).
      Execute()
    if err != nil {
      utils.Error("SQLite::syncAndUpdateContacts - Failed to fetch User - %w", err)
      return errors.New("Failed to Fetch TermLink User")
    }
    var user struct {
      username  string
      user_id   string
      status    Status
    }
    if err := json.Unmarshal(res, &user); err != nil {
      utils.Error("SQLite::syncAndUpdateContacts - Failed to Unmarshal User - %w", err)
      return errors.New("Failed to Fetch TermLink User")
    }

    _, err = db.sqlite.Exec(`
      INSERT INTO contacts (id, username, status)
      VALUES (?, ?, ?)
      ON CONFLICT(id) DO UPDATE SET
        username = excluded.username
        status = excluded.status;
    `, user.user_id, user.username, user.status)
    if err != nil {
      utils.Error("SQLite::syncAndUpdateContacts - Failed to sync User - %w", err)
      return errors.New("Failed to sync TermLink User")
    }
  }

  db.updateSyncedKey("contacts")

  return nil
}
