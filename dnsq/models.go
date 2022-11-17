package dnsq

type Result struct {
	Domain     string
	Nameserver string
	Serial     uint32
	Diff       uint32
}

type Config struct {
	InitialNameserver string
	Nameservers       []string
	Domains           []string
}
