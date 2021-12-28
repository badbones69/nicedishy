package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/marc-campbell/nicedishy/pkg/dishy"
	"github.com/marc-campbell/nicedishy/pkg/logger"
	"github.com/marc-campbell/nicedishy/pkg/persistence"
	"github.com/marc-campbell/nicedishy/pkg/stores"
	"go.uber.org/zap"
)

type StoreDataStatusDeviceStateRequest struct {
	UptimeSeconds int `json:"uptimeSeconds"`
}

type StoreDataStatusDeviceInfoRequest struct {
	HardwareVersion string `json:"hardwareVersion"`
	SoftwareVersion string `json:"softwareVersion"`
}

type StoreDataStatusRequest struct {
	DeviceInfo            StoreDataStatusDeviceInfoRequest  `json:"deviceInfo"`
	DeviceState           StoreDataStatusDeviceStateRequest `json:"deviceState"`
	SNR                   float64                           `json:"snr"`
	DownlinkThroughputBps float64                           `json:"downlinkThroughputBps"`
	UplinkThroughputBps   float64                           `json:"uplinkThroughputBps"`
	PopPingLatencyMs      float64                           `json:"popPingLatencyMs"`
	PopPingDropRate       float64                           `json:"popPingDropRate"`
	PercentObstructed     float64                           `json:"percentObstructed"`
	SecondsObstructed     float64                           `json:"secondsObstructed"`
}

type StoreDataSpeedRequest struct {
	Download float64 `json:"download"`
	Upload   float64 `json:"upload"`
}

type StoreDataRequest struct {
	When   string                 `json:"when"`
	Status StoreDataStatusRequest `json:"status"`
	Speed  StoreDataSpeedRequest  `json:"speed"`
}

type StoreDataResponse struct {
	Error string `json:"error,omitempty"`
}

func StoreData(w http.ResponseWriter, r *http.Request) {
	storeDataResponse := StoreDataResponse{}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		JSON(w, http.StatusInternalServerError, err)
		return
	}

	storeDataRequest := StoreDataRequest{}
	if err := json.Unmarshal(payload, &storeDataRequest); err != nil {
		logger.Error(err)
		storeDataResponse.Error = err.Error()
		JSON(w, http.StatusInternalServerError, storeDataResponse)
		return
	}

	ipAddress := r.Header.Get("True-Client-IP")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Error(err)
			storeDataResponse.Error = err.Error()
			JSON(w, http.StatusInternalServerError, storeDataResponse)
			return
		}

		ipAddress = ip
	}

	// 206.214.226.67,172.70.210.63
	// split and takje the first one
	if strings.Contains(ipAddress, ",") {
		ipAddress = strings.Split(ipAddress, ",")[0]
	}

	// update the "last received data from" date of the dish
	d := DishyFromTokenContext(r)
	if d == nil {
		JSON(w, http.StatusInternalServerError, nil)
		return
	}

	when, err := time.Parse(time.RFC3339, storeDataRequest.When)
	if err != nil {
		logger.Error(err)
		storeDataResponse.Error = err.Error()
		JSON(w, http.StatusInternalServerError, storeDataResponse)
		return
	}

	// if it's been over 30 days, lets do a new geocheck of this source
	if d.LastGeocheckAt == nil || time.Since(*d.LastGeocheckAt) > 30*24*time.Hour {
		logger.Info("geochecking",
			zap.String("dishyID", d.ID))

		geocheck, err := dishy.Geocheck(d.ID, ipAddress)
		if err != nil {
			logger.Error(err)
			storeDataResponse.Error = err.Error()
			JSON(w, http.StatusInternalServerError, storeDataResponse)
			return
		}

		if geocheck.Org != "SpaceX Services, Inc." {
			fmt.Printf("THIS IS NOT A DISHY: org = %s\n", geocheck.Org)
		}

		if err := stores.GetStore().UpdateDishyGeo(context.TODO(), d.ID, when, geocheck); err != nil {
			logger.Error(err)
			storeDataResponse.Error = err.Error()
			JSON(w, http.StatusInternalServerError, storeDataResponse)
			return
		}
	}

	metricsDB := persistence.MustGetMetricsDBSession()
	query := `insert into dishy_data (
time, dishy_id, ip_address, snr, downlink_throughput_bps, uplink_throughput_bps,
pop_ping_latency_ms, pop_ping_drop_rate, percent_obstructed, seconds_obstructed,
download_speed, upload_speed)
values
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = metricsDB.Exec(context.Background(), query, when, d.ID, ipAddress, storeDataRequest.Status.SNR,
		storeDataRequest.Status.DownlinkThroughputBps, storeDataRequest.Status.UplinkThroughputBps,
		storeDataRequest.Status.PopPingLatencyMs, storeDataRequest.Status.PopPingDropRate,
		storeDataRequest.Status.PercentObstructed, storeDataRequest.Status.SecondsObstructed,
		storeDataRequest.Speed.Download, &storeDataRequest.Speed.Upload)
	if err != nil {
		logger.Error(err)
		storeDataResponse.Error = err.Error()
		JSON(w, http.StatusInternalServerError, storeDataResponse)
		return
	}

	if err := stores.GetStore().SetDishyLastReceivedStats(context.Background(), d.ID, when); err != nil {
		logger.Error(err)
		storeDataResponse.Error = err.Error()
		JSON(w, http.StatusInternalServerError, storeDataResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
