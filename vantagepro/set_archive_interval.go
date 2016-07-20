package vantagepro

import (
	"fmt"
	"log"
)

type Interval int

const (
	ONEMINUTE      Interval = 1
	FIVEMINUTES    Interval = 5
	TENMINUTES     Interval = 10
	FIFTEENMINUTES Interval = 15
	THIRTYMINUTES  Interval = 30
	ONEHOUR        Interval = 60
	TWOHOURS       Interval = 120
)

func (vp *VPConsole) setArchiveInterval(interval Interval) error {
	vp.wakeUpConsole()
	getAck := make([]byte, 100)
	fullMessage := []byte{}
	_, err := vp.portconnection.Write([]byte(fmt.Sprintf("SETPER %d\n", interval)))
	if err != nil {
		log.Fatal(err)
	}
	readData := true
	for readData {
		ackSize, _ := vp.portconnection.Read(getAck)
		if ackSize == 0 {
			readData = false
			continue
		}
		fullMessage = append(fullMessage, getAck[0:ackSize]...)
	}
	return nil
}
