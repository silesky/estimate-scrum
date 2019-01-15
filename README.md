# DOCS

## Install dependencies:
```
go get -d ./...
```

## Start
```
go install estimate/server
estimate
```

## Develop
```
go get github.com/oxequa/realize
realize start
docker-compose build --no-cache
docker-compose up
```


# Test
- you can only test on a package
```
go test estimate/router
```
