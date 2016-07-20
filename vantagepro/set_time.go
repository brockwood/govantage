package vantagepro

import (
	"errors"
	"time"
)

func (vp *VPConsole) SetTime() error {
	vp.wakeUpConsole()
	cmderr := vp.sendCommand("SETTIME")
	if cmderr != nil {
		return cmderr
	}
	timePacket := make([]byte, 6)
	now := time.Now()
	timePacket[0] = byte(now.Second())
	timePacket[1] = byte(now.Minute())
	timePacket[2] = byte(now.Hour())
	timePacket[3] = byte(now.Day())
	timePacket[4] = byte(now.Month())
	timePacket[5] = byte(now.Year() - 1900)
	crcPacket := getCRC(timePacket)
	timePacket = append(timePacket, crcPacket...)
	_, err := vp.portconnection.Write(timePacket)
	if err != nil {
		return err
	}
	getAck := make([]byte, 1)
	ackSize, ackError := vp.portconnection.Read(getAck)
	if ackSize != 1 || getAck[0] != ACK || ackError != nil {
		return errors.New("Did not get ACK back from time packet")
	}
	return nil
}
