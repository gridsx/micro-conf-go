package mc

type ChangeListener interface {

	// OnChange data will be sent via this method
	OnChange(data []byte) error

	// RetryTimes   when getting error from onChange call
	RetryTimes() int
}
