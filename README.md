## indexes generator

This project is an attempt to create an search engine, I decide to start by creating an indexes generation, this way I will have a very basic searching mechanism. Right now it is a work in progress, but is possible to generate indexes from a fake data set.

### how to use

```bash
# setup project
go mod tidy

# gerate data sample

go run seed/seed.go > sample

# generate indexes

go run main.go
```
