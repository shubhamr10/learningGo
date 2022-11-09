# Bread and Breakfast Application
github.com/shubhamr10/learningGo

## Planning
* Deciding what to build   : Booking & Reservations
* Project Scope            : A Bed & Breakfast with two rooms
* Key functionality        : What do we need to do ?

### Key Functionality
* Showcase the property
* Allow for booking a room for one or more nights
* Check a room's availability
* Book the room
* Notify guest, and notify property owner  

* Have a backend that owner logs into
* Renew existing bookings
* Show a calendar of bookings
* Change or cancel a booking

### What will we need ?
* An authentication system
* Somewhere to store information (database)
* A means of sending notification (email/text)

### Other templating engine
https://github.com/CloudyKit/jet

## See your test coverage in details
go test -coverprofile=coverage.out && go tool cover -html=coverage.out

# Database integration

We are using :postgres as a database through docker.
Dbeaver is the CLI to connect to the postgres as it is a lightweight one.

We are using "Soda" a.k.a "Pop" by buffalo to create our tables using migrations.


```text
soda generate fizz <table_name>

soda migrate

soda migrate down

```

indices make search faster.


soda reset
This will delete all the tables and create again
I mean it runs all down migrations and then up migrations, but no client should be connected
to database. not even Dbeaver.

# Object Relational Model for db
1. GORM
2. upper/db

Don't use these because it slows down your code.


```go
// this function allows to cancel the transaction if there is any problem
// In case of the web application it is very unpredictable that things would work in a 
// better way
context.WithTimeout()
db.ExecContext
```

soda generate sql <filename>
generates a sql migration file