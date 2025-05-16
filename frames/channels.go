// https://github.com/crsf-wg/crsf/wiki/CRSF_FRAMETYPE_RC_CHANNELS_PACKED
package frames

import (
	"fmt"
	"log/slog"
)

const (
	ChannelsFrameLength        = 22 + 2 //Payload + Type + CRC
	MaxChannels                = 16
	ChannelsMask        uint16 = 0x07ff // The maximum 11-bit channel value

	//(0-1984)
	ChannelsMin = 172  //988us
	ChannelsMid = 992  //1500us
	ChannelsMax = 1811 //2012us
)

type ChannelsData struct {
	Channels []uint16
}

func UnmarshalChannels(data []byte) (ChannelsData, error) {
	d := ChannelsData{
		Channels: make([]uint16, MaxChannels),
	}

	if len(data) != ChannelsFrameLength {
		slog.Error("Invalid frame length")
		return d, ErrFrameLength
	}
	if !ValidateFrame(data) {
		return d, ErrInvalidCRC8
	}
	//TODO check correct type?

	d.Channels[0] = ((uint16(data[1]) | uint16(data[2])<<8) & ChannelsMask)
	d.Channels[1] = ((uint16(data[2])>>3 | uint16(data[3])<<5) & ChannelsMask)
	d.Channels[2] = ((uint16(data[3])>>6 | uint16(data[4])<<2 | uint16(data[5])<<10) & ChannelsMask)
	d.Channels[3] = ((uint16(data[5])>>1 | uint16(data[6])<<7) & ChannelsMask)
	d.Channels[4] = ((uint16(data[6])>>4 | uint16(data[7])<<4) & ChannelsMask)
	d.Channels[5] = ((uint16(data[7])>>7 | uint16(data[8])<<1 | uint16(data[9])<<9) & ChannelsMask)
	d.Channels[6] = ((uint16(data[9])>>2 | uint16(data[10])<<6) & ChannelsMask)
	d.Channels[7] = ((uint16(data[10])>>5 | uint16(data[11])<<3) & ChannelsMask)
	d.Channels[8] = ((uint16(data[12]) | uint16(data[13])<<8) & ChannelsMask)
	d.Channels[9] = ((uint16(data[13])>>3 | uint16(data[14])<<5) & ChannelsMask)
	d.Channels[10] = ((uint16(data[14])>>6 | uint16(data[15])<<2 | uint16(data[16])<<10) & ChannelsMask)
	d.Channels[11] = ((uint16(data[16])>>1 | uint16(data[17])<<7) & ChannelsMask)
	d.Channels[12] = ((uint16(data[17])>>4 | uint16(data[18])<<4) & ChannelsMask)
	d.Channels[13] = ((uint16(data[18])>>7 | uint16(data[19])<<1 | uint16(data[20])<<9) & ChannelsMask)
	d.Channels[14] = ((uint16(data[20])>>2 | uint16(data[21])<<6) & ChannelsMask)
	d.Channels[15] = ((uint16(data[21])>>5 | uint16(data[22])<<3) & ChannelsMask)

	return d, nil
}

// TODO: check if this is correct
func (d *ChannelsData) MarshalChannels() []byte {
	const PAYLOAD_LEN = 22
	payload := make([]byte, PAYLOAD_LEN)

	ch := d.Channels
	if len(ch) < MaxChannels {
		return nil // or handle error
	}

	payload[0] = byte(ch[0] & 0xFF)
	payload[1] = byte((ch[0] >> 8) | ((ch[1] & 0x07) << 3))
	payload[2] = byte((ch[1] >> 3) | ((ch[2] & 0x3F) << 5))
	payload[3] = byte((ch[2] >> 6) | ((ch[3] & 0x01FF) << 2))
	payload[4] = byte((ch[3] >> 6) | ((ch[4] & 0x0F) << 7))
	payload[5] = byte((ch[4] >> 4) | ((ch[5] & 0x7F) << 4))
	payload[6] = byte((ch[5] >> 7) | ((ch[6] & 0x3F) << 1))
	payload[7] = byte((ch[6] >> 6) | ((ch[7] & 0xFF) << 3))
	payload[8] = byte((ch[7] >> 5) | ((ch[8] & 0x1F) << 6))
	payload[9] = byte((ch[8] >> 2))
	payload[10] = byte((ch[8] >> 10) | ((ch[9] & 0xFF) << 1))
	payload[11] = byte((ch[9] >> 7) | ((ch[10] & 0x7F) << 4))
	payload[12] = byte((ch[10] >> 4) | ((ch[11] & 0x0F) << 7))
	payload[13] = byte((ch[11] >> 1))
	payload[14] = byte((ch[11] >> 9) | ((ch[12] & 0x3F) << 2))
	payload[15] = byte((ch[12] >> 6) | ((ch[13] & 0x1F) << 5))
	payload[16] = byte((ch[13] >> 3))
	payload[17] = byte((ch[14] & 0xFF))
	payload[18] = byte((ch[14] >> 8) | ((ch[15] & 0x07) << 3))
	payload[19] = byte((ch[15] >> 3))
	payload[20] = 0
	payload[21] = 0

	return payload
}

func (d *ChannelsData) String() string {
	builtString := ""
	for i := range d.Channels {
		if i != 0 {
			builtString = fmt.Sprintf("%s Channel%d: %d", builtString, i+1, d.Channels[i])
		} else {
			builtString = fmt.Sprintf("Channel%d: %d", i+1, d.Channels[i])
		}
	}

	return builtString
}
