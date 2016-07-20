package vantagepro

type loopPacketRaw struct {
	BarTrend        int8
	PacketType      int8
	NextRecord      uint16
	Barometer       uint16
	InsideTemp      int16
	InsideHumidity  uint8
	OutsideTemp     int16
	WindSpeed       uint8
	TenMinWindAvg   uint8
	WindDir         uint16
	ExtraTemp       [7]byte
	SoilTemp        [4]byte
	LeafTemp        [4]byte
	OutsideHumidity uint8
	ExtraHumidities [7]byte
	RainRate        uint16
	UV              uint8
	SolarRadiation  uint16
	StormRain       uint16
	StartDate       [2]byte
	DayRain         uint16
	MonthRain       uint16
	YearRain        uint16
}

func (lp *loopPacketRaw) GetInsideTemp() float64 {
	return float64(lp.InsideTemp) / 10.0
}

func (lp *loopPacketRaw) GetOutsideTemp() float64 {
	return float64(lp.OutsideTemp) / 10.0
}

func (lp *loopPacketRaw) GetBarometer() float64 {
	return float64(lp.Barometer) / 1000.0
}

func (lp *loopPacketRaw) GetRainRate() float64 {
	return float64(lp.RainRate) / 100.0
}

func (lp *loopPacketRaw) GetStormRain() float64 {
	return float64(lp.StormRain) / 100.0
}

func (lp *loopPacketRaw) GetDayRain() float64 {
	return float64(lp.DayRain) / 100.0
}

func (lp *loopPacketRaw) GetMonthRain() float64 {
	return float64(lp.MonthRain) / 100.0
}

func (lp *loopPacketRaw) GetYearRain() float64 {
	return float64(lp.YearRain) / 100.0
}
