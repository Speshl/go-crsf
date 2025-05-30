package crsf

import (
	"fmt"
	"strings"

	"github.com/Speshl/go-crsf/frames"
)

type CRSFData struct {
	Channels frames.ChannelsData
	CRSFTelemetry
}

type CRSFTelemetry struct {
	Gps           frames.GpsData
	Vario         frames.VarioData
	BatterySensor frames.BatterySensorData
	Barometer     frames.BarometerData
	LinkStats     frames.LinkStatsData
	LinkRx        frames.LinkRxData
	LinkTx        frames.LinkTxData
	Attitude      frames.AttitudeData
	FlightMode    frames.FlightModeData
}

func NewCRSFData() CRSFData {
	return CRSFData{}
}

func (d *CRSFData) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Channels: {%s}\n", d.Channels.String())
	fmt.Fprintf(&sb, "GPS: {%s}\n", d.Gps.String())
	fmt.Fprintf(&sb, "Vario: {%s}\n", d.Vario.String())
	fmt.Fprintf(&sb, "Battery: {%s}\n", d.BatterySensor.String())
	fmt.Fprintf(&sb, "Barometer: {%s}\n", d.Barometer.String())
	fmt.Fprintf(&sb, "LinkStats: {%s}\n", d.LinkStats.String())
	fmt.Fprintf(&sb, "LinkRx: {%s}\n", d.LinkRx.String())
	fmt.Fprintf(&sb, "LinkTx: {%s}\n", d.LinkTx.String())
	fmt.Fprintf(&sb, "Attitude: {%s}\n", d.Attitude.String())
	fmt.Fprintf(&sb, "FlightMode: {%s}", d.FlightMode.String())
	return sb.String()
}
