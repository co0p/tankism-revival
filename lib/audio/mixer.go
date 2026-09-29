package audio

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type loopPlayer struct {
	soundID string
	player  *audio.Player
}

type EffectMixer struct {
	context *audio.Context
	sounds  map[string][]byte
	active  []*audio.Player
	loops   map[string]loopPlayer
}

func NewEffectMixer(sampleRate int, sounds map[string][]byte) *EffectMixer {
	cloned := make(map[string][]byte, len(sounds))
	for id, pcm := range sounds {
		cloned[id] = pcm
	}

	return &EffectMixer{
		context: audio.NewContext(sampleRate),
		sounds:  cloned,
		loops:   make(map[string]loopPlayer),
	}
}

func (m *EffectMixer) Play(soundID string, volume float64) error {
	pcm, ok := m.sounds[soundID]
	if !ok {
		return fmt.Errorf("unknown sound id %q", soundID)
	}

	player := m.context.NewPlayerFromBytes(pcm)
	if volume < 0 {
		volume = 0
	}
	player.SetVolume(volume)
	player.Play()

	m.active = append(m.active, player)
	return nil
}

func (m *EffectMixer) SetLoop(loopID, soundID string, volume float64) error {
	if loopID == "" {
		return fmt.Errorf("loop id is required")
	}

	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}

	if existing, ok := m.loops[loopID]; ok {
		if existing.soundID == soundID {
			existing.player.SetVolume(volume)
			if !existing.player.IsPlaying() {
				existing.player.Play()
			}
			return nil
		}

		_ = existing.player.Close()
		delete(m.loops, loopID)
	}

	pcm, ok := m.sounds[soundID]
	if !ok {
		return fmt.Errorf("unknown sound id %q", soundID)
	}
	if len(pcm) == 0 {
		return fmt.Errorf("sound %q has empty PCM data", soundID)
	}

	loopStream := audio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm)))
	player, err := m.context.NewPlayer(loopStream)
	if err != nil {
		return fmt.Errorf("create loop player for %q: %w", soundID, err)
	}

	player.SetVolume(volume)
	player.Play()

	m.loops[loopID] = loopPlayer{
		soundID: soundID,
		player:  player,
	}

	return nil
}

func (m *EffectMixer) RemoveLoop(loopID string) {
	loop, ok := m.loops[loopID]
	if !ok {
		return
	}

	_ = loop.player.Close()
	delete(m.loops, loopID)
}

func (m *EffectMixer) Update() {
	kept := m.active[:0]
	for _, player := range m.active {
		if player.IsPlaying() {
			kept = append(kept, player)
			continue
		}

		_ = player.Close()
	}

	m.active = kept
}
