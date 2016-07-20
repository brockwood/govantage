package vantagepro

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

func (vp *VPConsole) GetLoopPacket() (*loopPacketRaw, error) {
	vp.wakeUpConsole()
	getAck := make([]byte, 1)
	_, err := vp.portconnection.Write([]byte("LOOP 1\n"))
	if err != nil {
		return nil, err
	}
	ackSize, ackError := vp.portconnection.Read(getAck)
	if ackSize != 1 || getAck[0] != 0x06 || ackError != nil {
		return nil, errors.New("Did not get ACK back from loop command")
	}
	loop, commerr := vp.getBytes(99, 10)
	if commerr != nil {
		return nil, commerr
	}
	vlp := loopPacketRaw{}
	crcVal := checkCRC(loop)
	if crcVal != 0 {
		return nil, errors.New("CRC check failed; the packet from the console was bad.")
	}
	b := loop[3:99]
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &vlp)
	if err != nil {
		return nil, err
	}
	vp.portconnection.Flush()
	vp.wakeUpConsole()
	return &vlp, nil
}

func (vp *VPConsole) cancelLoop() {
	_, err := vp.portconnection.Write([]byte("\r"))
	if err != nil {
		log.Fatal(err)
	}
}
