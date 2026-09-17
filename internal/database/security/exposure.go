package security

type ExposureRisk struct {
	Service  string
	Port     int
	Exposed  bool
	Severity string
}

func CheckExposure(service string, port int, isPublic bool) ExposureRisk {
	sev := "INFO"
	if isPublic {
		sev = "HIGH"
	}
	return ExposureRisk{
		Service:  service,
		Port:     port,
		Exposed:  isPublic,
		Severity: sev,
	}
}
