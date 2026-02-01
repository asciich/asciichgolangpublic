package dnsinterfaces

type DnsRecord interface{
	GetRecordType() (string, error)
	GetName() (string, error)
	GetContent() (string, error)
	GetTtl() (int, error)
}