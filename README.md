# Earl

Link shortener and view logger

# Endpoints

### POST /accounts

Create an account, receive a token back

Example request:
```
curl -X POST http://localhost:8000/accounts`
```

Example response:
```
{"token":"WS1kVlNQRk5fV003Nzh2MFozSUl3"}
```

### POST /links

Create links associated with an account, or not.

Example request, with account:
```
curl -X POST http://localhost:8000/links \
-d "url=https://jonathanthom.com" \
-H "Authorization: basic token-from-account-creation"
```

Example request, without account:
```
curl -X POST http://localhost:8000/links \
-d "url=https://jonathanthom.com" 
```

Example response:
```
{"original":"https://jonathanthom.com","shortened":"http://localhost:8000/orxHsI","views":null}
```

### GET /links

Inspect Links & Views (only works with account)

Example request:
```
curl http://localhost:8000/links \
-H "Authorization: basic token-from-account-creation"
```

Example response:
```
[{"original":"https://jonathanthom.com","shortened":"http://localhost:8000/R-KMIa","views":[{"createdAt":"2020-08-29T21:42:56.706419Z","referer":"","country":"United States of America","city":"Bellingham","zipCode":"98225"}]}]
```

### Visit Link

`curl -L http://localhost:8000/{identifier}` or just use a browser, of
course!
