package nmap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	family "github.com/Ullaakut/nmap/pkg/osfamilies"
)

func TestParseTime(t *testing.T) {
	ts := Timestamp{}

	err := ts.ParseTime("invalid")
	if err == nil {
		t.Errorf("expected strconv.ParseInt: parsing \"invalid\": invalid syntax got %v", err)
	}
}

func TestFormatTime(t *testing.T) {
	originalStr := "123456789"
	ts := Timestamp{}

	err := ts.ParseTime(originalStr)
	if err != nil {
		panic(err)
	}

	result := ts.FormatTime()

	if result != originalStr {
		t.Errorf("expected %s got %s", originalStr, result)
	}
}

func TestOSFamily(t *testing.T) {
	osc := OSClass{
		Family: "Linux",
	}

	if osc.OSFamily() != family.Linux {
		t.Errorf("expected OSClass.OSFamily() to be equal to %v, got %v", family.Linux, osc.OSFamily())
	}
}

func TestParseTableXML(t *testing.T) {
	expectedTable := map[string]string{
		"key":         "AAAAB3NzaC1yc2EAAAABIwAAAQEAwVKoTY/7GFG7BmKkG6qFAHY/f3ciDX2MXTBLMEJP0xyUJsoy/CVRYw2b4qUB/GCJ5lh2InP+LVnPD3ZdtpyIvbS0eRZs/BH+mVLGh9xA/wOEUiiCfzQRsHj1xn7cqeWViAzQtdGluk/5CVAvr1FU3HNaaWkg7KQOSiKAzgDwCBtQhlgI40xdXgbqMkrHeP4M1p4MxoEVpZMe4oObACWwazeHP/Xas1vy5rbnmE59MpEZaA8t7AfGlW4MrVMhAB1JsFMdd0qFLpy/l93H3ptSlx1+6PQ5gUyjhmDUjMR+k6fb0yOeGdOrjN8IrWPmebZRFBjK5aCJwubgY/03VsSBMQ==",
		"fingerprint": "79f809acd4e232421049d3bd208285ec",
		"type":        "ssh-rsa",
		"bits":        "2048",
	}

	input := []byte(fmt.Sprintf(
		`<table><elem key="key">%s</elem><elem key="fingerprint">%s</elem><elem key="type">%s</elem><elem key="bits">%s</elem></table>`,
		expectedTable["key"],
		expectedTable["fingerprint"],
		expectedTable["type"],
		expectedTable["bits"],
	))

	var table Table

	err := xml.Unmarshal(input, &table)
	if err != nil {
		panic(err)
	}

	if table["key"] != expectedTable["key"] {
		t.Errorf("expected %v got %v", expectedTable["key"], table["key"])
	}

	if table["fingerprint"] != expectedTable["fingerprint"] {
		t.Errorf("expected %v got %v", expectedTable["fingerprint"], table["fingerprint"])
	}

	if table["type"] != expectedTable["type"] {
		t.Errorf("expected %v got %v", expectedTable["type"], table["type"])
	}

	if table["bits"] != expectedTable["bits"] {
		t.Errorf("expected %v got %v", expectedTable["bits"], table["bits"])
	}
}

