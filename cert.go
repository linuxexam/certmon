package certmon

import (
	"crypto/tls"
	"crypto/x509"
	"time"
)

type Cert struct {
	ID           int      `json:"id"`
	Addr         string   `json:"addr"`
	DNS          string   `json:"dns"` // for customized DNS resolution, optional
	UpdateTime   JSONTime `json:"updateTime"`
	DaysLeft     int      `json:"daysLeft"`
	UpdateStatus string   `json:"updateStatus"`
}

func (c *Cert) Update() {
	cert, err := CheckCert(c.Addr, c.DNS)
	c.UpdateStatus = "ok"
	c.UpdateTime = JSONTime(time.Now())

	if err != nil {
		c.UpdateStatus = err.Error()
		return
	}
	c.DaysLeft = int(cert.NotAfter.Sub(time.Now()) / (time.Hour * 24))
}

func CheckCert(addr, dns string) (*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	return cert, nil
}
