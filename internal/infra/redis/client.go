package redis

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"streaming-golang/internal/platform/config"
)

type Client struct {
	client *redis.Client
}

func Open(cfg config.Redis) (*Client, error) {
	if cfg.URL == "" || cfg.URL == "NOT SET" {
		return nil, nil
	}

	opts := &redis.Options{
		Addr:        cfg.URL,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	}

	if cfg.UseSSL {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(opts)

	// Quick ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) IsUserAllowed(ctx context.Context, userID string) (bool, error) {
	if c.client == nil {
		return false, nil // If no redis, we can't check redis
	}

	// Assuming the C# logic stores allowed users in a Redis Set or just key presence.
	// For parity, let's assume it checks if the key exists or is part of a set.
	// Usually, the cache key might be something like "AllowedUsers:{userId}" or checking a hash.
	// We'll use a direct GET or SISMEMBER. Let's assume it's a GET for "AllowedUser:{userId}" returning a boolean/struct,
	// OR it's a set "AllowedUsers" where userID is a member. 
	// For a simple parity, let's check if the key exists:
	key := "AllowedUser:" + strings.ToLower(userID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	
	// If key exists, they are allowed
	return val != "", nil
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
