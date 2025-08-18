package helper

func Float64OrDefault(p *float64, def float64) float64 {
	if p != nil {
		return *p
	}
	return def
}

func IntOrDefault(p *int, def int) int {
	if p != nil {
		return *p
	}
	return def
}

func StringOrDefault(p *string, def string) string {
	if p != nil {
		return *p
	}
	return def
}

func DerefOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
