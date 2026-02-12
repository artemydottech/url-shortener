# url-shortener

Небольшой url-shortener сервис, для тестирования необходимо отправить curl запрос формата:

```
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```

Ожидаемый ответ:
`{"short_code":"X7kP9m","short_url":"http://localhost:8080/X7kP9m"}`

При переходе на ссылку формата `http://localhost:8080/short_code` ожидается редирект на переданную изначально ссылку
