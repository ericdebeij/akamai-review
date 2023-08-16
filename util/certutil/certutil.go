package certutil

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"
)

func Loadcerts(hostname string) (certs []*x509.Certificate, err error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	d := &net.Dialer{Timeout: 2 * time.Second}
	conn, err := tls.DialWithDialer(d, "tcp", hostname+":443", conf)
	if err != nil {
		return
	}
	defer conn.Close()

	certs = conn.ConnectionState().PeerCertificates
	return
}
