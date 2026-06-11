package datetime

import "time"

func GetCurrentTimeAsSortableString() (currentTime string) {
	return time.Now().Format("20060102_150405")
}
