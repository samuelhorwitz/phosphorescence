package phosphor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/samuelhorwitz/phosphorescence/api/common"
	"github.com/samuelhorwitz/phosphorescence/api/middleware"
	"github.com/samuelhorwitz/phosphorescence/api/models"
	"github.com/samuelhorwitz/phosphorescence/api/session"
)

func TransferPlayback(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(middleware.SessionContextKey).(*session.Session)
	if !ok {
		common.Fail(w, errors.New("No session on request context"), http.StatusUnauthorized)
		return
	}
	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		common.Fail(w, errors.New("Must specify device ID"), http.StatusBadRequest)
		return
	}
	var playState models.PlayState
	switch strings.ToLower(r.URL.Query().Get("playState")) {
	case "play":
		playState = models.PlayStatePlay
	case "pause":
		playState = models.PlayStatePause
	}
	models.TransferPlayback(r.Context(), sess, deviceID, playState)
	time.Sleep(250 * time.Millisecond)
	type returnDevicesEnvelope struct {
		devices models.SpotifyDevices
		success bool
	}
	returnDevicesChan := make(chan returnDevicesEnvelope)
	errChan := make(chan error)
	done := make(chan struct{})
	go func() {
		common.ExponentialBackoff(100*time.Millisecond, 750*time.Millisecond, func() bool {
			devices, err := models.GetDevices(r.Context(), sess)
			if err != nil {
				errChan <- err
				return true
			}
			for _, device := range devices.Devices {
				if device.ID == deviceID && device.IsActive {
					returnDevicesChan <- returnDevicesEnvelope{devices, true}
					return true
				}
			}
			returnDevicesChan <- returnDevicesEnvelope{devices, false}
			return false
		})
		close(done)
	}()
	var returnDevices models.SpotifyDevices
	for {
		select {
		case devicesEnvelope := <-returnDevicesChan:
			if devicesEnvelope.success {
				common.JSON(w, devicesEnvelope.devices)
				return
			}
			returnDevices = devicesEnvelope.devices
		case err := <-errChan:
			common.Fail(w, fmt.Errorf("Could not get devices: %s", err), http.StatusInternalServerError)
			return
		case <-done:
			common.FailWithJSON(w, fmt.Errorf("Device %s does not exist", deviceID), returnDevices, http.StatusNotFound)
			return
		}
	}
}
