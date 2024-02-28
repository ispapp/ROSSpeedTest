package speedtest

type Test struct {
	TestID    string
	TX        []int64
	RX        []int64
	TxAvg     int64
	RxAvg     int64
	PING      int64
	Size      int64
	CreatedAt int64
}
