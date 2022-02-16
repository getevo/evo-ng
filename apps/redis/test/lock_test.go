package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/getevo/evo-ng/apps/redis"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	Init()

	// try to obtain again
	lock, err := redis.Locker.Obtain("key1", 10*time.Second, nil, nil)
	if exp, got := redis.ErrLockNotObtained, err; !errors.Is(got, exp) {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	fmt.Println(lock.Key())
	// manually unlock
	if err := lock.Release(ctx); err != nil {
		t.Fatal(err)
	}

	// lock again
	lock, err = redis.Locker.Obtain("key1", 10*time.Second, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Release(ctx)
}

/*func TestObtain(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	lock := quickObtain(t, rc, time.Hour)
	if err := lock.Release(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestObtain_metadata(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	meta := "my-data"
	lock, err := Obtain(ctx, rc, lockKey, time.Hour, &redis.Options{Metadata: meta})
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Release(ctx)

	if exp, got := meta, lock.Metadata(); exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestObtain_retry_success(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	// obtain for 20ms
	lock1 := quickObtain(t, rc, 20*time.Millisecond)
	defer lock1.Release(ctx)

	// lock again with linar retry - 3x for 20ms
	lock2, err := Obtain(ctx, rc, lockKey, time.Hour, &redis.Options{
		RetryStrategy: redis.LimitRetry(redis.LinearBackoff(20*time.Millisecond), 3),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer lock2.Release(ctx)
}

func TestObtain_retry_failure(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	// obtain for 50ms
	lock1 := quickObtain(t, rc, 50*time.Millisecond)
	defer lock1.Release(ctx)

	// lock again with linar retry - 2x for 5ms
	_, err := Obtain(ctx, rc, lockKey, time.Hour, &redis.Options{
		RetryStrategy: redis.LimitRetry(redis.LinearBackoff(5*time.Millisecond), 2),
	})
	if exp, got := ErrNotObtained, err; !errors.Is(got, exp) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestObtain_concurrent(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	numLocks := int32(0)
	numThreads := 100
	wg := new(sync.WaitGroup)
	errs := make(chan error, numThreads)
	for i := 0; i < numThreads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			wait := rand.Int63n(int64(10 * time.Millisecond))
			time.Sleep(time.Duration(wait))

			_, err := Obtain(ctx, rc, lockKey, time.Minute, nil)
			if err == ErrNotObtained {
				return
			} else if err != nil {
				errs <- err
			} else {
				atomic.AddInt32(&numLocks, 1)
			}
		}()
	}
	wg.Wait()

	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	if exp, got := 1, int(numLocks); exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLock_Refresh(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	lock := quickObtain(t, rc, time.Hour)
	defer lock.Release(ctx)

	// check TTL
	assertTTL(t, lock, time.Hour)

	// update TTL
	if err := lock.Refresh(ctx, time.Minute, nil); err != nil {
		t.Fatal(err)
	}

	// check TTL again
	assertTTL(t, lock, time.Minute)
}

func TestLock_Refresh_expired(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	lock := quickObtain(t, rc, 5*time.Millisecond)
	defer lock.Release(ctx)

	// try releasing
	time.Sleep(10 * time.Millisecond)
	if exp, got := ErrNotObtained, lock.Refresh(ctx, time.Minute, nil); !errors.Is(got, exp) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLock_Release_expired(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	lock := quickObtain(t, rc, 5*time.Millisecond)
	defer lock.Release(ctx)

	// try releasing
	time.Sleep(10 * time.Millisecond)
	if exp, got := redis.ErrLockNotHeld, lock.Release(ctx); !errors.Is(got, exp) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func TestLock_Release_not_own(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(redisOpts)
	defer teardown(t, rc)

	lock := quickObtain(t, rc, time.Hour)
	defer lock.Release(ctx)

	// manually transfer ownership
	if err := rc.Set(ctx, lockKey, "ABCD", 0).Err(); err != nil {
		t.Fatal(err)
	}

	// try releasing
	if exp, got := redis.ErrLockNotHeld, lock.Release(ctx); !errors.Is(got, exp) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}

func quickObtain(t *testing.T, rc *redis.Client, ttl time.Duration) *redis.Lock {
	t.Helper()

	lock, err := Obtain(context.Background(), rc, lockKey, ttl, nil)
	if err != nil {
		t.Fatal(err)
	}
	return lock
}

func assertTTL(t *testing.T, lock *redis.Lock, exp time.Duration) {
	t.Helper()

	ttl, err := lock.TTL(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	delta := ttl - exp
	if delta < 0 {
		delta = 1 - delta
	}
	if delta > time.Second {
		t.Fatalf("expected ~%v, got %v", exp, ttl)
	}
}

func teardown(t *testing.T, rc *redis.Client) {
	t.Helper()

	if err := rc.Del(context.Background(), lockKey).Err(); err != nil {
		t.Fatal(err)
	}
	if err := rc.Close(); err != nil {
		t.Fatal(err)
	}
}
*/
