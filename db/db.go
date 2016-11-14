package db

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"time"

	pg "gopkg.in/pg.v5"
)

type Conversation struct {
	TableName struct{}  `sql:"conversations,alias:c" json:"-"`
	ID        int       `db:"id" json:"id"`
	Subject   string    `db:"subject" json:"subject" binding:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Messages  []Message `json:"messages,omitempty"`
}

type Message struct {
	TableName      struct{}      `sql:"messages,alias:m" json:"-"`
	ID             int           `db:"id" json:"id"`
	UserName       string        `db:"user_name" json:"user_name" binding:"required"`
	Txt            string        `db:"txt" json:"txt" binding:"required"`
	CreatedAt      time.Time     `db:"created_at" json:"created_at"`
	ConversationID int           `db:"conversation_id" json:"conversation_id"`
	Conversation   *Conversation `db:"conversation" json:"conversation,omitempty"`
}

var db *pg.DB

func init() {
	pg.SetQueryLogger(log.New(os.Stdout, "", log.LstdFlags))

	url, err := url.Parse(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	} else if url.Host == "" {
		panic("Invalid host")
	} else if url.EscapedPath() == "" {
		panic("Invalid database name")
	}

	o := &pg.Options{
		Addr:     url.Host,
		Database: url.EscapedPath()[1:],
	}
	if url.User != nil {
		pwd, _ := url.User.Password()
		o.User = url.User.Username()
		o.Password = pwd
	} else {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		o.User = u.Username
	}
	db = pg.Connect(o)

	createTables := []string{
		`create table if not exists conversations (
			id serial primary key,
			subject text unique not null,
			created_at timestamp without time zone not null default now()
		)`,
		`create table if not exists messages (
			id serial primary key,
			conversation_id integer references conversations on delete cascade,
			user_name text not null,
			txt text not null,
			created_at timestamp without time zone not null default now()
		)`,
	}
	for _, sql := range createTables {
		if _, err := db.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// NOTE(trey): _only_ use this for testing
func Truncate() {
	for _, table := range []string{"messages", "conversations"} {
		if _, err := db.Exec(fmt.Sprintf("truncate %s cascade", table)); err != nil {
			panic(err)
		}
	}
}

func (c *Conversation) Insert() error { return db.Insert(c) }
func (c *Conversation) Delete() error { return db.Delete(c) }
func (m *Message) Insert() error      { return db.Insert(m) }
func (m *Message) Delete() error      { return db.Delete(m) }

func Conversations() ([]*Conversation, error) {
	convs := []*Conversation{}
	err := db.Model(&convs).
		Column("c.id", "c.subject", "c.created_at", "Messages").
		Select()
	return convs, err
}

func GetConversation(id int, withMessages bool) (*Conversation, error) {
	c := &Conversation{}
	cols := []string{"c.id", "c.subject", "c.created_at"}
	if withMessages {
		cols = append(cols, "Messages")
	}
	err := db.Model(&c).
		Column(cols...).
		Where("c.id = ?", id).
		Select()
	return c, err
}

func GetMessage(id int) (*Message, error) {
	c := &Message{}
	err := db.Model(&c).
		Column("m.id", "m.user_name", "m.txt", "m.created_at",
			"m.conversation_id", "Conversation").
		Where("m.id = ?", id).
		Select()
	return c, err
}
