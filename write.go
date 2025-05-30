package crsf

import (
	"fmt"
	"time"

	"github.com/Speshl/go-crsf/frames"
)

func (c *CRSF) startWriter() error {
	ticker := time.NewTicker(c.opts.WriterInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-ticker.C:
			c.dataLock.RLock()
			channelsData := c.data.Channels
			c.dataLock.RUnlock()

			fullFrame, err := c.buildFrame(frames.FrameTypeChannels, channelsData.MarshalChannels())
			if err != nil {
				return fmt.Errorf("failed building frame: %w", err)
			}

			_, err = c.port.Write(fullFrame)
			if err != nil {
				return fmt.Errorf("failed writing to %s: %w", c.path, err)
			}
		}
	}
}

func (c *CRSF) buildFrame(frameType frames.FrameType, payload []byte) ([]byte, error) {
	frame := make([]byte, 2+len(payload)+1)
	frame[0] = byte(frameType)
	copy(frame[1:], payload)
	frame[len(frame)-1] = frames.GenerateCrc8Value(frame[:len(frame)-1])
	return frame, nil
}
