package crsf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/albenik/go-serial/v2"
	"golang.org/x/sync/errgroup"
)

type CRSF struct {
	path string
	opts CRSFOptions

	port        *serial.Port
	readChan    chan []byte // for testing purposes
	readBuff    []byte      // buffer for reading from serial port
	readBuffIdx int

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

func (c *CRSF) String() string {
	c.dataLock.RLock()
	defer c.dataLock.RUnlock()
	return c.data.String()
}

func (c *CRSF) Stop() {
	c.ctx.Done()
}

func (c *CRSF) Start(ctx context.Context) error {
	var err error

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

	c.readChan = make(chan []byte, 1024)

	c.crsfGroup.Go(c.startReader)
	c.crsfGroup.Go(c.startReadParser)

	if !c.opts.ReadOnly {
		c.crsfGroup.Go(c.startWriter)
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
