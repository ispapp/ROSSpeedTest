package speedtest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type SpeedTest struct {
	Actives map[string]Test
	sync.Mutex
}

//	func oldHandler(res http.ResponseWriter, req *http.Request) {
//		// keep old logics for authorisation
//		speed := SpeedTest{}
//		speed.Handler(res, req)
//	}
func (t *SpeedTest) Handler(res http.ResponseWriter, req *http.Request) {
	newtime := time.Now()
	switch req.Method {
	case "GET":
		id := req.URL.Query().Get("id")
		seq := req.URL.Query().Get("seq")
		ping := req.URL.Query().Get("ping")
		log.Println("id", id, "seq", seq, "ping", ping)
		if req.UserAgent() == "MikroTik" && id != "" && seq != "" && ping != "" {
			_ping, err := ConvertToMilliseconds(ping)
			if err != nil {
				res.Write([]byte{})
			}
			tt := Test{
				TX:     []int64{},
				RX:     []int64{},
				PING:   _ping,
				TestID: id,
			}
			if old, ok := t.Actives[id]; ok {
				dsince := time.Since(time.Unix(old.CreatedAt, 0)).Milliseconds() - time.Since(newtime).Milliseconds()
				if i, err := strconv.ParseInt(seq, 10, 32); err == nil {
					if i >= 4 {
						delete(t.Actives, id) // cleanup for less memory usage
					}
				}
				size := getPacketSize(req)
				// timeTx/timeRx=sizeTx/sizeRX
				// timeTx/(timeRxTx - timeTx)=sizeTx/sizeRX
				// timeTx =sizeTx.(timeRxTx - timeTx)/sizeRX
				// timeTx + timeTx/sizeRX=sizeTx.timeRxTx/sizeRX
				// timeTx(sizeRX + 1)/sizeRX=sizeTx.timeRxTx/sizeRX
				// timeTx=(sizeTx.timeRxTx.sizeRX)/(sizeRX.(sizeRX + 1))
				timeTx := dsince / 100
				TX := int64(size) * 8 * 1000 / timeTx
				fmt.Println(timeTx)
				RX := (_ping * 74 * 1000) - int64(TX)
				_tx := append(old.TX, int64(TX))
				_rx := append(old.RX, RX)
				txAvg := Avg(_tx)
				rxAvg := Avg(_rx)
				// send time of sending as created ..
				// and for rx i need time of sending and time of recept and tx rate to calculate recept time
				tt = Test{
					TX:        _tx,
					RX:        _rx,
					TxAvg:     int64(txAvg),
					RxAvg:     int64(rxAvg),
					PING:      _ping,
					TestID:    id,
					Size:      int64(len(fmt.Sprintf(":return %s", ROString(tt))) + size + 8),
					CreatedAt: time.Now().Unix(),
				}
				cmd := fmt.Sprintf(":return %s", ROString(tt))
				res.WriteHeader(http.StatusOK)
				io.WriteString(res, cmd)
			}
			if t.TryLock() {
				tt.CreatedAt = time.Now().Unix()
				t.Actives[id] = tt
				t.Unlock()
			}
			res.Write([]byte{})
		}
		res.Write([]byte{})
	}
}