func TestFormatTableXML(t *testing.T) {
	table := Table(map[string]string{
		"key":         "AAAAB3NzaC1yc2EAAAABIwAAAQEAwVKoTY/7GFG7BmKkG6qFAHY/f3ciDX2MXTBLMEJP0xyUJsoy/CVRYw2b4qUB/GCJ5lh2InP+LVnPD3ZdtpyIvbS0eRZs/BH+mVLGh9xA/wOEUiiCfzQRsHj1xn7cqeWViAzQtdGluk/5CVAvr1FU3HNaaWkg7KQOSiKAzgDwCBtQhlgI40xdXgbqMkrHeP4M1p4MxoEVpZMe4oObACWwazeHP/Xas1vy5rbnmE59MpEZaA8t7AfGlW4MrVMhAB1JsFMdd0qFLpy/l93H3ptSlx1+6PQ5gUyjhmDUjMR+k6fb0yOeGdOrjN8IrWPmebZRFBjK5aCJwubgY/03VsSBMQ==",
		"fingerprint": "79f809acd4e232421049d3bd208285ec",
		"type":        "ssh-rsa",
		"bits":        "2048",
	})

	expectedXML := [][]byte{
		[]byte("<Table>"),
		[]byte(fmt.Sprintf(`<elem key="key">%s</elem>`, table["key"])),
		[]byte(fmt.Sprintf(`<elem key="fingerprint">%s</elem>`, table["fingerprint"])),
		[]byte(fmt.Sprintf(`<elem key="type">%s</elem>`, table["type"])),
		[]byte(fmt.Sprintf(`<elem key="bits">%s</elem>`, table["bits"])),
		[]byte("</Table>"),
	}

	XML, err := xml.Marshal(table)
	if err != nil {
		panic(err)
	}

	for _, expectedXMLElement := range expectedXML {
		if !bytes.Contains(XML, expectedXMLElement) {
			t.Errorf("missing %s in %s", expectedXMLElement, XML)
		}
	}
}

