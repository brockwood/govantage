package main

import (
	"log"
	"math"
	"time"

	"github.com/brockwood/govantage/vantagepro"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	radcon  = math.Pi / 180.0
	degcon  = 180.0 / math.Pi
	mphtoms = 0.44704
	MyDB    = `davis`
)

func main() {
	c, cerr := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://wx.rockhouse.org:8086",
	})
	if cerr != nil {
		log.Fatalln("Error: ", cerr)
	}
	defer c.Close()
	timeBetweenPackets := 2 * time.Second
	console, consoleerr := vantagepro.ConnectToConsole("/dev/ttyUSB0", time.Millisecond*500, 19200)
	if consoleerr != nil {
		log.Fatal("Error setting up connection to Vantage Pro 2:  " + consoleerr.Error())
	} else {
		log.Println("Connected to Vantage Pro 2... streaming data to InfluxDB.")
	}

	for true {
		starttime := time.Now()
		packet, packeterr := console.GetLoopPacket()
		if packeterr != nil {
			log.Fatal("Boom on packet!")
		}
		// Prep influxdb data
		observationTime := time.Now()
		bp, bperr := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  MyDB,
			Precision: "s",
		})
		if bperr != nil {
			log.Fatalln("Error: ", bperr)
		}
		// ** Create points **
		tags := map[string]string{"station": "davis1"}
		// Temperature
		fields := map[string]interface{}{
			"inside":  packet.GetInsideTemp(),
			"outside": packet.GetOutsideTemp(),
		}
		pt, pterr := client.NewPoint("temperature", tags, fields, observationTime)
		if pterr != nil {
			log.Fatalln("Error: ", pterr)
		}
		bp.AddPoint(pt)
		// Humidity
		fields = map[string]interface{}{
			"inside":  float64(packet.InsideHumidity),
			"outside": float64(packet.OutsideHumidity),
		}
		pt, pterr = client.NewPoint("humidity", tags, fields, observationTime)
		if pterr != nil {
			log.Fatalln("Error: ", pterr)
		}
		bp.AddPoint(pt)
		// Wind Stuff
		msWindSpeed := float64(packet.WindSpeed) * mphtoms
		dirToRad := float64(packet.WindDir) * radcon
		ucomp := (-math.Abs(msWindSpeed)) * math.Sin(dirToRad)
		vcomp := (-math.Abs(msWindSpeed)) * math.Cos(dirToRad)
		fields = map[string]interface{}{
			"windspeed":       float64(packet.WindSpeed),
			"winddirection":   float64(packet.WindDir),
			"tenminwindspeed": float64(packet.TenMinWindAvg),
			"ucomp":           ucomp,
			"vcomp":           vcomp,
		}
		pt, pterr = client.NewPoint("wind", tags, fields, observationTime)
		if pterr != nil {
			log.Fatalln("Error: ", pterr)
		}
		bp.AddPoint(pt)
		// Rain Stuff
		fields = map[string]interface{}{
			"rainrate":  packet.GetRainRate(),
			"stormrain": packet.GetStormRain(),
			"rainday":   packet.GetDayRain(),
			"rainmonth": packet.GetMonthRain(),
			"rainyear":  packet.GetYearRain(),
		}
		pt, pterr = client.NewPoint("rain", tags, fields, observationTime)
		if pterr != nil {
			log.Fatalln("Error: ", pterr)
		}
		bp.AddPoint(pt)
		// Pressure Stuff
		fields = map[string]interface{}{
			"barometer": packet.GetBarometer(),
			"trend":     packet.BarTrend,
		}
		pt, pterr = client.NewPoint("barometer", tags, fields, observationTime)
		if pterr != nil {
			log.Fatalln("Error: ", pterr)
		}
		bp.AddPoint(pt)
		// Writing to influx
		cwrite := c.Write(bp)
		if cwrite != nil {
			log.Fatalln("Error: ", cwrite)
		}
		duration := time.Now().Sub(starttime)
		timeToSleep := timeBetweenPackets - duration
		time.Sleep(timeToSleep)
	}
}
