package external

import (
	"fmt"
	"math/rand"
	"time"
)

func WowSuperLongRunningFunction() int {
	val := rand.Intn(200) + 500
	dur := time.Duration(val) * time.Millisecond
	time.Sleep(dur)

	fmt.Println("WowSuperLongRunningFunction ran for", dur)

	return val
}
