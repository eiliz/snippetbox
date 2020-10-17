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