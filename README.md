[![License BSD](https://img.shields.io/badge/License-BSD-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/nats-io/go-nats)](https://goreportcard.com/badge/github.com/bradclawsie/rscs)
[![GoDoc](https://godoc.org/github.com/bradclawsie/rscs?status.svg)](http://godoc.org/github.com/bradclawsie/rscs)

# The Ridiculously Simple Configuration System

**RSCS** is a configuration database that is ridiculously simple. 

**RSCS** allows you to store text values for keys in a SQLite
database. That is *all* it does. You will not find any of the features
offered by more complex clustered configuration databases. 

### What's wrong with Consul and Vault?

Nothing! I've used them in the past. I simply find they are too much
for my requirements and I don't want deal with complexity I will
never use. **RSCS** is for people who want a single, simple, available
local key/value store and nothing more. 

### Why do you use SQLite instead of $X?

SQLite is a proven technology that provides that features needed by
**RSCS** and none of the features that are not needed. SQLite provides
a better, transactional equivalent of `fopen`, which is the level of
simplicity and reliability suitable here. Furthermore, if you decide
you don't like **RSCS** anymore, you can take your database file and
use some other SQLite-supporting tool with it.

### Are there a bunch of complicated tables?

```
CREATE TABLE kv (key VARCHAR(255) PRIMARY KEY, value TEXT NOT NULL)
```

That is it. If you want more, extend the codebase yourself.

### How do you achieve clustering? Do you support the Raft protocol?

There is no clustering. If you want clustering, you can build it
on top of **RSCS** because the codebase is intended to be very simple
to read and understand in just a few hours.

### Can I store passwords? How are they protected?

You can store whatever you want (within the limits of SQLite)
and secure your data however you wish. If you want a value to be
encrypted, encrypt it. If you don't want people snooping at the SQLite
file, use user/group file permissioning to give you the level of
security you need. If you want to support temporary passwords or some
other nifty Vault feature, extend the code and manage the rows yourself.

### What about binary data?

Encode it as text using base64 or another textual encoding.

### How is output delivered?

JSON which is trivial: `{"Value":"your value"}`. Extend it if you want.

### The default daemon is http! Yuck!

**RSCS** is intended to be run on your local machine and not accept
external traffic. If you still believe you want the extra assurances
of https, then go in to the `rscs.go` file and change it.

### You keep saying "change the code"...

Yes. The **RSCS** codebase is intended to be very simple and only
provide the most basic features. You should read the code instead of
relying on `godoc`. If you want something more complex, just take
ownership of your fork and make changes to the source. You will be far
happier with source-code modifications than having to master tuning a large
number of optional features. This was what drove me away from Consul
and Vault...they are both too complex and seem intent on solving
everyone's problem from one code base. I am confident that even a
novice Go programmer can modify the **RSCS** codebase to suit their
particular needs.

### I read the code...there's almost nothing there. What's the point?

The point is to provide the simplest tool that covers basic needs and
gives you a clear path to extending it for your needs.

### I tried compiling it and it didn't build! WTF!

You must use Go 1.8 or higher.

### Why do you use a router like Chi if you want to keep things "ridiculously simple"?

[Chi](https://github.com/pressly/chi) is a very simple router that
allows our codebase to be much more compact and readable without
imposing a complex framework model.

### Okay, I get it, just show me how it is used.

*create an empty db:*

`$ rscs --db=/tmp/test.sqlite3 --create-only`

*run the daemon:*

`$ rscs --db=/tmp/test.sqlite3`

*now do some transactions:*

*status:*

`$ curl http://localhost:8081/v1/status`

output:

`{"Alive":true,"DBFile":"/tmp/test.sqlite3","Uptime":"4.30598268s"}`

*create a new row:*

`curl -X POST -d '{"Value":"value1"}' http://localhost:8081/v1/kv/key1`

*read it:*

`curl -X GET http://localhost:8081/v1/kv/key1`

output:

`{"Value":"value1"}`

*update:*

`curl -X PUT -d '{"Value":"value1-new"}' http://localhost:8081/v1/kv/key1`

*read it:*

`curl -X GET http://localhost:8081/v1/kv/key1`

output:

`{"Value":"value1-new"}`

*delete:*

`curl -X DELETE http://localhost:8081/v1/kv/key1`

(and test)

`curl -X GET http://localhost:8081/v1/kv/key1`

output:

`no value found`

*stop the daemon:*

Send any signal to the process that satisfies `os.Interrupt` (on Linux
systems, `SIGINT`). The daemon uses the graceful stopping feature made
available in Go 1.8 in the standard library.

