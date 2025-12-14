package redis

import (
	"context"
	"log"
	"yus/internal/models"

	"github.com/redis/go-redis/v9"
)

func (r *RedisStore) RateLimiter(ctx context.Context, rateLimit *models.RateLimit) (int, error) {

	rateLimiter := redis.NewScript(`

	-- KEYS[1] = rate limit key (per IP)
	-- ARGV[1] = bucket capacity
	-- ARGV[2] = refill rate (tokens per second)
	-- ARGV[3] = current time (milliseconds)

	--get values from go
	local key = KEYS[1]
	local capacity = tonumber(ARGV[1])
	local refill_rate = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])

	-- Get tokens from redis
	local data = redis.call("HMGET",key,"tokens","timestamp")
	local tokens = tonumber(data[1])
	local last_time = tonumber(data[2])

	-- If bucket does not exist
	if tokens == nil then
		tokens = capacity
		last_time = now
	end

	-- Refill tokens
	local elapsed_time = math.max(0,now - last_time)
	local refill = (elapsed_time/1000) * refill_rate
	tokens = math.min(capacity,math.floor(tokens + refill))

	-- find allowed status
	local allowed = 0
	if tokens >=1 then
		tokens = tokens - 1
		allowed = 1
	end

	-- Save updated state
	redis.call("HMSET",key,"tokens",tokens,"timestamp",now)
	redis.call("EXPIRE",key,60)
	
	return allowed

	`)

	allowed, err := rateLimiter.Run(ctx, r.RedisClient,
		[]string{rateLimit.Key},
		rateLimit.Capacity,
		rateLimit.RefillPerSecond,
		rateLimit.TimeStamp,
	).Int()

	if err != nil {
		log.Println("error while run rate limit script - ", err)
	}

	return allowed, err
}
