package kafka

import (
	"math/rand"

	"github.com/Pallinder/go-randomdata"
)

const SrcIOS = "iOS"
const SrcBrowser = "JavaScript"
const SrcWeb = "PHP"

type Message struct {
	WorkerID int
	Type     string
	IP       string
	Msg      string
}

func Generate(wid int) Message {
	var msgType string
	switch rand.Intn(3) {
	case 0:
		msgType = SrcIOS
	case 1:
		msgType = SrcBrowser
	case 2:
		msgType = SrcWeb
	}

	src := randomdata.IpV4Address()
	msg := randomdata.FirstName(randomdata.RandomGender)

	return Message{
		WorkerID: wid,
		Type:     msgType,
		IP:       src,
		Msg:      msg,
	}
}