func TestParseRunXML(t *testing.T) {
	tests := []struct {
		inputFile string

		expectedResult *Run
		expectedError  error
	}{
		{
			inputFile: "test_xml/scan01.xml",

			expectedResult: &Run{
				Args:             "nmap -A -v -oX sample-03.xml freshmeat.net sourceforge.net nmap.org kernel.org openbsd.org netbsd.org google.com gmail.com",
				Scanner:          "nmap",
				StartStr:         "Sun Jan 27 21:10:02 2008",
				Version:          "4.53",
				XMLOutputVersion: "1.01",
				ScanInfo: ScanInfo{
					NumServices: 1714,
					Protocol:    "tcp",
					Services:    "1-1027,1029-1033,1040,1043,1050,1058-1059,1067-1068,1076,1080,1083-1084,1103,1109-1110,1112,1127,1139,1155,1158,1178,1212,1214,1220,1222,1234,1241,1248,1270,1337,1346-1381,1383-1552,1600,1650-1652,1661-1672,1680,1720,1723,1755,1761-1764,1827,1900,1935,1984,1986-2028,2030,2032-2035,2038,2040-2049,2053,2064-2065,2067-2068,2105-2106,2108,2111-2112,2120-2121,2201,2232,2241,2301,2307,2401,2430-2433,2500-2501,2564,2600-2605,2627-2628,2638,2766,2784,2809,2903,2998,3000-3001,3005-3006,3025,3045,3049,3052,3064,3086,3128,3141,3264,3268-3269,3292,3299,3306,3333,3372,3389,3397-3399,3421,3455-3457,3462,3531,3632,3689,3900,3984-3986,3999-4000,4002,4008,4045,4125,4132-4133,4144,4199,4224,4321,4333,4343,4444,4480,4500,4557,4559,4660,4662,4672,4899,4987,4998,5000-5003,5009-5011,5050,5060,5100-5102,5145,5190-5193,5232,5236,5300-5305,5308,5400,5405,5432,5490,5500,5510,5520,5530,5540,5550,5555,5560,5631-5632,5679-5680,5713-5717,5800-5803,5900-5903,5977-5979,5997-6009,6017,6050,6101,6103,6105-6106,6110-6112,6141-6148,6222,6346-6347,6400-6401,6502,6543-6544,6547-6548,6558,6588,6662,6665-6670,6699-6701,6881,6969,7000-7010,7070,7100,7200-7201,7273,7326,7464,7597,7937-7938,8000,8007,8009,8021,8076,8080-8082,8118,8123,8443,8770,8888,8892,9040,9050-9051,9090,9100-9107,9111,9152,9535,9876,9991-9992,9999-10000,10005,10082-10083,11371,12000,12345-12346,13701-13702,13705-13706,13708-13718,13720-13722,13782-13783,14141,15126,15151,16080,16444,16959,17007,17300,18000,18181-18185,18187,19150,20005,22273,22289,22305,22321,22370,26208,27000-27010,27374,27665,31337,31416,32770-32780,32786-32787,38037,38292,43188,44334,44442-44443,47557,49400,50000,50002,54320,61439-61441,65301",
					Type:        "syn",
				},
				Start: Timestamp(time.Unix(1201479002, 0)),
				Verbose: Verbose{
					Level: 1,
				},
				Stats: Stats{
					Finished: Finished{
						Time:    Timestamp(time.Unix(1201481569, 0)),
						TimeStr: "Sun Jan 27 21:52:49 2008",
					},
					Hosts: HostStats{
						Up:    8,
						Total: 8,
						Down:  0,
					},
				},
				Hosts: []Host{
					{
						IPIDSequence: IPIDSequence{
							Class:  "All zeros",
							Values: "0,0,0,0,0,0",
						},
						OS: OS{
							PortsUsed: []PortUsed{
								{
									State: "open",
									Proto: "tcp",
									ID:    80,
								},
								{
									State: "closed",
									Proto: "tcp",
									ID:    443,
								},
							},
							Matches: []OSMatch{
								{
									Name:     "MicroTik RouterOS 2.9.46",
									Accuracy: 94,
									Line:     14788,
								},
								{
									Name:     "Linksys WRT54GS WAP (Linux kernel)",
									Accuracy: 94,
									Line:     8292,
								},
								{
									Name:     "Linux 2.4.18 - 2.4.32 (likely embedded)",
									Accuracy: 94,
									Line:     8499,
								},
								{
									Name:     "Linux 2.4.21 - 2.4.33",
									Accuracy: 94,
									Line:     8624,
								},
								{
									Name:     "Linux 2.4.27",
									Accuracy: 94,
									Line:     8675,
								},
								{
									Name:     "Linux 2.4.28 - 2.4.30",
									Accuracy: 94,
									Line:     8693,
								},
								{
									Name:     "Linux 2.6.5 - 2.6.18",
									Accuracy: 94,
									Line:     11411,
								},
								{
									Name:     "Linux 2.6.8",
									Accuracy: 94,
									Line:     11485,
								},
								{
									Name:     "WebVOIZE 120 IP phone",
									Accuracy: 94,
									Line:     18921,
								},
								{
									Name:     "Linux 2.4.2 (Red Hat 7.1)",
									Accuracy: 91,
									Line:     8533,
								},
							},
							Fingerprints: []OSFingerprint{
								{
									Fingerprint: fingerprint,
								},
							},
							Classes: []OSClass{
								{
									Vendor:       "MikroTik",
									OSGeneration: "2.X",
									Type:         "software router",
									Accuracy:     94,
									Family:       "RouterOS",
								},
								{
									Vendor:       "Linksys",
									OSGeneration: "2.4.X",
									Type:         "WAP",
									Accuracy:     94,
									Family:       "Linux",
								},
								{
									Vendor:       "Linux",
									OSGeneration: "2.4.X",
									Type:         "general purpose",
									Accuracy:     94,
									Family:       "Linux",
								},
								{
									Vendor:       "Linux",
									OSGeneration: "2.6.X",
									Type:         "general purpose",
									Accuracy:     94,
									Family:       "Linux",
								},
								{
									Vendor:   "WebVOIZE",
									Type:     "VoIP phone",
									Accuracy: 94,
									Family:   "embedded",
								},
								{
									Vendor:       "D-Link",
									OSGeneration: "2.4.X",
									Type:         "WAP",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:   "Inventel",
									Type:     "WAP",
									Accuracy: 91,
									Family:   "embedded",
								},
								{
									Vendor:   "USRobotics",
									Type:     "broadband router",
									Accuracy: 91,
									Family:   "embedded",
								},
								{
									Vendor:       "Linux",
									OSGeneration: "2.4.X",
									Type:         "broadband router",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:       "Linux",
									OSGeneration: "2.4.X",
									Type:         "WAP",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:       "Linux",
									OSGeneration: "2.4.X",
									Type:         "media device",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:       "Linux",
									OSGeneration: "2.4.X",
									Type:         "VoIP gateway",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:   "Netgear",
									Type:     "WAP",
									Accuracy: 91,
									Family:   "embedded",
								},
								{
									Vendor:   "QLogic",
									Type:     "switch",
									Accuracy: 91,
									Family:   "embedded",
								},
								{
									Vendor:       "Sharp",
									OSGeneration: "2.4.X",
									Type:         "PDA",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:       "FON",
									OSGeneration: "2.6.X",
									Type:         "WAP",
									Accuracy:     91,
									Family:       "Linux",
								},
								{
									Vendor:       "FON",
									OSGeneration: "2.4.X",
									Type:         "WAP",
									Accuracy:     90,
									Family:       "Linux",
								},
								{
									Vendor:   "Belkin",
									Type:     "WAP",
									Accuracy: 90,
									Family:   "embedded",
								},
								{
									Vendor:   "Asus",
									Type:     "WAP",
									Accuracy: 90,
									Family:   "embedded",
								},
								{
									Vendor:       "Netgear",
									OSGeneration: "2.4.X",
									Type:         "WAP",
									Accuracy:     90,
									Family:       "Linux",
								},
								{
									Vendor:   "Xerox",
									Type:     "printer",
									Accuracy: 90,
									Family:   "embedded",
								},
								{
									Vendor:       "Aladdin",
									OSGeneration: "2.4.X",
									Type:         "security-misc",
									Accuracy:     89,
									Family:       "Linux",
								},
								{
									Vendor:   "Occam",
									Type:     "VoIP gateway",
									Accuracy: 89,
									Family:   "embedded",
								},
								{
									Vendor:   "Roku",
									Type:     "media device",
									Accuracy: 89,
									Family:   "embedded",
								},
								{
									Vendor:   "Siemens",
									Type:     "WAP",
									Accuracy: 89,
									Family:   "Linux",
								},
								{
									Vendor:       "3Com",
									OSGeneration: "2.4.X",
									Type:         "broadband router",
									Accuracy:     89,
									Family:       "Linux",
								},
								{
									Vendor:       "Dream Multimedia",
									OSGeneration: "2.6.X",
									Type:         "media device",
									Accuracy:     89,
									Family:       "Linux",
								},
								{
									Vendor:       "Iomega",
									OSGeneration: "2.6.X",
									Type:         "storage-misc",
									Accuracy:     89,
									Family:       "Linux",
								},
							},
						},
						Status: Status{
							State:  "up",
							Reason: "reset",
						},
						TCPSequence: TCPSequence{
							Index:      242,
							Difficulty: "Good luck!",
							Values:     "457B276,4584FC8,161C122C,161B185F,1605EA95,1614C498",
						},
						TCPTSSequence: TCPTSSequence{
							Class:  "other",
							Values: "3FB03AA9,3FB03C75,45B26360,45B2636A,45B26374,45B2637E",
						},
						Times: Times{
							SRTT: "269788",
							RTT:  "41141",
							To:   "434352",
						},
						Trace: Trace{
							Proto: "tcp",
							Port:  80,
							Hops: []Hop{
								{
									TTL:    1,
									RTT:    "1.83",
									IPAddr: "192.168.254.254",
								},
								{
									TTL:    2,
									RTT:    "18.95",
									IPAddr: "200.217.89.32",
								},
								{
									TTL:    3,
									RTT:    "18.33",
									IPAddr: "200.217.30.250",
									Host:   "gigabitethernet5-1.80-cto-rn-rotd-02.telemar.net.br",
								},
								{
									TTL:    4,
									RTT:    "45.05",
									IPAddr: "200.97.65.250",
									Host:   "pos15-1-nbv-pe-rotd-03.telemar.net.br",
								},
								{
									TTL:    5,
									RTT:    "43.49",
									IPAddr: "200.223.131.13",
									Host:   "pos6-0-nbv-pe-rotn-01.telemar.net.br",
								},
								{
									TTL:    6,
									RTT:    "91.27",
									IPAddr: "200.223.131.205",
									Host:   "so-0-2-0-0-arc-rj-rotn-01.telemar.net.br",
								},
								{
									TTL:    8,
									RTT:    "191.87",
									IPAddr: "200.223.131.110",
									Host:   "PO0-3.ARC-RJ-ROTN-01.telemar.net.br",
								},
								{
									TTL:    9,
									RTT:    "177.30",
									IPAddr: "208.173.90.89",
									Host:   "bpr2-so-5-2-0.miamimit.savvis.net",
								},
								{
									TTL:    10,
									RTT:    "181.50",
									IPAddr: "208.172.97.169",
									Host:   "cr2-pos-0-3-1-0.miami.savvis.net",
								},
								{
									TTL:    11,
									RTT:    "336.43",
									IPAddr: "206.24.210.70",
									Host:   "cr1-loopback.sfo.savvis.net",
								},
								{
									TTL:    12,
									RTT:    "245.32",
									IPAddr: "204.70.200.229",
									Host:   "er1-te-1-0-1.SanJose3Equinix.savvis.net",
								},
								{
									TTL:    13,
									RTT:    "238.47",
									IPAddr: "204.70.200.210",
									Host:   "hr1-te-2-0-0.santaclarasc4.savvis.net",
								},
								{
									TTL:    14,
									RTT:    "322.90",
									IPAddr: "204.70.200.217",
									Host:   "hr1-te-2-0-0.santaclarasc9.savvis.net",
								},
								{
									TTL:    15,
									RTT:    "330.96",
									IPAddr: "204.70.203.146",
								},
								{
									TTL:    16,
									RTT:    "342.57",
									IPAddr: "66.35.194.59",
									Host:   "csr2-ve242.santaclarasc8.savvis.net",
								},
								{
									TTL:    17,
									RTT:    "248.22",
									IPAddr: "66.35.210.202",
								},
								{
									TTL:    18,
									RTT:    "238.36",
									IPAddr: "66.35.250.168",
									Host:   "freshmeat.net",
								},
							},
						},
						Uptime: Uptime{
							Seconds:  206,
							Lastboot: "Sun Jan 27 21:43:11 2008",
						},
						Addresses: []Address{
							{
								Addr:     "66.35.250.168",
								AddrType: "ipv4",
							},
						},
						ExtraPorts: []ExtraPort{
							{
								State: "filtered",
								Count: 1712,
								Reasons: []Reason{
									{
										Reason: "host-prohibiteds",
										Count:  1712,
									},
								},
							},
						},
						Hostnames: []Hostname{
							{
								Name: "freshmeat.net",
								Type: "PTR",
							},
						},
						Ports: []Port{
							{
								ID:       80,
								Protocol: "tcp",
								Service: Service{
									Name:          "http",
									ExtraInfo:     "(Unix) PHP/4.4.7",
									Method:        "probed",
									Product:       "Apache httpd",
									Version:       "1.3.39",
									Configuration: 10,
								},
								State: State{
									State:     "open",
									Reason:    "syn-ack",
									ReasonTTL: 45,
								},
								Scripts: []Script{
									{
										ID:     "robots.txt",
										Output: "User-Agent: * /img/ /redir/ ",
									},
									{
										ID:     "HTML title",
										Output: "freshmeat.net: Welcome to freshmeat.net",
									},
								},
							},
							{
								ID:       443,
								Protocol: "tcp",
								Service: Service{
									Name:          "https",
									Method:        "table",
									Configuration: 3,
								},
								State: State{
									State:     "closed",
									Reason:    "reset",
									ReasonTTL: 46,
								},
							},
						},
					},
				},
				TaskBegin: []Task{
					{
						Time: Timestamp(time.Unix(1201479013, 0)),
						Task: "Ping Scan",
					},
					{
						Time: Timestamp(time.Unix(1201479014, 0)),
						Task: "Parallel DNS resolution of 8 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201479015, 0)),
						Task: "System CNAME DNS resolution of 4 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201479016, 0)),
						Task: "SYN Stealth Scan",
					},
					{
						Time: Timestamp(time.Unix(1201480879, 0)),
						Task: "Service scan",
					},
					{
						Time: Timestamp(time.Unix(1201481006, 0)),
						Task: "Traceroute",
					},
					{
						Time: Timestamp(time.Unix(1201481028, 0)),
						Task: "Traceroute",
					},
					{
						Time: Timestamp(time.Unix(1201481059, 0)),
						Task: "Parallel DNS resolution of 85 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201481070, 0)),
						Task: "System CNAME DNS resolution of 8 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201481086, 0)),
						Task: "SCRIPT ENGINE",
					},
				},
				TaskProgress: []TaskProgress{
					{
						Percent:   3.22,
						Remaining: 903,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201479949, 0)),
						Time:      Timestamp(time.Unix(1201479046, 0)),
					},
					{
						Percent:   56.66,
						Remaining: 325,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201479767, 0)),
						Time:      Timestamp(time.Unix(1201479442, 0)),
					},
					{
						Percent:   77.02,
						Remaining: 225,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201479995, 0)),
						Time:      Timestamp(time.Unix(1201479770, 0)),
					},
					{
						Percent:   81.95,
						Remaining: 215,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480212, 0)),
						Time:      Timestamp(time.Unix(1201479996, 0)),
					},
					{
						Percent:   86.79,
						Remaining: 182,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480395, 0)),
						Time:      Timestamp(time.Unix(1201480213, 0)),
					},
					{
						Percent:   87.84,
						Remaining: 172,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480433, 0)),
						Time:      Timestamp(time.Unix(1201480260, 0)),
					},
					{
						Percent:   91.65,
						Remaining: 129,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480564, 0)),
						Time:      Timestamp(time.Unix(1201480435, 0)),
					},
					{
						Percent:   94.43,
						Remaining: 91,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480656, 0)),
						Time:      Timestamp(time.Unix(1201480565, 0)),
					},
					{
						Percent:   96.35,
						Remaining: 62,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480720, 0)),
						Time:      Timestamp(time.Unix(1201480658, 0)),
					},
					{
						Percent:   97.76,
						Remaining: 39,
						Task:      "SYN Stealth Scan",
						Etc:       Timestamp(time.Unix(1201480760, 0)),
						Time:      Timestamp(time.Unix(1201480721, 0)),
					},
				},
				TaskEnd: []Task{
					{
						Time:      Timestamp(time.Unix(1201479014, 0)),
						Task:      "Ping Scan",
						ExtraInfo: "8 total hosts",
					},
					{
						Time: Timestamp(time.Unix(1201479015, 0)),
						Task: "Parallel DNS resolution of 8 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201479016, 0)),
						Task: "System CNAME DNS resolution of 4 hosts.",
					},
					{
						Time:      Timestamp(time.Unix(1201480878, 0)),
						Task:      "SYN Stealth Scan",
						ExtraInfo: "8570 total ports",
					},
					{
						Time:      Timestamp(time.Unix(1201480984, 0)),
						Task:      "Service scan",
						ExtraInfo: "20 services on 5 hosts",
					},
					{
						Time: Timestamp(time.Unix(1201481028, 0)),
						Task: "Traceroute",
					},
					{
						Time: Timestamp(time.Unix(1201481059, 0)),
						Task: "Traceroute",
					},
					{
						Time: Timestamp(time.Unix(1201481070, 0)),
						Task: "Parallel DNS resolution of 85 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201481086, 0)),
						Task: "System CNAME DNS resolution of 8 hosts.",
					},
					{
						Time: Timestamp(time.Unix(1201481197, 0)),
						Task: "SCRIPT ENGINE",
					},
				},
			},

			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.inputFile, func(t *testing.T) {
			rawXML, err := ioutil.ReadFile(test.inputFile)
			if err != nil {
				t.Fatal(err)
			}

			result, err := Parse(rawXML)

			// Remove rawXML before comparing
			result.rawXML = []byte{}

			compareResults(t, test.expectedResult, result)

			if err != test.expectedError {
				t.Errorf("expected %v got %v", test.expectedError, err)
			}
		})
	}
}

