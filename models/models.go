package models

type Service struct {
	Ip       string
	Port     int
	Username string
	Password string
}

type ScanResult struct {
	Service Service
	Result  bool
}

type IpAddr struct {
	Ip   string
	Port int
}
