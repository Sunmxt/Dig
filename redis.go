package discover

import (
    "github.com/gomodule/redigo/redis"
)

type RedisRegistry struct {
    redis   *redis.Pool
    prefix  string

    service map[string]*RedisServiceEntry
}

type RedisServiceEntry struct {
    node        map[string]map[string]string
    registry    *RedisRegistry
}

func NewRedisPoolRegistry(pool *redis.Pool, prefix string) (Registry, error) {
    return RedisRegistry{
        redis: pool,
        prefix: prefix,
        service: make(map[string]*RedisServiceEntry, 0),
    }, nil
}

func NewRedisRegistry(network, prefix string, maxIdle, maxActive int) (Registry, error) {
    return NewRedisRegistry(&redis.Pool{
        Dail: func() (redis.Conn, error) {
            return redis.Dail("tcp", network)
        },
        MaxIdle: maxIdle,
        MaxActive: maxActive,
        Wait:       true,
        MaxConnLifetime: 0,
        IdleTimeout: 0,
    }, prefix)
}

func Connect(args ...interface{}) (Registry, error) {
    if len(args) < 2 {
        return nil, ErrInvalidArguments
    }
    switch first := args[0].(type) {
    case *redis.Pool:
        if len(args) != 2 {
            return nil, ErrInvalidArguments
        }
        prefix, ok := args[1].(string)
        if !ok {
            return nil, ErrInvalidArguments
        }
        return NewRedisPoolRegistry(first, prefix)
    case string:
        var maxIdle, maxActive int
        if len(args) != 4 {
            return nil, ErrInvalidArguments
        }
        prefix, ok := args[1].(string)
        if !ok {
            return nil, ErrInvalidArguments
        }
        if maxIdle, ok = args[2].(int); !ok {
            return nil, ErrInvalidArguments
        }
        if maxActive, ok = args[3].(int); !ok {
            return nil, ErrInvalidArguments
        }
        return NewRedisPoolRegistry(first, prefix, maxIdle, maxActive)
    }

    return nil, ErrInvalidArguments
}

func (r *RedisRegistry) Service(name string) (Service, error) {
}

func (r *RedisRegistry) Poll() (bool, error) {
    conn := redis.Get()
    defer conn.Close()

    svcs, err := redis.Strings(conn.Do("SMEMBERS", r.prefix + "{dig-services}"), nil)
    if err != nil {
        return false, err
    }
    changed := false
    for idx := range svcs {
        entry, ok := r.service[svcs[idx]]
        if !ok {
            r.service[svcs[idx]] = nil
            changed = true
        }
        if entry != nil && entry.poll(conn) {
            changed = true
        }
    }
    return changed, false
}

func (r *RedisRegistry) Close() {
}
