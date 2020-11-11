# ServeMux

### The DefaultServeMux

The reason why you can skip creating a mux and register routes with just http.Handle() or http.HandleFunc() is because Go creates a global var where it stores a default mux instance.

```
var DefaultServeMux = NewServeMux()
```

It's not a good idea to rely on this feature in prod envs since a global mux would be a security liability.



### Types of paths

Go's servemux supports 2 types of paths: fixed paths and subtree paths.

Fixed paths don't have the final slash and are the ones you'd like to be exactly matched: /snippet will match requests from /snippet but not snippet/create.

Subtree paths end in a slash and will match whenever the URL request start matches it.


### Features

* longer URL patterns take precedence over shorter ones, no matter the order they're in
* request URL paths are automatically sanitized (“For example, if a user makes a request to /foo/bar/..//baz they will automatically be sent a 301 Permanent Redirect to /foo/baz instead.” - Excerpt From: Alex Edwards. “Let's Go” )
* If a subtree path has been registered and a request is received for that subtree path without a trailing slash, then the user will automatically be sent a 301 Permanent Redirect to the subtree path with the slash added. For example, if you have registered the subtree path /foo/, then any request to /foo will be redirected to /foo/


### Headers

Go only allows you to call WriteHeader once per response
The first call to Write will call WriteHeader with a 200 OK automatically before writing any data, if WriteHeader has not been called yet.

So WriteHeader is usually present when you want to send errors and must come before Write.

Header gives you access to the response header map to add new headers and writing to it must always come before calls to WriteHeader and Write or else the headers added will be ignored.

```
w.Header().Set("Allow", http.MethodPost)
```

WriteHeader will actually send the client the status code header.
```
w.WriteHeader(405)
```
Then the client gets his body data as well.
w.Write([]byte("Method not allowed"))


Go will sniff the content type without setting the content-type header.
But in the case of JSON we need to explicitly do it as the sniffer
mistakenly sets it to text/plain.

The header name is auto canonicalized by Go.
For HTTP2 conns, Go sets header names and values to lowercase as per the HTTP2 spec.
```
w.Header().Set("content-type", "application/json")
```

Set the header directly into the header map to escape the auto canon. behavior.
```
w.Header()["X-XSS-Protection"] = []string{"1;mode=block"}
```

Go auto sets the content-type, content-length and date headers.
You cannot remove them with Del() but you can directly set them to nil inside http.Header().
```
w.Header()["Date"] = nil
```