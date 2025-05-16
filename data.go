package crsf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Speshl/go-crsf/frames"
)

var (
	ErrNoPayloadLength = errors.New("payload has no length")
	ErrPaylodTooLong   = errors.New("payload length too high")
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

func (c *CRSF) String() string {
	c.dataLock.RLock()
	defer c.dataLock.RUnlock()
	return c.data.String()
}

func (c *CRSF) validateFrame() ([]byte, error) {
	addressByte, err := getByte(c.ctx, c.readChan)
	if err != nil {
		return []byte{}, err
	}

	if !frames.AddressType(addressByte).IsValid() {
		return []byte{}, frames.ErrInvalidAddressType
	}

	//slog.Info("found address type", "type", AddressType(addressByte).String())

	//[sync] [len] [type] [payload] [crc8]
	//next byte should be the length of the payload
	lengthByte, err := getByte(c.ctx, c.readChan)
	if err != nil {
		return []byte{}, err
	}

	if lengthByte == 0 {
		return []byte{}, ErrNoPayloadLength
	}

	if lengthByte > 62 {
		return []byte{}, ErrPaylodTooLong
	}

	//length should be the type + payload + CRC
	return getBytes(c.ctx, c.readChan, int(lengthByte))
}

func (c *CRSF) applyFrame(frame []byte) error {
	var err error

	if len(frame) == 0 {
		return fmt.Errorf("frame is empty")
	}

	//first byte of the full payload should be the frame type
	//slog.Debug("update looking for frame", "length", int(lengthByte), "frame", fullPayload[0], "type", FrameType(fullPayload[0]))
	switch frames.FrameType(frame[0]) {
	case frames.FrameTypeChannels:
		if !c.opts.ReadChannels {
			return nil
		}
		err = c.updateChannels(frame)
	//telemetry
	case frames.FrameTypeGPS:
		err = c.updateGps(frame)
	case frames.FrameTypeVario:
		err = c.updateVario(frame)
	case frames.FrameTypeBatterySensor:
		err = c.updateBatterySensor(frame)
	case frames.FrameTypeBarometer:
		err = c.updateBarometer(frame)
	case frames.FrameTypeLinkStats:
		err = c.updateLinkStats(frame)
	case frames.FrameTypeLinkRx:
		err = c.updateLinkRx(frame)
	case frames.FrameTypeLinkTx:
		err = c.updateLinkTx(frame)
	case frames.FrameTypeAttitude:
		err = c.updateAttitude(frame)
	case frames.FrameTypeFlightMode:
		err = c.updateFlightMode(frame)
	default:
		err = fmt.Errorf("unsupported frame type: %s", frames.FrameType(frame[0]).String())
	}

	if err != nil {
		err = fmt.Errorf("failed parsing frame: %w", err)
	}

	return err
}
