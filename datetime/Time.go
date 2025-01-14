package datetime

import "time"

type TimeService struct{}

func NewTimeService() (t *TimeService) {
	return new(TimeService)
}

func Time() (t *TimeService) {
	return NewTimeService()
}

func (t *TimeService) GetCurrentTimeAsSortableString() (currentTime string) {
	return time.Now().Format("20060102_150405")
}