func compareResults(t *testing.T, expected, got *Run) {
	if got.Args != expected.Args {
		t.Errorf("unexpected arguments, expected %v got %v", expected.Args, got.Args)
	}

	if got.ProfileName != expected.ProfileName {
		t.Errorf("unexpected arguments, expected %v got %v", expected.ProfileName, got.ProfileName)
	}

	if got.Scanner != expected.Scanner {
		t.Errorf("unexpected arguments, expected %v got %v", expected.Scanner, got.Scanner)
	}

	if got.StartStr != expected.StartStr {
		t.Errorf("unexpected arguments, expected %v got %v", expected.StartStr, got.StartStr)
	}

	if !reflect.DeepEqual(got.Debugging, expected.Debugging) {
		t.Errorf("unexpected debugging, expected %+v got %+v", expected.Debugging, got.Debugging)
	}

	if !reflect.DeepEqual(got.ScanInfo, expected.ScanInfo) {
		t.Errorf("unexpected scan info, expected %+v got %+v", expected.ScanInfo, got.ScanInfo)
	}

	if !reflect.DeepEqual(got.Start, expected.Start) {
		t.Errorf("unexpected start time, expected %+v got %+v", expected.Start, got.Start)
	}

	if !reflect.DeepEqual(got.Targets, expected.Targets) {
		t.Errorf("unexpected targets, expected %+v got %+v", expected.Targets, got.Targets)
	}

	if len(expected.TaskBegin) != len(got.TaskBegin) {
		t.Errorf("unexpected tasks begin entries, expected to have %d entries, got %d instead", len(expected.TaskBegin), len(got.TaskBegin))
	} else {
		for idx := range expected.TaskBegin {
			if !reflect.DeepEqual(got.TaskBegin[idx], expected.TaskBegin[idx]) {
				t.Errorf("unexpected task begin entry, expected %+v got %+v", expected.TaskBegin[idx], got.TaskBegin[idx])
			}
		}
	}

	if len(expected.TaskProgress) != len(got.TaskProgress) {
		t.Errorf("unexpected tasks progress entries, expected to have %d entries, got %d instead", len(expected.TaskProgress), len(got.TaskProgress))
	} else {
		for idx := range expected.TaskProgress {
			if !reflect.DeepEqual(got.TaskProgress[idx], expected.TaskProgress[idx]) {
				t.Errorf("unexpected task progress entry, expected %+v got %+v", expected.TaskProgress[idx], got.TaskProgress[idx])
			}
		}
	}

	if len(expected.TaskEnd) != len(got.TaskEnd) {
		t.Errorf("unexpected tasks end entries, expected to have %d entries, got %d instead", len(expected.TaskEnd), len(got.TaskEnd))
	} else {
		for idx := range expected.TaskEnd {
			if !reflect.DeepEqual(got.TaskEnd[idx], expected.TaskEnd[idx]) {
				t.Errorf("unexpected task end entry, expected %+v got %+v", expected.TaskEnd[idx], got.TaskEnd[idx])
			}
		}
	}

	if len(expected.Hosts) != len(got.Hosts) {
		t.Errorf("unexpected number of hosts, expected to have %d hosts, got %d instead", len(expected.Hosts), len(got.Hosts))
	} else {
		for idx := range expected.Hosts {
			if expected.Hosts[idx].Comment != got.Hosts[idx].Comment {
				t.Errorf("unexpected host comment, expected %v got %v", expected.Hosts[idx].Comment, got.Hosts[idx].Comment)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Addresses, got.Hosts[idx].Addresses) {
				t.Errorf("unexpected host addresses, expected %+v got %+v", expected.Hosts[idx].Addresses, got.Hosts[idx].Addresses)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Distance, got.Hosts[idx].Distance) {
				t.Errorf("unexpected host distance, expected %+v got %+v", expected.Hosts[idx].Distance, got.Hosts[idx].Distance)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].ExtraPorts, got.Hosts[idx].ExtraPorts) {
				t.Errorf("unexpected host extra ports, expected %+v got %+v", expected.Hosts[idx].ExtraPorts, got.Hosts[idx].ExtraPorts)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].HostScripts, got.Hosts[idx].HostScripts) {
				t.Errorf("unexpected host host scripts, expected %+v got %+v", expected.Hosts[idx].HostScripts, got.Hosts[idx].HostScripts)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Hostnames, got.Hosts[idx].Hostnames) {
				t.Errorf("unexpected host host names, expected %+v got %+v", expected.Hosts[idx].Hostnames, got.Hosts[idx].Hostnames)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].IPIDSequence, got.Hosts[idx].IPIDSequence) {
				t.Errorf("unexpected host IPIDSequence, expected %+v got %+v", expected.Hosts[idx].IPIDSequence, got.Hosts[idx].IPIDSequence)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].OS, got.Hosts[idx].OS) {
				t.Errorf("unexpected host OS, expected %+v got %+v", expected.Hosts[idx].OS, got.Hosts[idx].OS)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Ports, got.Hosts[idx].Ports) {
				t.Errorf("unexpected host ports, expected %+v got %+v", expected.Hosts[idx].Ports, got.Hosts[idx].Ports)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Smurfs, got.Hosts[idx].Smurfs) {
				t.Errorf("unexpected host smurfs, expected %+v got %+v", expected.Hosts[idx].Smurfs, got.Hosts[idx].Smurfs)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].StartTime, got.Hosts[idx].StartTime) {
				t.Errorf("unexpected host start time, expected %+v got %+v", expected.Hosts[idx].StartTime, got.Hosts[idx].StartTime)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Status, got.Hosts[idx].Status) {
				t.Errorf("unexpected host status, expected %+v got %+v", expected.Hosts[idx].Status, got.Hosts[idx].Status)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].TCPSequence, got.Hosts[idx].TCPSequence) {
				t.Errorf("unexpected host TCPSequence, expected %+v got %+v", expected.Hosts[idx].TCPSequence, got.Hosts[idx].TCPSequence)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].TCPTSSequence, got.Hosts[idx].TCPTSSequence) {
				t.Errorf("unexpected host TCPTSSequence, expected %+v got %+v", expected.Hosts[idx].TCPTSSequence, got.Hosts[idx].TCPTSSequence)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Times, got.Hosts[idx].Times) {
				t.Errorf("unexpected host times, expected %+v got %+v", expected.Hosts[idx].Times, got.Hosts[idx].Times)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Trace, got.Hosts[idx].Trace) {
				t.Errorf("unexpected host trace, expected %+v got %+v", expected.Hosts[idx].Trace, got.Hosts[idx].Trace)
			}

			if !reflect.DeepEqual(expected.Hosts[idx].Uptime, got.Hosts[idx].Uptime) {
				t.Errorf("unexpected host uptime, expected %+v got %+v", expected.Hosts[idx].Uptime, got.Hosts[idx].Uptime)
			}
		}
	}
}

