package main

import "github.com/skythen/apdu"

func SelectAIDAPDU() apdu.Capdu {
	return apdu.Capdu{
		Cla:  0x00,
		Ins:  0xA4,
		P1:   0x00,
		P2:   0x00,
		Data: []byte{0xA0, 0x00, 0xDA, 0xDA, 0xDA, 0xDA, 0xDA},
	}
}

func TokenReadyAPDU(chunk_size int) apdu.Capdu {
	return apdu.Capdu{
		Cla: 0xD0,
		Ins: 0x01,
		P1:  byte(chunk_size),
		P2:  0x00,
		Ne:  0x02,
	}
}

func GetTokenAPDU(id int, chunk_size int) apdu.Capdu {
	return apdu.Capdu{
		Cla: 0xD0,
		Ins: 0x02,
		P1:  byte(id),
		P2:  0x00,
		Ne:  chunk_size + 2,
	}
}
