package vantagepro

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"time"
)

type archivePacketRaw struct {
	Datestamp        uint16
	Timestamp        uint16
	Outsidetemp      int16
	Highoutsidetemp  int16
	Lowoutsidetemp   int16
	Rainfall         uint16
	HighRainRate     uint16
	Barometer        uint16
	Solarradiation   uint16
	Numofwindsamples uint16
	Insidetemp       int16
	Insidehumidity   uint8
	Outsidehumidity  uint8
	Avgwindspeed     uint8
	Highwindspeed    uint8
	Highwinddir      uint8
	Prvlingwinddir   uint8
	Averageuvindex   uint8
}

func (vp *VPConsole) DumpArchiveAfter(target time.Time) {
	vp.wakeUpConsole()
}

func (vp *VPConsole) DumpArchive() ([]archivePacketRaw, error) {
	vp.wakeUpConsole()
	dmperr := vp.sendCommand("DMP")
	if dmperr != nil {
		return nil, dmperr
	}
	var pagedump []archivePacketRaw
	for pagecount := 0; pagecount < 512; pagecount++ {
		log.Println("Attempting to get page...")
		vp.portconnection.Write([]byte{ACK})
		fullpage, pageerror := vp.getBytes(267, 10)
		if pageerror != nil {
			return nil, pageerror
		}
		pageVerify := checkCRC(fullpage)
		if pageVerify != 0 {
			vp.portconnection.Write([]byte{NAK})
			return nil, errors.New("A page's CRC did not verify.")
		}
		for page := 0; page < 5; page++ {
			arp := archivePacketRaw{}
			offset := page*52 + 1
			offsetEnd := offset + 52
			pagebuf := fullpage[offset:offsetEnd]
			buf := bytes.NewReader(pagebuf)
			err := binary.Read(buf, binary.LittleEndian, &arp)
			if err != nil {
				return nil, err
			}
			year := arp.Datestamp >> 9
			month := (arp.Datestamp >> 5) & 0x0F
			day := arp.Datestamp & 0x1F
			hour := arp.Timestamp / 100
			min := arp.Timestamp % 100
			fmt.Printf("%02d/%02d/%02d %02d:%02d\n", month, day, year, hour, min)
			hiWindDir := WindDir(arp.Highwinddir)
			prevWindDir := WindDir(arp.Prvlingwinddir)
			fmt.Printf("Wind out of the %s at %dmph\n", prevWindDir.String(), arp.Avgwindspeed)
			fmt.Printf("Wind gust out of the %s at %dmph\n", hiWindDir.String(), arp.Highwindspeed)
			pagedump = append(pagedump, arp)
		}
		log.Printf("Page %d received...\n", pagecount)
	}
	return pagedump, nil
}
