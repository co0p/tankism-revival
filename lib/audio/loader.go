package audio

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type SoundAsset struct {
	ID   string
	Path string
}

func LoadOGGPCM(assetsFS fs.FS, sampleRate int, path string) ([]byte, error) {
	file, err := assetsFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open audio file %q: %w", path, err)
	}
	defer file.Close()

	decoded, err := vorbis.DecodeWithSampleRate(sampleRate, file)
	if err != nil {
		return nil, fmt.Errorf("decode ogg %q: %w", path, err)
	}

	pcm, err := io.ReadAll(decoded)
	if err != nil {
		return nil, fmt.Errorf("read decoded audio %q: %w", path, err)
	}

	return pcm, nil
}

func LoadSoundBank(assetsFS fs.FS, sampleRate int, assets []SoundAsset) (map[string][]byte, error) {
	bank := make(map[string][]byte, len(assets))

	for _, asset := range assets {
		pcm, err := LoadOGGPCM(assetsFS, sampleRate, asset.Path)
		if err != nil {
			return nil, fmt.Errorf("load sound %q: %w", asset.ID, err)
		}
		bank[asset.ID] = pcm
	}

	return bank, nil
}
