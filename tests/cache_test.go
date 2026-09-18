package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"grip/config"
	"grip/dto"
	"grip/repository"
	services "grip/service"
	"grip/utils"

	"github.com/redis/go-redis/v9"
)

// requireTestRedis points config.RDB at an isolated Redis database (15)
// and skips when Redis is unreachable. State is flushed before and after.
func requireTestRedis(t *testing.T) {
	t.Helper()
	addr := testRedisAddr(t)

	client := redis.NewClient(&redis.Options{Addr: addr, DB: 15})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		t.Skipf("redis not available at %s (%v), skipping", addr, err)
	}

	config.RDB = client
	flushRedis(t, client)
	t.Cleanup(func() {
		flushRedis(t, client)
		_ = client.Close()
		config.RDB = nil
	})
}

func testRedisAddr(t *testing.T) string {
	t.Helper()
	if v := os.Getenv("TEST_REDIS_ADDR"); v != "" {
		return v
	}
	return "localhost:6379"
}

func flushRedis(t *testing.T, client *redis.Client) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("flush test redis: %v", err)
	}
}

func TestCacheSetGetRoundTrip(t *testing.T) {
	requireTestRedis(t)
	ctx := context.Background()

	type profile struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	if err := utils.SetJSON(ctx, "grip:test:roundtrip", profile{Name: "Kai", Age: 30}, time.Minute); err != nil {
		t.Fatalf("set: %v", err)
	}
	var got profile
	hit, err := utils.GetJSON(ctx, "grip:test:roundtrip", &got)
	if err != nil || !hit {
		t.Fatalf("get hit=%v err=%v", hit, err)
	}
	if got.Name != "Kai" || got.Age != 30 {
		t.Fatalf("got %+v", got)
	}

	hit, err = utils.GetJSON(ctx, "grip:test:missing", &got)
	if err != nil || hit {
		t.Fatalf("missing hit=%v err=%v, want miss", hit, err)
	}
}

func TestCacheDelAndExpire(t *testing.T) {
	requireTestRedis(t)
	ctx := context.Background()

	if err := utils.SetJSON(ctx, "grip:test:del", map[string]string{"a": "b"}, time.Minute); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := utils.Del(ctx, "grip:test:del"); err != nil {
		t.Fatalf("del: %v", err)
	}
	var dst map[string]string
	if hit, _ := utils.GetJSON(ctx, "grip:test:del", &dst); hit {
		t.Fatal("expected miss after del")
	}

	if err := utils.SetJSON(ctx, "grip:test:exp", "v", 0); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := utils.Expire(ctx, "grip:test:exp", time.Second); err != nil {
		t.Fatalf("expire: %v", err)
	}
	time.Sleep(1200 * time.Millisecond)
	var s string
	if hit, _ := utils.GetJSON(ctx, "grip:test:exp", &s); hit {
		t.Fatal("expected key to expire")
	}
}

func TestCacheIncr(t *testing.T) {
	requireTestRedis(t)
	ctx := context.Background()

	n, err := utils.Incr(ctx, "grip:test:counter", time.Minute)
	if err != nil || n != 1 {
		t.Fatalf("incr = %d, %v; want 1", n, err)
	}
	n, err = utils.Incr(ctx, "grip:test:counter", time.Minute)
	if err != nil || n != 2 {
		t.Fatalf("incr = %d, %v; want 2", n, err)
	}
}

func TestCacheDelPattern(t *testing.T) {
	requireTestRedis(t)
	ctx := context.Background()

	for _, k := range []string{"grip:test:pat:a", "grip:test:pat:b", "grip:test:other"} {
		if err := utils.SetJSON(ctx, k, "v", time.Minute); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	if err := utils.DelPattern(ctx, "grip:test:pat:*"); err != nil {
		t.Fatalf("del pattern: %v", err)
	}
	var s string
	if hit, _ := utils.GetJSON(ctx, "grip:test:pat:a", &s); hit {
		t.Fatal("expected pattern keys deleted")
	}
	if hit, _ := utils.GetJSON(ctx, "grip:test:other", &s); !hit {
		t.Fatal("expected non-matching key to survive")
	}
}

func TestCacheDisabledIsSafe(t *testing.T) {
	old := config.RDB
	config.RDB = nil
	defer func() { config.RDB = old }()
	ctx := context.Background()

	if utils.CacheEnabled() {
		t.Fatal("expected CacheEnabled=false with nil client")
	}
	var dst map[string]string
	if hit, err := utils.GetJSON(ctx, "x", &dst); hit || err != nil {
		t.Fatalf("disabled get hit=%v err=%v", hit, err)
	}
	if err := utils.SetJSON(ctx, "x", "v", 0); err == nil {
		t.Fatal("expected disabled set to error")
	}
	if err := utils.Del(ctx, "x"); err != nil {
		t.Fatalf("disabled del: %v", err)
	}
	if err := utils.Expire(ctx, "x", time.Minute); err != nil {
		t.Fatalf("disabled expire: %v", err)
	}
	if _, err := utils.Incr(ctx, "x", time.Minute); err == nil {
		t.Fatal("expected disabled incr to error")
	}
	if err := utils.DelPattern(ctx, "x*"); err != nil {
		t.Fatalf("disabled del pattern: %v", err)
	}
}

func TestUserCacheHitAndInvalidate(t *testing.T) {
	requireTestDB(t)
	requireTestRedis(t)
	ctx := context.Background()

	created, err := services.CreateUser(dto.CreateUserRequest{
		Name: "Cache User", Email: "cacheuser@example.com", Age: 30,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// First read populates the cache.
	if _, err := services.GetUserByID(created.ID); err != nil {
		t.Fatalf("get: %v", err)
	}
	if exists, _ := config.RDB.Exists(ctx, utils.UserKey(created.ID)).Result(); exists != 1 {
		t.Fatal("expected cache key to exist after read")
	}

	// Update invalidates: the next read reflects the DB, never stale cache.
	if _, err := services.UpdateUser(created.ID, dto.CreateUserRequest{
		Name: "Cache Renamed", Email: "cacheuser@example.com", Age: 31,
	}); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, err := services.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if got.Name != "Cache Renamed" {
		t.Fatalf("name = %q, want updated value (stale cache?)", got.Name)
	}

	// Delete invalidates too.
	if err := services.DeleteUser(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if exists, _ := config.RDB.Exists(ctx, utils.UserKey(created.ID)).Result(); exists != 0 {
		t.Fatal("expected cache key gone after delete")
	}
	if _, err := repository.GetUserByID(created.ID); err == nil {
		t.Fatal("expected user gone from db")
	}
}
