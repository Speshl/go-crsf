package crsf

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Speshl/go-crsf/frames"
)

var (
	ErrNoPayloadLength = errors.New("payload has no length")
	ErrPaylodTooLong   = errors.New("payload length too high")
)

func (c *CRSF) startReader() error {
	buff := make([]byte, 128)
	for {
		n, err := c.port.Read(buff)
		if err != nil {
			return fmt.Errorf("failed reading from %s: %w", c.path, err)
		}

		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case c.readChan <- buff[:n]:
			continue
		default:
			slog.Warn("read channel is full, dropping data", "path", c.path, "data_length", n)
			continue
		}
	}
}

func (c *CRSF) startReadParser() error {
	for {
		if c.ctx.Err() != nil {
			return c.ctx.Err()
		}

		addressByte, err := c.getNextByte()
		if err != nil {
			slog.Warn("failed to get next byte", "error", err)
			continue
		}

		if !frames.AddressType(addressByte).IsValid() {
			continue
		}

		//[sync] [len] [type] [payload] [crc8]
		//next byte should be the length of the payload
		lengthByte, err := c.getNextByte()
		if err != nil {
			slog.Warn("failed to get next byte", "error", err)
			continue
		}

		if lengthByte == 0 {
			return ErrNoPayloadLength
		}

		if lengthByte > 62 {
			return ErrPaylodTooLong
		}

		//length should be the type + payload + CRC
		frameBytes, err := c.getNextBytes(int(lengthByte))
		if err != nil {
			slog.Warn("failed to get next bytes", "error", err)
			continue
		}

		err = c.applyFrame(frameBytes)
		if err != nil {
			slog.Warn("failed to apply frame", "error", err, "address", frames.AddressType(addressByte), "length", lengthByte, "frame", frameBytes)
			continue
		}
	}
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

func (c *CRSF) getNextByte() (byte, error) {
	if c.readBuffIdx < len(c.readBuff) {
		b := c.readBuff[c.readBuffIdx]
		c.readBuffIdx++
		return b, nil
	}

	select {
	case <-c.ctx.Done():
		return 0, c.ctx.Err()
	case newBuf, ok := <-c.readChan:
		if !ok {
			return 0, fmt.Errorf("read channel closed")
		}
		c.readBuff = newBuf
		if len(c.readBuff) == 0 {
			c.readBuffIdx = 0
			return 0, fmt.Errorf("read buffer is empty")
		}
		c.readBuffIdx = 1 // Reset index to 1 since showing first byte from next buffer
		return c.readBuff[0], nil
	}

}

func (c *CRSF) getNextBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	copied := 0
	for copied < length {
		remaining := len(c.readBuff) - c.readBuffIdx
		if remaining == 0 {
			select {
			case <-c.ctx.Done():
				return nil, c.ctx.Err()
			case newBuff, ok := <-c.readChan:
				if !ok {
					return nil, fmt.Errorf("read channel closed")
				}
				if len(newBuff) == 0 {
					return nil, fmt.Errorf("read buffer is empty")
				}
				c.readBuff = <-c.readChan
				c.readBuffIdx = 0
				remaining = len(c.readBuff)
			}
		}

		toCopy := min(length-copied, remaining)
		copy(bytes[copied:], c.readBuff[c.readBuffIdx:c.readBuffIdx+toCopy])
		c.readBuffIdx += toCopy
		copied += toCopy
	}
	return bytes, nil
}
