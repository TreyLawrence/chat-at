package main

import (
	"chatat/db"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/parnurzeal/gorequest"
)

var (
	rq   = gorequest.New().SetDebug(true)
	host = "http://localhost:8080"
)

func TestIntegration(t *testing.T) {
	go main()
	defer db.Truncate()

	c1 := createConv(t, "idk")
	c2 := createConv(t, "idk2")
	fmt.Println("WTF", c1, c2)

	m11 := createMsg(t, c1, "msg11")
	m12 := createMsg(t, c1, "msg12")
	m21 := createMsg(t, c2, "msg21")
	m22 := createMsg(t, c2, "msg22")

	getC1 := getConv(t, c1.ID)
	if len(getC1.Messages) != 2 {
		t.Error()
	}
	if getC1.Messages[0].ID != m11.ID {
		t.Error()
	}
	if getC1.Messages[1].ID != m12.ID {
		t.Error()
	}

	getC2 := getConv(t, c2.ID)
	if len(getC2.Messages) != 2 {
		t.Error()
	}
	if getC2.Messages[0].ID != m21.ID {
		t.Error()
	}
	if getC2.Messages[1].ID != m22.ID {
		t.Error()
	}

	deleteMsg(t, c1, m11)
	getC1 = getConv(t, c1.ID)
	if len(getC1.Messages) != 1 {
		t.Error()
	}
	if getC1.Messages[0].ID != m12.ID {
		t.Error()
	}

	deleteConv(t, c1.ID)
	cs := getConvs(t)
	if len(cs) != 1 {
		t.Error()
	}
	if cs[0].ID != c2.ID {
		t.Error()
	}
}

func createConv(t *testing.T, subject string) *db.Conversation {
	resp, body, errs := rq.Post(host + "/conversations").
		Send(fmt.Sprintf(`{"subject": "%s"}`, subject)).
		End()
	if len(errs) > 0 {
		t.Error(errs)
	} else if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	c := db.Conversation{}
	if err := json.Unmarshal([]byte(body), &c); err != nil {
		t.Error(err)
	}
	return &c
}

func getConv(t *testing.T, id int) *db.Conversation {
	resp, body, errs := rq.Get(fmt.Sprintf("%s/conversations/%d", host, id)).
		End()
	if len(errs) > 0 {
		t.Error(errs)
	} else if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	c := db.Conversation{}
	if err := json.Unmarshal([]byte(body), &c); err != nil {
		t.Error(err)
	}
	return &c
}

func getConvs(t *testing.T) []*db.Conversation {
	resp, body, errs := rq.Get(fmt.Sprintf("%s/conversations", host)).
		End()
	if len(errs) > 0 {
		t.Error(errs)
	} else if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	cs := []*db.Conversation{}
	if err := json.Unmarshal([]byte(body), &cs); err != nil {
		t.Error(err)
	}
	return cs
}

func deleteConv(t *testing.T, id int) {
	resp, _, errs := rq.Delete(fmt.Sprintf("%s/conversations/%d", host, id)).
		End()
	if len(errs) > 0 {
		t.Error(errs)
	} else if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
}

func createMsg(t *testing.T, c *db.Conversation, msg string) *db.Message {
	resp, body, errs := rq.Post(fmt.Sprintf("%s/conversations/%d/messages",
		host, c.ID)).
		Send(fmt.Sprintf(`{"user_name": "trey", "txt": "%s"}`, msg)).
		End()
	if len(errs) > 0 {
		t.Error(errs)
	} else if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	m := db.Message{}
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		t.Error(err)
	}
	return &m
}

func deleteMsg(t *testing.T, c *db.Conversation, msg *db.Message) {
	resp, _, errs := rq.Delete(fmt.Sprintf("%s/conversations/%d/messages/%d",
		host, c.ID, msg.ID)).
		End()
	if len(errs) > 0 {
		t.Error(errs)
	} else if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
}
