package gonmap

const (
	Closed    Status = 0x000a1
	Open      Status = 0x000b2
	Matched   Status = 0x000c3
	NoMatched Status = 0x000d4
	Unknow    Status = 0x000e5
)

type Status int

func (s Status) String() string {
	switch s {
	case Closed:
		return "Closed"
	case Open:
		return "Open"
	case Matched:
		return "Matched"
	case NoMatched:
		return "NoMatched"
	case Unknow:
		return "Unknow"
	default:
		return "Unknow"
	}
}

type Response struct {
	Raw         string
	Tls         bool
	Fingerprint *FingerPrint
}

var dnsResPonse = Response{
	Raw: "DnsServer", Tls: false, Fingerprint: &FingerPrint{Service: "dns"},
}
