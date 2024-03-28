package main

import (
	"os"

	"github.com/clausecker/nfc/v2"
	log "github.com/sirupsen/logrus"
	"github.com/skythen/apdu"
)

const chunk_size = 128

var modulation = nfc.Modulation{
	Type:     nfc.ISO14443a,
	BaudRate: nfc.Nbr106,
}

func sendCommand(device nfc.Device, command apdu.Capdu) (response *apdu.Rapdu, err error) {
	tx, err := command.Bytes()
	if err != nil {
		return &apdu.Rapdu{}, err
	}

	rx := make([]byte, command.Ne+2)

	_, err = device.InitiatorTransceiveBytes(tx, rx, 500)
	if err != nil {
		return &apdu.Rapdu{}, err
	}

	rapdu, err := apdu.ParseRapdu(rx)
	if err != nil {
		return &apdu.Rapdu{}, err
	}

	return rapdu, nil
}

func main() {
	dev, err := nfc.Open(os.Getenv(("NFC_DEVICE")))
	if err != nil {
		log.Fatal("error opening nfc device: ", err)
	}

	_, err = dev.InitiatorSelectPassiveTarget(
		modulation,
		[]byte{},
	)
	if err != nil {
		log.Fatal("error selecting target: ", err)
	}

	res, err := sendCommand(dev, SelectAIDAPDU())
	if err != nil {
		log.Fatal("error sending Select AID Command: ", err)
	}
	log.Info("App responded: ", res)

	var chunk_count = 0
	var remainder = 0

	res, err = sendCommand(dev, TokenReadyAPDU(chunk_size))
	if err != nil {
		log.Fatal("Error sending TokenReadyAPDU: ", err)
	}

	if res.SW1 == 0x6F {
		log.Fatal("token not ready")
	}

	chunk_count = int(res.Data[0])
	remainder = int(res.Data[1])

	if chunk_count == 0 && remainder == 0 {
		log.Fatal("TokenReadyAPDU timed out.")
	}

	token := []byte{}

	for i := 0; i < chunk_count; i += 1 {
		res, err := sendCommand(dev, GetTokenAPDU(i, chunk_size))
		if err != nil {
			log.Fatalf("error reading chunk #%d: %s", i, err)
		}
		if res.SW1 == 0x6F {
			log.Fatal("app respondend with error")
		}

		token = append(token, res.Data...)
	}

	res, err = sendCommand(dev, GetTokenAPDU(chunk_size, remainder))
	if err != nil {
		log.Fatalf("error reading chunk #%d: %s", chunk_count, err)
	}
	if res.SW1 == 0x6F {
		log.Fatal("app respondend with error")
	}

	token = append(token, res.Data...)

	log.Infof("token: %s", string(token))

}
