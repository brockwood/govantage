package vantagepro

type BarametricTrend int8

func (bt BarametricTrend) String() string {
	switch bt {
	case -60:
		return "Falling Rapidly"
	case -20:
		return "Falling Slowly"
	case 0:
		return "Steady"
	case 20:
		return "Rising Slowly"
	case 60:
		return "Rising Rapidly"
	case 80:
		return "P"
	default:
		return "Not enough data to determine the trend."
	}
}

type WindDir uint8

func (wd WindDir) String() string {
	switch wd {
	case 0:
		return "N"
	case 1:
		return "NNE"
	case 2:
		return "NE"
	case 3:
		return "ENE"
	case 4:
		return "E"
	case 5:
		return "ESE"
	case 6:
		return "SE"
	case 7:
		return "SSE"
	case 8:
		return "S"
	case 9:
		return "SSW"
	case 10:
		return "SW"
	case 11:
		return "WSW"
	case 12:
		return "W"
	case 13:
		return "WNW"
	case 14:
		return "NW"
	case 15:
		return "NNW"
	default:
		return "--"
	}
}
