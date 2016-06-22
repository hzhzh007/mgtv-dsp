package logic

type Creative struct {
	Id          int
	Url         string
	Height      int
	Width       int
	Type        int
	Duration    int
	LandingPage string
	Click       []string    `yaml:"click_url"`
	MonitorUrl  Impressions `yaml:"monitor_url"`
}
