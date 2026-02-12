# url-shortener

Небольшой url-shortener сервис, для тестирования необходимо отправить curl запрос формата:

```
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```
