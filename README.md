# chat-at
Simple rest API for conversations

There are two resources for this API: conversations and messages. Both
resources have create, read, and delete. The conversation endpoints are:
```
Create:     POST /conversations
Get all:    GET /conversations
Get one:    GET /conversations/{conversation_id}
Delete one: DELETE /conversations/{conversation_id}
```

When creating a conversation, you must pass in a json object with a `subject` field.
This subject field must be unique, otherwise a status of 409 will be returned.

And since a conversation is a group of messages, the message endpoints are hierarchical:
```
Create:     POST /conversations/{conversation_id}/messages
Get all:    GET /conversations/{conversation_id}/messages
Get one:    GET /conversations/{conversation_id}/messages/{message_id}
Delete one: DELETE /conversations/{conversation_id}/messages/{message_id}
```

When creating a message, you must pass in a json object with two required fields:
`user_name` and `txt`. There currently isn't an authentication layer, so `user_name` uniqueness is not yet supported.

### Running
Make sure to set your `GOPATH` first, and add `$GOPATH/bin` to your `PATH`. You will also need to install postgres if you want to run locally.

Download this project `go get github.com/TreyLawrence/chat-at` and `cd $GOPATH/src/github.com/TreyLawrence/chat-at`

Save the postgres url of your local database in the `DATABASE` environment variable.

To run tests `go test chatat`

To run locally `go install chatat && chatat`

### TODO
Add user authentication and websocket API
