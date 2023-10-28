# go-book-script

A rewrite of [this Python notebook](https://github.com/yeargin/book-search) using Golang.

## How to use

- [Install Go locally](https://go.dev/doc/install) if you haven't already.
- Clone this repository.
- Export your Goodreads library as a CSV.  Do nothing to clean up the data.  Save it to the `/imports` directory (I named mine `goodreads-export.csv`).
- `go run main.go`.  This may take a minute.
- Your output will be at `goodreads.json`.  It will only include books that have a `Date Read` value.  If the thought of declaring `To Read` amnesty bothers you, you can pass in the `-include-tbr` flag:  `go run main.go -include-tbr`

## Output 

```json
  [
    {
        "title": "",
        "authors": [],
        "date_finished": "",
        "isbn": "", // sourced either from your Goodreads export csv or from the Google Books API
    }
  ]
```

## Possibly useful documentation

- [Google Books API documentation](https://developers.google.com/books/) - No API key is required at this time.
- [OpenLibrary API documentation](https://openlibrary.org/developers/api) - I queried this in an earlier iteration to get a cover image filename, but it turns out you can just use `https://covers.openlibrary.org/b/$key/$value-$size.jpg` where `key` is `isbn`, `value` is the ISBN, and `size` is either `S`, `M`, or `L`.