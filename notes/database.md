# Connecting to a database

We're using a MySQL driver along with the database/sql package from the std
library.

```
“import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)”
```

We assign the driver import to _ because main doesn't actually use it however we
need the init of the package to run so that it can register itself with the
database/sql package.

To connect to the db we use a database source name (DSN) that's customized for
the specific driver we're using: `web:pass@/snippetbox?parseTime=true`.

parseTime instructs the driver to convert SQL TIME and DATE fields to go
time.Time objects.

```
db, err := sql.Open("mysql", "web:pass@/snippetbox?parseTime=true")
```

sql.Open returns an sql.DB object which is not a connection but a pool of many
connections to hold future connections. The actual connections are established
lazily, the first moment one's needed.

The connection pool is safe for concurrent access and it's intended to be
long-lived that's why it should be initialized in main and passed on to the
handlers. Should not be init in handlers.

To check everything work we use the db.Ping method that creates a connection and
checks for errors.

### Prepared statements

In database management systems (DBMS), a prepared statement or parameterized statement is a feature used to execute the same or similar database statements repeatedly with high efficiency. Typically used with SQL statements such as queries or updates, the prepared statement takes the form of a template into which certain constant values are substituted during each execution.

The typical workflow of using a prepared statement is as follows:

Prepare: At first, the application creates the statement template and sends it
to the DBMS. Certain values are left unspecified, called parameters,
placeholders or bind variables (labelled "?" below):

```INSERT INTO products (name, price) VALUES (?, ?);```

Then, the DBMS compiles (parses, optimizes and translates) the statement template, and stores the result without executing it.
Execute: At a later time, the application supplies (or binds) values for the parameters of the statement template, and the DBMS executes the statement (possibly returning a result). The application may execute the statement as many times as it wants with different values. In the above example, it might initially supply "bike" for the first parameter and "10900" for the second parameter, and then later supply "shoes" for the first parameter and "7400" for the second parameter.
As compared to executing statements directly, prepared statements offer two main
advantages:
* The overhead of compiling the statement is incurred only once, although the statement is executed multiple times. However not all optimization can be performed at the time the statement template is compiled, for two reasons: the best plan may depend on the specific values of the parameters, and the best plan may change as tables and indexes change over time.

* Prepared statements are resilient against SQL injection because values which are transmitted later using a different protocol are not compiled like the statement template. If the statement template is not derived from external input, SQL injection cannot occur.