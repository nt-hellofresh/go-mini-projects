package external

import (
	"fmt"
	"math/rand"
	"time"
)

func GetValueLongRunningTask() int {
	val := rand.Intn(200) + 1500
	dur := time.Duration(val) * time.Millisecond
	time.Sleep(dur)

	fmt.Println("GetValueLongRunningTask ran for", dur)

	return val
}
