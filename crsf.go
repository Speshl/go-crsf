package crsf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Speshl/go-crsf/frames"
	"github.com/albenik/go-serial/v2"
	"golang.org/x/sync/errgroup"
)

type CRSF struct {
	path string
	opts CRSFOptions

	port     *serial.Port
	readChan chan byte

	running   bool
	ctx       context.Context
	crsfGroup *errgroup.Group

	dataLock sync.RWMutex
	data     CRSFData
}

// NewCRSF("somepath", WithBaudRate(115200), WithTimeout(1000))
func NewCRSF(path string, opts ...Option) *CRSF {
	return &CRSF{
		path: path,
		opts: getOptions(opts),
	}
}

func (c *CRSF) Stop() error {
	if !c.running {
		return fmt.Errorf("crsf not started")
	}
	c.running = false
	c.ctx.Done()
	return nil
}

func (c *CRSF) Start(ctx context.Context) error {
	if c.running {
		return fmt.Errorf("crsf already started")
	}
	var err error
	c.running = true

	c.port, err = serial.Open(c.path,
		serial.WithBaudrate(c.opts.BaudRate),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(c.opts.ReadTimeout),
	)
	if err != nil {
		return fmt.Errorf("failed opening crsf %s: %w", c.path, err)
	}

	crsfGroup, groupCtx := errgroup.WithContext(ctx)
	c.ctx = groupCtx

	c.readChan = make(chan byte, 4096)

	c.crsfGroup.Go(func() error {
		defer close(c.readChan)
		slog.Info("start reading from crsf", "path", c.path)
		return c.startReader()
	})

	c.crsfGroup.Go(func() error {
		return c.startReadParser()
	})

	if !c.opts.ReadOnly {
		c.crsfGroup.Go(func() error {
			slog.Info("start writing to crsf", "path", c.path)
			return c.startWriter()
		})
	}

	crsfGroup.Wait()
	if err := c.crsfGroup.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("crsf context was cancelled", "path", c.path)
			return nil
		}
		return fmt.Errorf("crsf group error: %w", err)
	}

	return nil
}

func (c *CRSF) startReader() error {
	buff := make([]byte, 128)
	for {
		if c.ctx.Err() != nil {
			slog.Info("crsf reader context was cancelled", "path", c.path)
			return c.ctx.Err()
		}
		n, err := c.port.Read(buff)
		if err != nil {
			return fmt.Errorf("failed reading from %s: %w", c.path, err)
		}
		//slog.Debug("read bytes", "num", n, "bytes", buff[:n] /*PrintBytes(buff[:n])*/)
		for i := range buff[:n] {
			c.readChan <- buff[i]
		}
	}
}

func (c *CRSF) startReadParser() error {
	for {
		frame, err := c.validateFrame()
		if err != nil {
			if !errors.Is(err, frames.ErrInvalidAddressType) {
				slog.Warn("failed validating frame", "error", err)
			}
			continue
		}
		c.applyFrame(frame)
	}
}

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

func getBytes(ctx context.Context, readChan <-chan byte, n int) ([]byte, error) {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		b, err := getByte(ctx, readChan)
		if err != nil {
			return buf[:i], err // return what we have and the error
		}
		buf[i] = b
	}
	return buf, nil
}

func getByte(ctx context.Context, readChan <-chan byte) (byte, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case b, ok := <-readChan:
		if !ok {
			return 0, fmt.Errorf("reader stopped")
		}
		return b, nil
	}
}

func (c *CRSF) buildFrame(frameType frames.FrameType, payload []byte) ([]byte, error) {
	frame := make([]byte, 2+len(payload)+1)
	frame[0] = byte(frameType)
	copy(frame[1:], payload)
	frame[len(frame)-1] = frames.GenerateCrc8Value(frame[:len(frame)-1])
	return frame, nil
}
