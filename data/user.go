package data

import (
  _ "github.com/lib/pq"
  "database/sql"
  "log"
 )

type User struct {
   Id int
   Name string
   AccountKey string
 }

type UserDb interface {
  Close() error
  Authenticate(string) (User, error)
  RecordUsage(User)
}

// Implements the UserDB interface
type SqlUserDb struct {
  db *sql.DB
}

// Construct a new UserDB, connect to the database, etc
// It returns a connected UsetDB or an error if there is a problem connecting
func NewUserDb() (UserDb, error) {
  db, err := sql.Open("postgres", "user=geoff dbname=timezone sslmode=disable")
  err = db.Ping()
  return &SqlUserDb{db}, err
}

// Close the DB connection
func (udb *SqlUserDb) Close() error {
  return udb.db.Close()
}

// Authenticate a user
// Returns a User or an error if no user can be found with the token
func (udb *SqlUserDb) Authenticate(token string) (User, error) {
  var id int
  var name string
  err:= udb.db.QueryRow("SELECT id, accountname FROM account WHERE accountkey = $1", token).Scan(&id, &name)

  switch {
  case err == sql.ErrNoRows:
    log.Printf("No user with accountkey '%s'", token)
    return User{}, err
  case err != nil:
    log.Printf("Failed to exec query: %s", err)
    return User{}, err
  default:
    return User{id, name, token}, nil
  }
}

// Record the fact that the user used the service
func (udb *SqlUserDb) RecordUsage(u User) {
  log.Printf("User '%s:%d' used the service", u.Name, u.Id)
  _, err := udb.db.Exec("INSERT INTO account_usage (account_id) VALUES($1)", u.Id)
  if nil != err {
    log.Printf("Failed to record account usage: %s", err)
  }
}
