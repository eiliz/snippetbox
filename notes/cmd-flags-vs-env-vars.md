# Managing configuration

The idiomatic way to manage config vars is with command line flags. While you can use env vars, there's advantages in using the cli flags:

1. cli flags can have default values, env vars cannot
2. cli flags get automatic type conversions (flag,.String(), flag.Int())
3. they are documented when starting the app with the -help flag

The best of both worlds is to setup env vars and then feed them into the cli flags:
```
$ export SNIPPETBOX_ADDR=":9999"
$ go run ./cmd/web -addr=$SNIPPETBOX_ADDR
```

### Bool flags
These need to be explicitly set to false or it's interpreted as being true:
```
$ go run example.go -flag=true
$ go run example.go -flag
```

### Storing config settings in a single struct

```
type Config struct {
    Addr      string
    StaticDir string
}

...

cfg := new(Config)
flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
flag.Parse()
```