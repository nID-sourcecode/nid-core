package redisclient

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	client *Client
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

func (s *RedisTestSuite) SetupTest() {
	// Init Redis client. Mini-redis does not support expiration, so expiry must be tested
	// in integration tests.
	mr, err := miniredis.Run()
	s.Require().NoError(err)

	options := &redis.Options{
		Addr:     mr.Addr(),
		Password: "",
		DB:       0,
	}
	redisClient, err := NewRedisClient(options)
	s.Require().NoError(err)

	// Connecting to faulty address should not produce a pingable client.
	_, err = NewRedisClient(&redis.Options{
		Addr: "faultyaddr:1234",
	})
	s.Require().ErrorIs(err, ErrRedisConnection)

	s.client = redisClient
}

func (s *RedisTestSuite) TestRedisSetGetString() {
	key := "keystring"
	value := "valuestring"

	ctx := context.Background()

	// Mini-redis does not support expiration, so this functionality must be tested in integration tests.
	err := s.client.SetString(ctx, key, value, 2*time.Second)
	s.Require().NoError(err)

	returnValue, err := s.client.GetString(ctx, key)

	s.Require().NoError(err)
	s.Require().Equal(value, returnValue)

	wrongKey, err := s.client.GetString(ctx, "wrongkey")
	s.Require().ErrorIs(err, ErrKeyDoesNotExist)
	s.Empty(wrongKey)
}

func (s *RedisTestSuite) TestRedisSetGetInteger() {
	key := "one-hundred"
	value := 100

	ctx := context.Background()

	// Mini-redis does not support expiration, so this functionality must be tested in integration tests.
	err := s.client.SetInteger(ctx, key, value, 4*time.Second)
	s.Require().NoError(err)

	returnValue, err := s.client.GetInteger(ctx, key)
	s.Require().NoError(err)

	s.Require().Equal(value, returnValue)

	wrongKey, err := s.client.GetInteger(ctx, "wrongkey")
	s.Require().ErrorIs(err, ErrKeyDoesNotExist)
	s.Empty(wrongKey)
}

func (s *RedisTestSuite) TestRedisSetGetStruct() {
	key := "one-hundred"
	type testStruct struct {
		Field1 int    `json:"field1"`
		Field2 string `json:"field2"`
	}
	value := testStruct{
		Field1: 1,
		Field2: "wowSuchAString",
	}
	var returnValue testStruct
	var returnValueWrongKey testStruct

	ctx := context.Background()

	// Mini-redis does not support expiration, so this functionality must be tested in integration tests.
	err := s.client.SetStruct(ctx, key, value, 10*time.Minute)
	s.Require().NoError(err)

	err = s.client.GetStruct(ctx, key, &returnValue)
	s.Require().NoError(err)

	s.Require().Equal(value, returnValue)

	err = s.client.GetStruct(ctx, "wrongkey", &returnValueWrongKey)
	s.Require().ErrorIs(err, ErrKeyDoesNotExist)
	s.Empty(returnValueWrongKey)
}

func (s *RedisTestSuite) TestPubSub() {
	// Init simple channel and messages for testing.
	const testChannel string = "channel1"
	messages := []string{
		"message1",
		"message2",
		"message3",
	}

	// Extra channel and messages to confirm messages are not mixed across channels.
	const testOtherChannel string = "channel2"
	otherMessages := []string{
		"message4",
		"message5",
		"message6",
	}
	ctx := context.Background()

	// Subscribe to a channel.
	pubsub := s.client.Subscribe(ctx, testChannel)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	s.Require().NoError(err)

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Publish a message.
	for _, msg := range messages {
		err = s.client.Publish(ctx, testChannel, msg).Err()
		s.Require().NoError(err)
	}

	// Publish other messages.
	for _, msg := range otherMessages {
		err = s.client.Publish(ctx, testOtherChannel, msg).Err()
		s.Require().NoError(err)
	}

	time.AfterFunc(time.Second, func() {
		// When pubsub is closed channel is closed too.
		err = pubsub.Close()
		s.Require().NoError(err)
	})

	// Consume messages and confirm they are what was sent to only this channel.
	var receivedMessages []string
	for msg := range ch {
		receivedMessages = append(receivedMessages, msg.Payload)
	}
	s.Equal(messages, receivedMessages)
}
