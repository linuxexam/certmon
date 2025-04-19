package main

import "time"

type Cert struct {
	Host       string    `json:"host"`
	Port       int       `json:"port"`
	UpdateTime time.Time `json:"updateTime"`
	DaysLeft   int       `json:"daysLeft"`
}

func (c *Cert) Update() error {
	return nil
}