const fingerprint = "SCAN(V=4.53%D=1/27%OT=80%CT=443%CU=%PV=N%G=N%TM=479D25ED%P=i686-pc-linux-gnu)\nSEQ(SP=F2%GCD=1%ISR=E9%TI=Z%TS=1C)\nOPS(O1=M5B4ST11NW0%O2=M5B4ST11NW0%O3=M5B4NNT11NW0%O4=M5B4ST11NW0%O5=M5B4ST11NW0%O6=M5B4ST11)\nWIN(W1=16A0%W2=16A0%W3=16A0%W4=16A0%W5=16A0%W6=16A0)\nECN(R=Y%DF=Y%TG=40%W=16D0%O=M5B4NNSNW0%CC=N%Q=)\nT1(R=Y%DF=Y%TG=40%S=O%A=S+%F=AS%RD=0%Q=)\nT2(R=N)\nT3(R=Y%DF=Y%TG=40%W=16A0%S=O%A=S+%F=AS%O=M5B4ST11NW0%RD=0%Q=)\nT4(R=Y%DF=Y%TG=40%W=0%S=A%A=Z%F=R%O=%RD=0%Q=)\nT5(R=Y%DF=Y%TG=40%W=0%S=Z%A=S+%F=AR%O=%RD=0%Q=)\nT6(R=Y%DF=Y%TG=40%W=0%S=A%A=Z%F=R%O=%RD=0%Q=)\nT7(R=Y%DF=Y%TG=40%W=0%S=Z%A=S+%F=AR%O=%RD=0%Q=)\nU1(R=N)\nIE(R=N)\n"
