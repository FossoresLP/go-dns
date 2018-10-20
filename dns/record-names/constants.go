package names

// TYPE is the type enum for DNS messages
type TYPE uint16

const (
	// A is a IPv4 DNS entry
	A TYPE = iota + 1

	// NS is a nameserver DNS entry
	NS

	// MD is obsolete => replaced by MX
	MD

	// MF is obsolete => replaced by MX
	MF

	// CNAME is a name forward
	CNAME

	// SOA defines the beginning of a zone of authority
	SOA

	// MB is experimental
	MB

	// MG is experimental
	MG

	// MR is experimental
	MR

	// NULL is experimental
	NULL

	// WKS is a well-known service entry
	WKS

	// PTR is a pointer entry
	PTR

	// HINFO is a hostinfo entry
	HINFO

	// MINFO is a mailbox info entry
	MINFO

	// MX is an email server entry
	MX

	// TXT is a generic text entry
	TXT

	// RP is the responsible person entry
	RP

	// AFSDB is the AFS database record entry
	AFSDB

	// X25 is not currently in use
	X25

	// ISDN in not currently in use
	ISDN

	// RT is not currently in use
	RT

	// NSAP is not currently in use
	NSAP

	// NSAP_PTR is not currently in use
	NSAP_PTR

	// SIG is an obsolete part of an older version of DNSSEC
	SIG

	// KEY is an obsolete part of an older version of DNSSEC
	KEY

	// PX is not currently in use
	PX

	// GPOS is a previous version of LOC
	GPOS

	// AAAA is the IPv6 DNS entry
	AAAA

	// LOC is the geographical location associated with the domain name
	LOC

	// NXT is an obsolete part of an older version of DNSSEC
	NXT

	// EID is not currently in use
	EID

	// NIMLOC is not currently in use
	NIMLOC

	// SRV is a service locator entry
	SRV

	// ATMA is not currently in use
	ATMA

	// NAPTR is a naming authority pointer
	NAPTR

	// KX is a key exchange entry
	KX

	// CERT is a certificate record
	CERT

	// A6 is an obsolete proposal for an IPv6 record
	A6

	// DNAME is similar to CNAME but forwards all subnames, too
	DNAME

	// SINK is an obsolete proposal for a generic DNS record
	SINK

	// OPT is part of EDNS
	OPT

	// APL is an experimental record for an address prefix list
	APL

	// DS is a DNSSEC delegation signer entry
	DS

	// SSHFP is a SSH public key fingerprint entry
	SSHFP

	// IPSECKEY is an IPsec key entry
	IPSECKEY

	// RRSIG is a DNSSEC signature
	RRSIG

	// NSEC is a DNSSEC next secure record entry
	NSEC

	// DNSKEY is a DNSSEC key
	DNSKEY

	// DHCID is a DHCP identifier
	DHCID

	// NSEC3 is a DNSSEC next secure record version 3 entry
	NSEC3

	// NSEC3PARAM is a DNSSEC next secure record version 3 parameters entry
	NSEC3PARAM

	// TLSA is a DANE TLSA certificate association
	TLSA

	// 53 is unassigned
	_

	// 54 is unassigned
	_

	// HIP is a host identity protocol entry
	HIP

	// 56 is unassigned
	_

	// 57 is unassigned
	_

	// 58 is unassigned
	_

	// CDS is a DNSSEC child DS entry
	CDS

	// CDNSKEY is a DNSSEC child DNSKEY entry
	CDNSKEY

	// OPENPGPKEY is a DANE OpenPGP public key entry
	OPENPGPKEY

	// SPF is an obsolete sender policy framework entry
	SPF TYPE = iota + 38

	// UINFO is an unused, IANA-reserved entry
	UINFO

	// UID is an unused, IANA-reserved entry
	UID

	// GID is an unused, IANA-reserved entry
	GID

	// UNSPEC is an unused, IANA-reserved entry
	UNSPEC

	// TKEY is a transaction key entry
	TKEY TYPE = 249

	// TSIG is a transaction signature
	TSIG TYPE = 250

	// URI is an uniform resource identifier entry
	URI TYPE = 256

	// CAA is a certificate authority authorization entry
	CAA TYPE = 257

	// TA is a DNSSEC trust authorities entry (currently only a proposal)
	TA TYPE = 32768

	// DLV is a DNSSEC lookaside validation record
	DLV TYPE = 32769
)

// TypeToInt converts a string to the according type value
func TypeToInt(name string) (uint16, bool) {
	m := map[string]uint16{
		"A":          1,
		"NS":         2,
		"MD":         3,
		"MF":         4,
		"CNAME":      5,
		"SOA":        6,
		"MB":         7,
		"MG":         8,
		"MR":         9,
		"NULL":       10,
		"WKS":        11,
		"PTR":        12,
		"HINFO":      13,
		"MINFO":      14,
		"MX":         15,
		"TXT":        16,
		"RP":         17,
		"AFSDB":      18,
		"X25":        19,
		"ISDN":       20,
		"RT":         21,
		"NSAP":       22,
		"NASP_PTR":   23,
		"SIG":        24,
		"KEY":        25,
		"PX":         26,
		"GPOS":       27,
		"AAAA":       28,
		"LOC":        29,
		"NXT":        30,
		"EID":        31,
		"NIMLOC":     32,
		"SRV":        33,
		"ATMA":       34,
		"NAPTR":      35,
		"KX":         36,
		"CERT":       37,
		"A6":         38,
		"DNAME":      39,
		"SINK":       40,
		"OPT":        41,
		"APL":        42,
		"DS":         43,
		"SSHFP":      44,
		"IPSECKEY":   45,
		"RRSIG":      46,
		"NSEC":       47,
		"DNSKEY":     48,
		"DHCID":      49,
		"NSEC3":      50,
		"NSEC3PARAM": 51,
		"TLSA":       52,
		"HIP":        55,
		"CDS":        59,
		"CDNSKEY":    60,
		"OPENPGPKEY": 61,
		"SPF":        99,
		"UINFO":      100,
		"UID":        101,
		"GID":        102,
		"UNSPEC":     103,
		"TKEY":       249,
		"TSIG":       250,
		"URI":        256,
		"CAA":        257,
		"TA":         32768,
		"DLV":        32769,
	}
	val, ok := m[name]
	return val, ok
}

