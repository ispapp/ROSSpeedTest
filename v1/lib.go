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

func (t *SpeedTest) Handler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		id := req.URL.Query().Get("id")
		seq := req.URL.Query().Get("seq")
		if req.UserAgent() == "MikroTik" && id != "" && seq != "" {
			// for name, values := range req.Header {
			// 	for _, value := range values {
			// 		fmt.Println(name, value)
			// 	}
			// } bytes / s
			since := time.Now().Unix()
			tt := Test{
				TX:        0,
				RX:        0,
				PING:      0,
				TestID:    id,
				Size:      0,
				CreatedAt: since,
			}
			if old, ok := t.Actives[id]; ok {
				dsince := time.Since(time.Unix(old.CreatedAt, 0)).Milliseconds()
				if i, err := strconv.ParseInt(seq, 10, 32); err == nil {
					if i >= 4 {
						delete(t.Actives, id)
					}
				}
				size := getRequestSize(req)
				TX := ((size * 8) / int(dsince/2)) * 1000
				RX := (((size + old.Size) * 8) / int(dsince/2)) * 1000
				PING := (74 * 1000 * 8) / (TX + RX) // ping packet frame size 74 bytes. which would be 20 bytes of IP header, 8 bytes of ICMP header + 32 data
				// send time of sending as created ..
				// and for rx i need time of sending and time of recept and tx rate to calculate recept time
				tt = Test{
					TX:        TX,
					RX:        RX,
					PING:      int(PING),
					TestID:    id,
					Size:      len(fmt.Sprintf(":return %s", ROString(old))),
					CreatedAt: since,
				}
				log.Default().Printf("seq: %s \t id: %s, reqsize: %dbytes", seq, id, size)
				cmd := fmt.Sprintf(":return %s", ROString(tt))
				res.WriteHeader(http.StatusOK)
				io.WriteString(res, cmd)
			}
			if t.TryLock() {
				t.Actives[id] = tt
				t.Unlock()
			}
			res.Write([]byte{})
		}
		res.Write([]byte{})
	}
}
