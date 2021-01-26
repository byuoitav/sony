package bravia

/*
func (t *Display) SetMute(ctx context.Context, block string, mute bool) error {
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
*/
