// Package redisclient provides implementation of a go redis client.
package redisclient

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

var (
	// ErrKeyDoesNotExist is thrown if the Redis Get says the key does not exist.
	ErrKeyDoesNotExist = errors.New("key does not exist")
	// ErrRedisGet is thrown if there is a different error in getting than the existence of the key.
	ErrRedisGet = errors.New("error in getting value")
	// ErrRedisSet is thrown if setting the value did not work. In practice this should never occur,
	// only perhaps when the process is out of memory.
	ErrRedisSet = errors.New("error in setting value")
	// ErrRedisConnection is thrown if the connection to redis with a ping request fails.
	ErrRedisConnection = errors.New("redis connection failed")
)

// RedisClient is an interface.
type RedisClient interface {
	Set(string, interface{}) error
	Get(string) (string, error)
	Subscribe(string) *redis.PubSub
	Publish(string, string) string
}

// Client contains the connection to the redis server.
type Client struct {
	redisClient *redis.Client
	expiration  time.Duration
}

// NewRedisClient can be called from a service to return a reusable redis connection.
func NewRedisClient(options *redis.Options) (*Client, error) {
	rc := &Client{
		redisClient: redis.NewClient(options),
		expiration:  0,
	}

	// Ping to check connection.
	_, err := rc.redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.CombineErrors(ErrRedisConnection, err)
	}

	return rc, nil
}

// setValue implements the Set method, taking a key and storing the value in the redis server.
func (rc *Client) setValue(ctx context.Context, key string, value interface{}, expireAfter time.Duration) error {
	err := rc.redisClient.Set(ctx, key, value, expireAfter).Err()
	if err != nil {
		return errors.CombineErrors(ErrRedisSet, err)
	}

	return nil
}

// SetTTL sets or resets the expiration of a given key.
func (rc *Client) SetTTL(ctx context.Context, key string, expireAfter time.Duration) error {
	err := rc.redisClient.Expire(ctx, key, expireAfter).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetTTL returns the expiration time of a given key.
func (rc *Client) GetTTL(ctx context.Context, key string) time.Duration {
	duration := rc.redisClient.TTL(ctx, key).Val()
	return duration
}

// SetString enforces a string value input. Underneath the Set method is called.
func (rc *Client) SetString(ctx context.Context, key, value string, expireAfter time.Duration) error {
	return rc.setValue(ctx, key, value, expireAfter)
}

// SetInteger enforces an integer value input. Underneath the Set method is called.
func (rc *Client) SetInteger(ctx context.Context, key string, value int, expireAfter time.Duration) error {
	return rc.setValue(ctx, key, value, expireAfter)
}

// SetStruct is meant for a struct value input. Underneath the Set method is called.
func (rc *Client) SetStruct(ctx context.Context, key string, value interface{}, expireAfter time.Duration) error {
	valueMarshalled, err := json.Marshal(value)
	if err != nil {
		return errors.CombineErrors(ErrRedisSet, err)
	}

	return rc.setValue(ctx, key, valueMarshalled, expireAfter)
}

// getValue implements the Get method, taking a key and calling and returning the value from the redis server.
func (rc *Client) getValue(ctx context.Context, key string) (value string, err error) {
	value, err = rc.redisClient.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", errors.CombineErrors(ErrKeyDoesNotExist, err)
	} else if err != nil {
		return "", errors.CombineErrors(ErrRedisGet, err)
	}

	return value, nil
}

// GetString calls the Get method, but parses the result as an string.
func (rc *Client) GetString(ctx context.Context, key string) (value string, err error) {
	value, err = rc.getValue(ctx, key)

	if err != nil {
		return "", err
	}

	return value, nil
}

// GetInteger calls the Get method, but parses the result as an integer.
func (rc *Client) GetInteger(ctx context.Context, key string) (value int, err error) {
	valueString, err := rc.getValue(ctx, key)
	if err != nil {
		return 0, err
	}

	value, err = strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// GetStruct calls the Get method, but unmarshals the result to the third argument, a struct.
func (rc *Client) GetStruct(ctx context.Context, key string, myStruct interface{}) (err error) {
	valueString, err := rc.getValue(ctx, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(valueString), myStruct)
	if err != nil {
		return err
	}

	return nil
}

// Subscribe implements the Subscribe method.
func (rc *Client) Subscribe(ctx context.Context, channel string) (sub *redis.PubSub) {
	return rc.redisClient.Subscribe(ctx, channel)
}

// Publish implements the Publish method.
func (rc *Client) Publish(ctx context.Context, channel, message string) (value *redis.IntCmd) {
	return rc.redisClient.Publish(ctx, channel, message)
}
