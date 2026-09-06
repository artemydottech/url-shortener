# url-shortener

A small url-shortener service. Solution for the [URL Shortening Service](https://roadmap.sh/projects/url-shortening-service) project from roadmap.sh.

Links live in memory, so they disappear when the process restarts.

## Run

```
cp .env.example .env
go run .
```

`PORT` accepts either form — `8080` or `:8080`. Without it the server listens on `:8080`.

## Use

```
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```

Expected response:
`{"short_code":"X7kP9m","short_url":"http://localhost:8080/X7kP9m"}`

`short_url` is built from the host the request came in on, so behind a proxy or on a real domain the link points there and not at localhost.

Visiting `http://localhost:8080/<short_code>` redirects to the original URL. An unknown code returns 404.

## Develop

```
go test -race ./...
gofmt -l .
go vet ./...
```
