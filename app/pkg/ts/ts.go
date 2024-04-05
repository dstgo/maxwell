package ts

import "time"

func Ts() int64 {
	return time.Now().UnixMicro()
}

func ParseTs(ts int64) time.Time {
	return time.UnixMicro(ts)
}
