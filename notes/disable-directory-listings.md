# Disabling directory listings

There are 2 approaches:

* have an empty index.html file inside each dir you want to disable listing in

```
find ./ui/static -type d -exec touch {}/index.html \;
```

* create a custom implementation of http.FileSystem and have it return an os.ErrNotExist for any dir.