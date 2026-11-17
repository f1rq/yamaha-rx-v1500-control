package yamaha

import (
	"log"
	"sync"
	"time"
)

type State struct {
	Power       bool    `json:"power"`
	SpeakerA    bool    `json:"speaker_a"`
	SpeakerB    bool    `json:"speaker_b"`
	Input       string  `json:"input"`
	VolumeDB    float64 `json:"volume_db"`
	SpotifyPlay bool    `json:"spotify_playing"`
}

type Controller struct {
	transport  *SerialTransport
	mu         sync.RWMutex
	state      State
	sleepTimer *time.Timer
}

func NewController(transport *SerialTransport) *Controller {
	return &Controller{
		transport: transport,
		state: State{
			Input:    "CDR",
			VolumeDB: -45.0,
		},
	}
}

func (c *Controller) GetState() State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

func (c *Controller) SetPower(on bool) error {
	cmd := CmdPowerOff
	if on {
		cmd = CmdPowerOn
	}
	if _, err := c.transport.Send(cmd); err != nil {
		return err
	}
	c.mu.Lock()
	c.state.Power = on
	c.mu.Unlock()
	return nil
}

func (c *Controller) SetSpeakerA(on bool) error {
	cmd := CmdSpeakerAOff
	if on {
		cmd = CmdSpeakerAOn
	}
	if _, err := c.transport.Send(cmd); err != nil {
		return err
	}
	c.mu.Lock()
	c.state.SpeakerA = on
	c.mu.Unlock()
	return nil
}

func (c *Controller) SetSpeakerB(on bool) error {
	cmd := CmdSpeakerBOff
	if on {
		cmd = CmdSpeakerBOn
	}
	if _, err := c.transport.Send(cmd); err != nil {
		return err
	}
	c.mu.Lock()
	c.state.SpeakerB = on
	c.mu.Unlock()
	return nil
}

func (c *Controller) SetVolume(db float64) error {
	if _, err := c.transport.Send(VolumeCommand(db)); err != nil {
		return err
	}
	c.mu.Lock()
	c.state.VolumeDB = db
	c.mu.Unlock()
	return nil
}

func (c *Controller) OnSpotifyPlay() error {
	c.mu.Lock()
	if c.sleepTimer != nil {
		c.sleepTimer.Stop()
		c.sleepTimer = nil
	}
	c.state.SpotifyPlay = true
	c.mu.Unlock()

	log.Println("[AUTO] Turning on Yamaha...")
	if err := c.SetPower(true); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	_, _ = c.transport.Send(CmdInputCDR)
	c.mu.Lock()
	c.state.Input = "CDR"
	c.mu.Unlock()

	st := c.GetState()
	if !st.SpeakerA && !st.SpeakerB {
		_ = c.SetSpeakerA(true)
	}
	return nil
}

func (c *Controller) OnSpotifyPause(idleTimeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.state.SpotifyPlay = false
	if c.sleepTimer != nil {
		c.sleepTimer.Stop()
	}

	c.sleepTimer = time.AfterFunc(idleTimeout, func() {
		log.Println("[AUTO] No music for 10 minutes - turning off Yamaha.")
		_ = c.SetPower(false)
	})
}
