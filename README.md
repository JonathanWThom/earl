# Earl

Link shortener and view logger

# Endpoints

### POST /accounts

Create an account, receive a token back

Example request:
```
curl -X POST https://earlurl.herokuapp.com/accounts`
```

Example response:
```
{"token":"WS1kVlNQRk5fV003Nzh2MFozSUl3"}
```

### POST /links

Create links associated with an account, or not.

Example request, with account:
```
curl -X POST https://earlurl.herokuapp.com/links \
-d "url=https://jonathanthom.com" \
-H "Authorization: basic token-from-account-creation"
```

Example request, without account:
```
curl -X POST https://earlurl.herokuapp.com/links \
-d "url=https://jonathanthom.com" 
```

Example response:
```
{"original":"https://jonathanthom.com","shortened":"https://earlurl.herokuapp.com/orxHsI","views":null}
```

### GET /links

Inspect Links & Views (only works with account)

Example request:
```
curl https://earlurl.herokuapp.com/links \
-H "Authorization: basic token-from-account-creation"
```

Example response:
```
[{"original":"https://jonathanthom.com","shortened":"https://earlurl.herokuapp.com/R-KMIa","views":[{"createdAt":"2020-08-29T21:42:56.706419Z","referer":"","country":"United States of America","city":"Bellingham","zipCode":"98225"}]}]
```

### Visit Link

`curl -L https://earlurl.herokuapp.com/{identifier}` or just use a browser, of
course!
