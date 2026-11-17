package yamaha

import "fmt"

const InitFrame = "\x11000\x03"

type Command string

const (
	CmdPowerOn     Command = "07E7E"
	CmdPowerOff    Command = "07E7F"
	CmdSpeakerAOn  Command = "07EAB"
	CmdSpeakerAOff Command = "07EAC"
	CmdSpeakerBOn  Command = "07EAD"
	CmdSpeakerBOff Command = "07EAE"
	CmdMuteOn      Command = "07EA2"
	CmdMuteOff     Command = "07EA3"
	CmdVolUp       Command = "07A1A"
	CmdVolDown     Command = "07A1B"
	CmdInputMDTape Command = "07A18"
	CmdInputCDR    Command = "07A19"
)

func (c Command) Frame() []byte {
	return []byte("\x02" + string(c) + "\x03")
}

func VolumeCommand(db float64) Command {
	if db < -80.0 {
		db = -80.0
	}
	if db > 0.0 {
		db = 0.0
	}

	// -80.0 dB = 0x28 (40), each 0.5 dB step is +1
	val := int((db+80.0)*2.0) + 0x28
	return Command(fmt.Sprintf("230%02X", val))
}
