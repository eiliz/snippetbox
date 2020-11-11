# go.sum

This file contains the crypto checksums representing the content of the required
packages in go.mod.

It serves 2 purposes:

1. you can run the go mod verify command to check that the checksums of the
   downloaded packages on your machine match the entries in go.sum, so you can
   be confident that they haven't been altered

2. if someone else needs to download all the deps for the project (go mod
   download), they will get an error if there is any mismatch between the deps
   they are downloading and the checksums in the file


## Upgrade package
Run go get with the u flag to upgrade to the latest minor or patch release.

```
$ go get -u github.com/foo/bar
```

Or upgrade to a specific version with:

```
$ go get -u github.com/foo/bar@v2.0.0
```


## Remove unused packages

Run go get with the postfix @none.

```
$ go get github.com/foo/bar@none
```

Or if all refs to the package have been removed run:
```
$ go mod tidy -v
```

to automatically remove all unused packages from go.mod and go.sum.