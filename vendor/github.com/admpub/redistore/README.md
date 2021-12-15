# redistore

A session store backend for admpub/sessions - [src](https://github.com/admpub/sessions).

## Requirements

Depends on the [Redigo](https://github.com/garyburd/redigo) Redis library.

## Installation

    go get gopkg.in/boj/redistore.v1

## Documentation

Available on [godoc.org](http://www.godoc.org/gopkg.in/boj/redistore.v1).

See http://www.gorillatoolkit.org/pkg/sessions for full documentation on underlying interface.

### Example
    //import "github.com/webx-top/echo"
    //var ctx echo.Context


    // Fetch new store.
    store, err := NewRediStore(10, "tcp", ":6379", "", []byte("secret-key"))
    if err != nil {
        panic(err)
    }
    defer store.Close()

    // Get a session.
	session, err = store.Get(ctx, "session-key")
	if err != nil {
        log.Error(err.Error())
    }

    // Add a value.
    session.Values["foo"] = "bar"

    // Save.
    if err = sessions.Save(ctx); err != nil {
        t.Fatalf("Error saving session: %v", err)
    }

    // Delete session.
    session.Options.MaxAge = -1
    if err = sessions.Save(ctx); err != nil {
        t.Fatalf("Error saving session: %v", err)
    }

    // Change session storage configuration for MaxAge = 10 days.
    store.SetMaxAge(10*24*3600)

