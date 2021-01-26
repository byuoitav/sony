package bravia

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func (d *Display) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	infos, err := d.getVolumeInformation(ctx)
	if err != nil {
		return nil, err
	}

	vols := make(map[string]int, len(blocks))

	for _, block := range blocks {
		found := false
		for _, info := range infos {
			if block == info.Target {
				found = true
				vols[block] = info.Volume
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("block %q not present", block)
		}
	}

	return vols, nil
}

func (d *Display) SetVolume(ctx context.Context, block string, vol int) error {
	req := request{
		Version: "1.2",
		Method:  "setAudioVolume",
		Params: []map[string]interface{}{
			{
				"target": block,
				"volume": strconv.Itoa(vol),
				"ui":     "off",
			},
		},
	}

	_, err := d.doRequest(ctx, "audio", req)
	return err
}

func (d *Display) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	infos, err := d.getVolumeInformation(ctx)
	if err != nil {
		return nil, err
	}

	mutes := make(map[string]bool, len(blocks))

	for _, block := range blocks {
		found := false
		for _, info := range infos {
			if block == info.Target {
				found = true
				mutes[block] = info.Mute
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("block %q not present", block)
		}
	}

	return mutes, nil
}

// SetMute sets mute on all blocks, not just the given block. The bravia API does not
// currently support setting mute on a specific block.
func (d *Display) SetMute(ctx context.Context, block string, mute bool) error {
	req := request{
		Version: "1.0",
		Method:  "setAudioMute",
		Params: []map[string]interface{}{
			{
				"status": mute,
			},
		},
	}

	_, err := d.doRequest(ctx, "audio", req)
	if err != nil {
		return err
	}

	// wait for display to mute
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mutes, err := d.Mutes(ctx, []string{block})
			switch {
			case err != nil:
				return fmt.Errorf("unable to confirm mute set: %w", err)
			case mutes[block] == mute:
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("unable to confirm mute set: %w", ctx.Err())
		}
	}
}

type volumeInformation struct {
	Target    string `json:"target"`
	Volume    int    `json:"volume"`
	Mute      bool   `json:"mute"`
	MaxVolume int    `json:"maxVolume"`
	MinVolume int    `json:"minVolume"`
}

func (d *Display) getVolumeInformation(ctx context.Context) ([]volumeInformation, error) {
	req := request{
		Version: "1.0",
		Method:  "getVolumeInformation",
		Params:  []map[string]interface{}{},
	}

	res, err := d.doRequest(ctx, "audio", req)
	switch {
	case err != nil:
		return nil, err
	case len(res) < 1:
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	list, ok := res[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response: %+v", res)
	}

	var infos []volumeInformation
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		var info volumeInformation

		info.Target, ok = m["target"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		vol, ok := m["volume"].(float64)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}
		info.Volume = int(vol)

		info.Mute, ok = m["mute"].(bool)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}

		vol, ok = m["maxVolume"].(float64)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}
		info.MaxVolume = int(vol)

		vol, ok = m["minVolume"].(float64)
		if !ok {
			return nil, fmt.Errorf("unexpected response: %+v", res)
		}
		info.MinVolume = int(vol)

		infos = append(infos, info)
	}

	return infos, nil
}
