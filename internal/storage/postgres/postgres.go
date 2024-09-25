package postgres

import (
	"database/sql/driver"
	"errors"
	"net"
)

type NullInet struct {
	IP    net.IP
	Valid bool
}

//it is necessary take 'this' in two different ways

func (ni *NullInet) Scan(value interface{}) error {
	if value == nil {
		ni.IP = nil
		ni.Valid = false
		return nil
	}

	ni.Valid = true

	ip, ok := value.(string)
	if !ok {
		return errors.New("failed to scan NullInet")
	}

	ni.IP = net.ParseIP(ip)

	return nil
}

func (ni NullInet) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}

	return ni.IP.String(), nil
}
