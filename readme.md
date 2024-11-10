# Golang microservices assignment

This is a take home test for <>.

## Build and Run

```
GOOS=linux go build -o ports ./cmd
docker build --no-cache -t ports .
docker run -m 100m ports .
```

## Test

```
go test ./... -v
```

A github actions linter and test runner is also provided.
