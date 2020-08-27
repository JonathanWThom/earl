# Earl

Link shortener

# Usage

1. `go run main.go &`
2. `curl -X POST localhost:8080/links -d "url=https://jonathanthom.com"`
3. `curl -L localhost:8080/foobar`

# Endpoints

### Create Account

`curl -X POST https://earlurl.herokuapp.com/accounts`

### Create Link

With account:
`curl -X POST https://earlurl.herokuapp.com/links -d
"url=https://jonathanthom.com" -H "Authorization: basic
token-from-account-creation"`

Without account:
`curl -X POST https://earlurl.herokuapp.com/links -d
"url=https://jonathanthom.com"`

### Inspect Links (only works with account)
`curl https://earlurl.herokuapp.com/links -H "Authorization: basic
token-from-account-creation"`

### Visit Link

`curl -L https://earlurl.herokuapp.com/{identifier}` or just use a browser, of
course!
