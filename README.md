## what it is

It is an simple application that reads data, creates an index and allows the user to search products. I built with a very basic understanding of how an search engine indexing would work, I ended up researching reading a few things about search engines, indexing strategies and fuzzy string searching etc. And I'm planning to build a more ambitious project on this domain.

## depends on

* [go-redis](https://github.com/redis/go-redis)
* redis
* docker

### how to run

```bash
# setup project
go mod tidy

# start docker container with redis
docker compose up -d

# extract data from sample file and creates index
go run cmd/main.go

```

### query examples
```bash
# by price, returns all products in the price range
Enter query: price:100

# by name, returns all products with the word in the name
Enter query: name:Ball

# by name, returns all products with the word in the description
Enter query: description:Sport

# by id, returns only the searched product
Enter query: id:3e2ff818-7512-4025-a026-07a1c7583cb1
```