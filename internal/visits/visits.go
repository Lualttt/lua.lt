package visits

import "sync"

type visits struct {
	sync.RWMutex
	amount int
}

var v = visits{}

func GetVisits() int {
	v.RLock()
	defer v.RUnlock()

	return v.amount
}

func SetVisits(amount int) {
	v.Lock()
	v.amount = amount
	v.Unlock()
}
