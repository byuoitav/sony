package bravia

import (
	"context"
	"errors"
	"strconv"
	"time"

	"encoding/json"

	"go.uber.org/zap"
)

func (t *TV) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	t.Log.Info("Getting volume for %v", zap.String("address", t.Address))
	toReturn := make(map[string]int)
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return toReturn, err
	}
	t.Log.Info("%v", zap.Any("parentResponse", parentResponse))

	for _, outerResult := range parentResponse.Result {

		for _, result := range outerResult {

			if result.Target == "speaker" {

				toReturn[""] = result.Volume
			}
		}
	}
	t.Log.Info("Done")

	return toReturn, nil
}

func (t *TV) SetVolume(ctx context.Context, block string, volume int) error {

	if volume > 100 || volume < 0 {
		return errors.New("Error: volume must be a value from 0 to 100!")
	}

	t.Log.Debug("Setting volume for %s to %v...", zap.String("address", t.Address), zap.Int("volume", volume))
	params := make(map[string]interface{})
	params["target"] = "speaker"
	params["volume"] = strconv.Itoa(volume)

	err := t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioVolume", params)
	if err != nil {
		return err
	}

	//do the same for the headphone
	params = make(map[string]interface{})
	params["target"] = "headphone"
	params["volume"] = strconv.Itoa(volume)

	err = t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioVolume", params)
	if err != nil {
		return err
	}

	t.Log.Debug("Done.")
	return nil
}

func (t *TV) getAudioInformation(ctx context.Context) (SonyAudioResponse, error) {
	payload := SonyTVRequest{
		Params:  []map[string]interface{}{},
		Method:  "getVolumeInformation",
		Version: "1.0",
		ID:      1,
	}

	t.Log.Info("%+v", zap.Any("payload", payload))

	resp, err := t.PostHTTPWithContext(ctx, "audio", payload)

	parentResponse := SonyAudioResponse{}

	t.Log.Info("%s", zap.Any("resp", resp))

	err = json.Unmarshal(resp, &parentResponse)
	return parentResponse, err

}

func (t *TV) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	toReturn := make(map[string]bool)
	t.Log.Info("Getting mute status for %v", zap.String("address", t.Address))
	parentResponse, err := t.getAudioInformation(ctx)
	if err != nil {
		return toReturn, err
	}

	for _, outerResult := range parentResponse.Result {
		for _, result := range outerResult {
			if result.Target == "speaker" {
				t.Log.Info("local mute: %v", zap.Bool("mute", result.Mute))
				toReturn[""] = result.Mute
			}
		}
	}

	t.Log.Info("Done")

	return toReturn, nil
}

func (t *TV) SetMute(ctx context.Context, block string, mute bool) error {
	params := make(map[string]interface{})
	params["status"] = mute

	err := t.BuildAndSendPayload(ctx, t.Address, "audio", "setAudioMute", params)
	if err != nil {
		return err
	}
	//we need to validate that it was actually muted
	blocks := []string{block}
	postStatus, err := t.Mutes(ctx, blocks)
	if err != nil {
		return err
	}

	if postStatus[""] == mute {
		return nil
	}

	//wait for a short time
	time.Sleep(10 * time.Millisecond)

	return nil
}
