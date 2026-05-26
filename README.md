# url-shortener

A small url-shortener service. Solution for the [URL Shortening Service](https://roadmap.sh/projects/url-shortening-service) project from roadmap.sh.

To test it, send a curl request in the following format:

```
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```

Expected response:
`{"short_code":"X7kP9m","short_url":"http://localhost:8080/X7kP9m"}`

Visiting a link of the form `http://localhost:8080/short_code` redirects to the original URL.

A custom port can also be provided via the `.env` file.
