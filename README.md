# Verisart API

# Requirements

- Golang 1.11.2
- Docker
- docker-compose
- [Go dep](https://github.com/golang/dep) (Golang dependency management)
- [Tavern](https://github.com/taverntesting/tavern) (API Testing)

# Running the API

Simply run `docker-compose up` and the api will be available at `http://localhost:8000`

# Owner identification

In order to identify the requesting users these endpoints will require an `X-Owner-ID` Header to be set, the endpoints in questions are:

- `/certificate[/:id]` with HTTP Methods `POST, PUT, DELETE`
- `/certificate/:id/transfer` with HTTP Method `POST`

A user account will need to be created at this endpoint:
- `/users/` send a POST request with a JSON body structured as such:
```
{
    "email": "user@email.com",
    "name": "Sample User
}
```

The response on successful creation will be :
```
{
    "id": "abcdef12344434243",
    "email": "user@email.com",
    "name": "Sample User
}
```

You can use the `id` as a `X-Owner-ID`

# Unit Tests

To run the unit tests run `go test ./...`, to add coverage details `go test ./... --cover`

Result should look something like this:

```
# ~/go/src/verisart-api $ go test ./... --cover
?       verisart-api    [no test files]
ok      verisart-api/internal   0.002s  coverage: 100.0% of statements
ok      verisart-api/internal/models    0.002s  coverage: 100.0% of statements
```

# Exploring the API

A set of [Postman](https://www.getpostman.com/) requests have been exported in the `examples/` folder.

**If you change the `email` value when creating a user you will have to update the `X-Owner-ID` values in other tests/examples to reflect the new ID**