package client

type Config struct {
	Hubble Hubble
}

type Hubble struct {
	Endpoint string
	Tls      TlsConfig
}
type TlsConfig struct {
	CertificateKey string
	Certificate    string
	IsInsecure     bool
}
