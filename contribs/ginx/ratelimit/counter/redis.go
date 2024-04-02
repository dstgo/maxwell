package counter

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

const countLuaScript = `
local count = 0
count = redis.call('get', KEYS[1])
-- exceed threshold
if not count then
    redis.call('set', KEYS[1], 0)
    redis.call('expire', KEYS[1], ARGV[2])
elseif tonumber(count) >= tonumber(ARGV[1]) then
    return 0
end
-- increase count
redis.call('incr', KEYS[1])
return 1
`

func NewRedisCountLimiter(prefix string, limit int, ttl time.Duration, client *redis.Client) *RedisCountLimiter {
	return &RedisCountLimiter{prefix: prefix, limit: limit, ttl: ttl, client: client}
}

// RedisCountLimiter implements Persistent with redis storage
type RedisCountLimiter struct {
	prefix string
	limit  int
	ttl    time.Duration
	client *redis.Client
	keyFn  func(ctx *gin.Context) string
}

func (r *RedisCountLimiter) Allow(ctx *gin.Context) (bool, error) {
	key := r.keyFn(ctx)
	key = r.prefix + ":" + key
	result := r.client.Eval(ctx, countLuaScript, []string{key}, []any{r.limit, int(r.ttl.Seconds())})
	if i, err := result.Int(); err != nil {
		return false, err
	} else {
		return i != 0, nil
	}
}
