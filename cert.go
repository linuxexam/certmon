package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

type Cert struct {
	Host         string   `json:"host"`
	Port         int      `json:"port"`
	UpdateTime   JSONTime `json:"updateTime"`
	DaysLeft     int      `json:"daysLeft"`
	UpdateStatus string   `json:"updateStatus"`
}

func (c *Cert) Update() {
	cert, err := CheckCert(fmt.Sprintf("%s:%d", c.Host, c.Port), "")
	c.UpdateStatus = "ok"
	c.UpdateTime = JSONTime(time.Now())

	if err != nil {
		c.UpdateStatus = err.Error()
		return
	}
	c.DaysLeft = int(cert.NotAfter.Sub(time.Now()) / (time.Hour * 24))
}

func CheckCert(addr, customizedIP string) (*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	return cert, nil
}
