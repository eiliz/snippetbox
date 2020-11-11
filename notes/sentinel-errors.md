# Sentinel errors

Sentinel errors are error objects stored in a global variable. They are created
using the errors.New function.

Some examples from the std library are io.ErrUnexpectedEOF and
bytes.ErrTooLarge.

The new idiomatic way of checking if an error matches a sentinel error (Go 1.13
and newer) is using errors.Is(). Before you'd just compare the values:

```
if err == sql.ErrNoRows {
    // Do something
} else {
    // Do something else
}
```

However that won't work for wrapped errors. Wrapped errors were also introduced
in Go 1.13 and allow you the ability to add extra information. Since wrapped
errors are not the same object as the original error, comparing them won't work.
However errors.Is will unwrap errors if necessary and then compare them.

```
if errors.Is(err, sql.ErrNoRows) {
    // Do something
} else {
    // Do something else
}
```

There's another function, errors.As() that allows you to check if a potentially
wrapped error has a specific type.
