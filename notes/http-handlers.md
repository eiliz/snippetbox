# The http.Handler interface

```
type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}
```

```
mux.HandleFunc("/snippet/create", createSnippet) will end up as
mux.Handle("/snippet/create", HandlerFunc(createSnippet))
```

We need to convert the regular func we have to a func that implements the Handler interface via HandlerFunc.

```
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

HandlerFunc is just an adapter that wraps over your normal func that already has the right signature. And then implements ServeHTTP which just calls your func.