// IntToType converts an uint16 to a type name
func IntToType(n uint16) (string, bool) {
	m := map[uint16]string{
		1:     "A",
		2:     "NS",
		3:     "MD",
		4:     "MF",
		5:     "CNAME",
		6:     "SOA",
		7:     "MB",
		8:     "MG",
		9:     "MR",
		10:    "NULL",
		11:    "WKS",
		12:    "PTR",
		13:    "HINFO",
		14:    "MINFO",
		15:    "MX",
		16:    "TXT",
		17:    "RP",
		18:    "AFSDB",
		19:    "X25",
		20:    "ISDN",
		21:    "RT",
		22:    "NSAP",
		23:    "NASP_PTR",
		24:    "SIG",
		25:    "KEY",
		26:    "PX",
		27:    "GPOS",
		28:    "AAAA",
		29:    "LOC",
		30:    "NXT",
		31:    "EID",
		32:    "NIMLOC",
		33:    "SRV",
		34:    "ATMA",
		35:    "NAPTR",
		36:    "KX",
		37:    "CERT",
		38:    "A6",
		39:    "DNAME",
		40:    "SINK",
		41:    "OPT",
		42:    "APL",
		43:    "DS",
		44:    "SSHFP",
		45:    "IPSECKEY",
		46:    "RRSIG",
		47:    "NSEC",
		48:    "DNSKEY",
		49:    "DHCID",
		50:    "NSEC3",
		51:    "NSEC3PARAM",
		52:    "TLSA",
		55:    "HIP",
		59:    "CDS",
		60:    "CDNSKEY",
		61:    "OPENPGPKEY",
		99:    "SPF",
		100:   "UINFO",
		101:   "UID",
		102:   "GID",
		103:   "UNSPEC",
		249:   "TKEY",
		250:   "TSIG",
		256:   "URI",
		257:   "CAA",
		32768: "TA",
		32769: "DLV",
	}
	val, ok := m[n]
	return val, ok
}

// QTYPE is the query type enum for DNS messages
type QTYPE TYPE

const (
	// AXFR is a zone transfer request
	AXFR QTYPE = iota + 252

	// MAILB requests the experimental MB, MG and MR records
	MAILB

	// MAILA requests the obsolete MD and MF records
	MAILA

	// QTYPE_ANY requests all records
	QTYPE_ANY
)

// QTypeToInt converts a string to the according type value
func QTypeToInt(name string) (uint16, bool) {
	m := map[string]uint16{
		"AXFR":  252,
		"MAILB": 253,
		"MAILA": 254,
		"ANY":   255,
	}
	val, ok := m[name]
	if !ok {
		val, ok = TypeToInt(name)
	}
	return val, ok
}

// IntToQType converts an uint16 to a type name
func IntToQType(n uint16) (string, bool) {
	m := map[uint16]string{
		252: "AXFR",
		253: "MAILB",
		254: "MAILA",
		255: "ANY",
	}
	val, ok := m[n]
	if !ok {
		val, ok = IntToType(n)
	}
	return val, ok
}

// CLASS is the class type enum for DNS messages
type CLASS uint16

const (
	// IN represents the internet
	IN CLASS = iota + 1

	// CS represents CSNET and is obsolete
	CS

	// CH is the CHAOS class
	CH

	// HS Hesiod [Dyer 87]
	HS
)

// ClassToInt converts a string to the according class value
func ClassToInt(name string) (uint16, bool) {
	m := map[string]uint16{
		"IN": 1,
		"CS": 2,
		"CH": 3,
		"HS": 4,
	}
	val, ok := m[name]
	return val, ok
}

// IntToClass converts an uint16 to a class name
func IntToClass(n uint16) (string, bool) {
	m := map[uint16]string{
		1: "IN",
		2: "CS",
		3: "CH",
		4: "HS",
	}
	val, ok := m[n]
	return val, ok
}

// QCLASS is the query class type enum for DNS messages
type QCLASS CLASS

const (
	// QCLASS_ANY requests any class
	QCLASS_ANY QCLASS = 255
)

// QClassToInt converts a string to the according class value
func QClassToInt(name string) (uint16, bool) {
	m := map[string]uint16{
		"ANY": 255,
	}
	val, ok := m[name]
	if !ok {
		val, ok = ClassToInt(name)
	}
	return val, ok
}

// IntToQClass converts an uint16 to a class name
func IntToQClass(n uint16) (string, bool) {
	m := map[uint16]string{
		255: "ANY",
	}
	val, ok := m[n]
	if !ok {
		val, ok = IntToClass(n)
	}
	return val, ok
}
