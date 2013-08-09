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

type UserDb struct {
  db *sql.DB
}

func NewUserDb() (UserDb, error) {
  db, err := sql.Open("postgres", "user=geoff dbname=timezone sslmode=disable")
  err = db.Ping()
  return UserDb{db}, err
}

func (udb *UserDb) Close() error {
  return udb.db.Close()
}

func (udb *UserDb) Authenticate(token string) (User, error) {
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

func (udb *UserDb) RecordUsage(u *User) {
  log.Printf("User '%s:%d' used the service", u.Name, u.Id)
  _, err := udb.db.Exec("INSERT INTO account_usage (account_id) VALUES($1)", u.Id)
  if nil != err {
    log.Printf("Failed to record account usage: %s", err)
  }
}
