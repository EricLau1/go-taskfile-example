package delay

import (
	"log"
	"math/rand"
	"time"
)

func Random(maxDelayInSeconds int) {
	delay := rand.Intn(maxDelayInSeconds)
	sleep := time.Second * time.Duration(delay)
	log.Println("Random Delay:", sleep.String())
	time.Sleep(sleep)
}
