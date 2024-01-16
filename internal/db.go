package internal

import (
	"encoding/json"

	"github.com/go-redis/redis"
)

type DB struct {
	redisClient *redis.Client
}

func NewDB(redisAddr string) *DB {
	return &DB{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (db *DB) Close() {
	db.redisClient.Close()
}

func (db *DB) SaveProduct(p Product) (bool, error) {
	serializedProduct, err := serialize(p)
	if err != nil {
		return false, err
	}

	err = db.redisClient.Set(p.ID, serializedProduct, 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *DB) GetProduct(productID string) (Product, error) {
	val, err := db.redisClient.Get(productID).Result()
	if err != nil {
		return Product{}, err
	}

	p, err := deserialize(val)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func serialize(p Product) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func deserialize(s string) (Product, error) {
	p := Product{}
	byt := []byte(s)

	err := json.Unmarshal(byt, &p)
	if err != nil {
		return Product{}, err
	}
	return p, err
}
