# Leveled logging

Logging information should be managed differently in terms of how much detail
should be included and what types of logging should be active in which scenario.


There are various levels like fatal (complete shutdown of the app), error
(something serious), warning (something took a few tries to get it to work,
should be looked at), info (a new user was created in the db), debug, trace
(include the most granular details you can).


Usually you'd have a logging framework and your code would make logging calls to
this framework. Based on the env settings these calls might do nothing. For
example you could have calls to trace and info inside a function and trace would
only work of the env was dev, and info would work if the env was production.


More detailed logging means more resources used so it's a balancing act to setup
the correct amount of logging for the specific env, as well as the correct
amount of granularity. It's a good idea to include as much info as possible and
make the levels manageable from the settings to make it possible to adjust the
logging levels in production.


“During development, it’s easy to view the log output because the standard streams are displayed in the terminal.

In staging or production environments, you can redirect the streams to a final destination for viewing and archival. This destination could be on-disk files, or a logging service such as Splunk. Either way, the final destination of the logs can be managed by your execution environment independently of the application.

For example, we could redirect the stdout and stderr streams to on-disk files when starting the application like so:

```
$ go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
```

Note: Using the double arrow >> will append to an existing file, instead of truncating it when starting the application.”

Excerpt From: Alex Edwards. “Let's Go”


There are two main output streams in Linux (and other OSs), standard output (stdout) and standard error (stderr). Error messages, like the ones you show, are printed to standard error. The classic redirection operator (command > file) only redirects standard output, so standard error is still shown on the terminal. To redirect stderr as well, you have a few choices:

Redirect stdout to one file and stderr to another file:
```
$ command > out 2>error
```

Redirect stdout to a file (>out), and then redirect stderr to stdout (2>&1):
```
$ command >out 2>&1
```

Redirect both to a file (this isn't supported by all shells, bash and zsh support it, for example, but sh and ksh do not):
```
$ command &> out
```

### Best practices
* Use Panic or Fatal inside main only and return errors from other places.
* Log output to std streams and redirect the output to a file at runtime.

### Concurrent logging
Custom loggers create with log.New are concurrency-safe, they can be shared by
multiple goroutines.

If you have multiple loggers writing to the same destination you need to make
sure that the destination's underlying Write method is concurrency-safe too.

