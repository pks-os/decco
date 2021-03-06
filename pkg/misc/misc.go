package misc

import (
	"github.com/cenkalti/backoff"
	"time"
)

func DefaultBackoff() backoff.BackOff {
	return &backoff.ExponentialBackOff{
		InitialInterval: time.Second * 3,
		RandomizationFactor: 0.0,
		Multiplier: 2,
		MaxInterval: time.Second * 60,
		MaxElapsedTime: time.Second * 300,
		Clock: backoff.SystemClock,
	}
}



