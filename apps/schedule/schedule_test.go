package schedule_test

import (
	"fmt"
	"github.com/getevo/evo-ng/apps/redis"
	"github.com/getevo/evo-ng/apps/redis/distributed"
	"github.com/getevo/evo-ng/apps/schedule"
	"github.com/getevo/evo/lib/text"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestSetInterval(t *testing.T) {
	schedule.Register()
	var job = schedule.SetInterval(func() error {
		fmt.Println("interval reached")
		return nil
	}, 1*time.Second)
	fmt.Println(*job)
	for {
		for _, item := range schedule.Jobs() {
			fmt.Printf("%+v \n", *item)
		}
		time.Sleep(5 * time.Second)
	}

}

func TestSetTimeout(t *testing.T) {
	schedule.Register()
	schedule.SetTimeout(func() error {
		fmt.Println("5s reached")
		return nil
	}, 5*time.Second)
	for {
		fmt.Println(schedule.Jobs())
		time.Sleep(1 * time.Second)
	}
}

func TestScheduler(t *testing.T) {
	redis.Connect(&redis.Config{
		Server: strings.Split("192.168.1.57:6379,192.168.1.58:6380,192.168.1.59:6381", ","),
	})
	schedule.Register()
	var blockingJob = schedule.NewSchedule("myjob", MyJob2).SetRecurring(1 * time.Second).SetBlocking(true)
	blockingJob.Start()

	var lock = distributed.NewMachineLock("test")
	var NonblockingJob = schedule.NewSchedule("myjob2", MyJob).SetRecurring(1 * time.Second).SetBlocking(false).SetLimiter(lock)
	NonblockingJob.Start()

	for {
		for _, item := range schedule.Jobs() {
			fmt.Println(text.ToJSON(item))
		}
		time.Sleep(1 * time.Second)
	}
}

func MyJob() error {

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) // n will be between 0 and 10
	//fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	//fmt.Println("Done Sleeping for ", n)
	return nil
}

func MyJob2() error {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) // n will be between 0 and 10
	fmt.Printf("Sleeping %d seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	//fmt.Println("Done Sleeping for ", n)
	return nil
}
