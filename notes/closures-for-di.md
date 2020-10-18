# Closures for dependency injection

When you have your handlers spread out over multiple packages you can't use the
simple pattern of declaring your handlers as methods against the application
struct (I don't see why at this point => I'm missing something; is it because
you can only declare the handlers against the app struct if they're in the same package?).

A solution for this case is to have a config package exporting the Application
struct and create closures over the app instance created in main.

```
func main() {
    app := &config.Application{
        ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    }

    mux.Handle("/", handlers.Home(app))
}
```

```
func Home(app *config.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ...
        ts, err := template.ParseFiles(files...)
        if err != nil {
            app.ErrorLog.Println(err.Error())
            http.Error(w, "Internal Server Error", 500)
            return
        }
        ...
    }
}
```