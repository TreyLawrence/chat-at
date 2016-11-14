package controllers

import (
	"chatat/db"
	"errors"
	"net/http"
	"strconv"

	gin "gopkg.in/gin-gonic/gin.v1"
	pg "gopkg.in/pg.v5"
)

type (
	Conversations struct{}
	Messages      struct{}
)

func (Conversations) CreateHandler(c *gin.Context) {
	conv := db.Conversation{}
	if err := c.BindJSON(&conv); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := conv.Insert(); err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.Field('C') == "23505" {
			c.AbortWithError(http.StatusConflict, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, conv)
	}
}

func (Conversations) ListHandler(c *gin.Context) {
	if convs, err := db.Conversations(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, convs)
	}
}

func (Conversations) TakeHandler(c *gin.Context) {
	conv, err := conversationFromPath(c, true)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, conv)
}

func (Conversations) DeleteHandler(c *gin.Context) {
	conv, err := conversationFromPath(c, false)
	if err != nil {
		return
	}
	if err := conv.Delete(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (Messages) CreateHandler(c *gin.Context) {
	conv, err := conversationFromPath(c, false)
	if err != nil {
		return
	}
	m := db.Message{}
	if err := c.BindJSON(&m); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	m.ConversationID = conv.ID
	if err := m.Insert(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (Messages) ListHandler(c *gin.Context) {
	conv, err := conversationFromPath(c, true)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, conv.Messages)
}

func (Messages) TakeHandler(c *gin.Context) {
	conv, err := conversationFromPath(c, false)
	if err != nil {
		return
	}
	m, err := messageFromPath(c)
	if err != nil {
		return
	}
	if m.Conversation.ID != conv.ID {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, m)
}

func (Messages) DeleteHandler(c *gin.Context) {
	conv, err := conversationFromPath(c, false)
	if err != nil {
		return
	}
	m, err := messageFromPath(c)
	if err != nil {
		return
	}
	if m.Conversation.ID != conv.ID {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := m.Delete(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func messageFromPath(c *gin.Context) (*db.Message, error) {
	id, err := idFromPath(c, "messages")
	if err != nil {
		return nil, err
	}
	m, err := db.GetMessage(id)
	if err == pg.ErrNoRows {
		c.AbortWithStatus(http.StatusNotFound)
		return nil, err
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	return m, nil
}

func conversationFromPath(c *gin.Context, withMessages bool) (*db.Conversation, error) {
	id, err := idFromPath(c, "conversations")
	if err != nil {
		return nil, err
	}
	conv, err := db.GetConversation(id, withMessages)
	if err == pg.ErrNoRows {
		c.AbortWithStatus(http.StatusNotFound)
		return nil, err
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil, err
	}
	return conv, nil
}

func idFromPath(c *gin.Context, resourceName string) (int, error) {
	idStr, ok := c.Params.Get(resourceName)
	if !ok {
		err := errors.New(resourceName + " path param not found")
		c.AbortWithError(http.StatusBadRequest, err)
		return 0, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return 0, err
	}
	return id, nil
}
