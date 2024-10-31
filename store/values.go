package store

import "net"

type DBvalues struct {
	ipAddress    net.IP
	port         int
	logicalClock int16
	endpoint     string
}

func NewDBValues(ipAddress net.IP, port int, clock int16, endpoint string) *DBvalues {
	return &DBvalues{
		ipAddress:    ipAddress,
		port:         port,
		logicalClock: clock,
		endpoint:     endpoint,
	}
}

func (d *DBvalues) GetIpAddress() net.IP {
	return d.ipAddress
}

func (d *DBvalues) GetListenPort() int {
	return d.port
}

func (d *DBvalues) GetLogicalClock() int16 {
	return d.logicalClock
}

func (d *DBvalues) GetOptionalEndpoint() string {
	return d.endpoint
}
