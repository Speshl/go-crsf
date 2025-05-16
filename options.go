package crsf

import "time"

type CRSFOptions struct {
	BaudRate       int
	ReadTimeout    int
	ReadOnly       bool
	ReadChannels   bool
	WriterInterval time.Duration
}

type Option func(*CRSFOptions)

func GetDefaultOptions() CRSFOptions {
	return CRSFOptions{
		BaudRate:       115200,
		ReadTimeout:    1000,
		ReadOnly:       false,
		ReadChannels:   true,
		WriterInterval: 5 * time.Millisecond,
	}
}

func WithBaudRate(baud int) Option {
	return func(o *CRSFOptions) {
		o.BaudRate = baud
	}
}

func WithTimeout(timeout int) Option {
	return func(o *CRSFOptions) {
		o.ReadTimeout = timeout
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(o *CRSFOptions) {
		o.ReadOnly = readOnly
	}
}

func WithReadChannels(readChannels bool) Option {
	return func(o *CRSFOptions) {
		o.ReadChannels = readChannels
	}
}

func WithWriterInterval(interval time.Duration) Option {
	return func(o *CRSFOptions) {
		o.WriterInterval = interval
	}
}

func getOptions(opts []Option) CRSFOptions {
	options := GetDefaultOptions()
	for i := range opts {
		opts[i](&options)
	}
	return options
}
