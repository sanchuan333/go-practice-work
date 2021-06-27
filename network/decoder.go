package network

type goImRaw struct {
	packageLength   string
	headerLength    string
	protocolVersion string
	operation       string
	sequence        string
	bodyRaw         string
}

func decoder(b []byte) goImRaw {
	res := goImRaw{
		packageLength:   string(b[0:3]),
		headerLength:    string(b[4:5]),
		protocolVersion: string(b[6:7]),
		operation:       string(b[8:11]),
		sequence:        string(b[12:15]),
		bodyRaw:         string(b[16:]),
	}
	return res
}
