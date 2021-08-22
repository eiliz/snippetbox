# Transactions

### Managing NULL values

Go is not good at managing NULL values in db records.

For example if a field contains a NULL value that is supposed to be converted
into a string that will fail. One solution is to change the field you're
scanning into from string to sql.NullString.

But the easiest is to simply avoid NULL values by making the columns NOT NULL
and using default values.


### Working with transactions

The calls to Exec, Query, QueryRow will use any available connection from the db
connection pool. That means that subsequent calls to Exec might not use the same
connection.

Sometimes that can be a problem, for example if you lock a table with MySQL's
LOCK TABLES command. Then you'd have to unlock it on the same connection but
that might not be possible and you'd get a deadlock.

The solution is to wrap multiple statements in a transaction, that will ensure
they'll use exactly the same connection.

```
type Model struct {
  DB *sql.DB
}

func (m *Model) Transaction() error {
  // Calling the Begin method on the connection pool creates
  // an sql.Tx object which represents the in-progress database transaction.
  tx, err:= m.DB.Begin()

  if err != nil {
    return err
  }

  // Calling Exec on the recently created transaction object and not on the connection pool.
  _, err = tx.Exec("INSERT INTO ...")
  if err != nil {
    // In case of any error calling tx.Rollback will abort the transaction
    // and no changes will be made to the database. This means that even if
    // this Exec would succeed while the next would fail, nothing would be
    // saved to the db as the transaction as a whole would fail.
    tx.Rollback()
    return err
  }

  _, err = tx.Exec("UPDATE ...)
  if err != nil {
    tx.Rollback()
    return err
  }

  // If there were no errors, the statements in the tx can be commited to the db.
  // You always have to make sure the connection is closed by either calling
  // tx.Rollback or tx.Commit before the function returns so that the connection
  // will be returned to the connection pool.
  err = tx.Commit()
  return err
}
```

Transactions are useful when you want to have the guarantee that multiple
statements are executed as an atomic action. That means that:
- either all statements succeed
- or if any of them fails, nothing gets commited to the DB, not even the
  successful ones


### Managing connections

The connection pool is made of idle and in-use connections. The default settings
have no limit on the number of concurrently open connections but the number of
idle connections is maxed at 2.

These settings can be changed with SetMaxOpenConns and SetMaxIdleConns.

There are also hard limits on the database itself. So for example even if the
connection pool has no limit for the number of open connections, MySQL does,
it's 151 by default. That means that you could get the database returning a "too
many connections" error under a high load.

If you set your connection pool limit under that, any request comming in while
there are no db connections left will be left waiting until one becomes free.

Sometimes that can be acceptable but for web apps it would be preferable to log
a "too many connections" error and send back a 500 Internal Server Error so that
the user is not left waiting and getting a timeout.

### Prepared statements

Exec, Query, and QueryRow use prepared statements to prevent SQL injections.
That means that that everytime you call them they setup a prepared statement on
the db connection, run with the parameters you give them and then close the
statement.

This might be considered ineficient to do on every run.

An alternative is to use DB.Prepare to create the prepared statement once and
store it for reuse. This would be particularly useful for complex statements
that are repeated a lot of times like multiple joins run on bulk inserts of
thousands of records.

```
type Model struct {
  DB *sql.DB
  InsertStmt *sql.Stmt
}

func NewModel(db *sql.DB) (*Model, error) {
  // Use the Prepare method to create a new prepared statement
  // for the current connection pool
  insertStmt, err := db.Prepare("INSERT INTO ...)
  if err != nil {
    return nil, err
  }

  return &Model{db, insertStmt}, nil
}

func (m *Model) Insert(args...) error {
  // Exec is run against the statement, not against the
  // connection pool like before
  _, err := m.InsertStmt.Exec(args...)

  return err
}

func main() {
  db, err := sql.Open(...)
  if err != nil {
    errorLog.Fatal(err)
  }

  defer db.Close()

  model, err := NewModel(db)
  if err != nil {
    errorLog.Fatal(err)
  }

  // The prepared statement must be closed before main returns
  defer model.InsertStmt.Close()
}
```

### Important to keep in mind

Prepared statements exist on db connections. The first time a prepared statement
is used it's attached to one of the db connections from the pool of connections.
The sql.Stmt object will remember this connection and attempt to use it again
next time. If it happens that the connection has been closed or is in use, the
statement has to be re-prepared on another connection.

Under heavy load, it is possible that many statements will be created on
multiple connections. This can lead to statements being prepared and re-prepared
multiple times or even running into server-side limits on the number of
statements (MySQL has a defualt limit of 16, 382 stmts).

For the simple use cases it's preferred to use the regular Exec, Query, QueryRow
methods that prepare the statements for you.
