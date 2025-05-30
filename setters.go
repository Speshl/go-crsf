package crsf

import (
	"github.com/Speshl/go-crsf/frames"
)

func (c *CRSF) SetData(data CRSFData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data = data
}

func (c *CRSF) updateChannels(data []byte) error {
	dataStruct, err := frames.UnmarshalChannels(data)
	if err != nil {
		return err
	}
	c.SetChannels(dataStruct)
	return nil
}

func (c *CRSF) SetChannels(data frames.ChannelsData) {
	//slog.Debug("setting channels", "data", data.String())
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.Channels = data
}

func (c *CRSF) updateGps(data []byte) error {
	dataStruct, err := frames.UnmarshalGps(data)
	if err != nil {
		return err
	}
	c.SetGps(dataStruct)
	return nil
}

func (c *CRSF) SetGps(data frames.GpsData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.Gps = data
}

func (c *CRSF) updateVario(data []byte) error {
	dataStruct, err := frames.UnmarshalVario(data)
	if err != nil {
		return err
	}
	c.SetVario(dataStruct)
	return nil
}

func (c *CRSF) SetVario(data frames.VarioData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.Vario = data
}

func (c *CRSF) updateBatterySensor(data []byte) error {
	dataStruct, err := frames.UnmarshalBatterySensor(data)
	if err != nil {
		return err
	}
	c.SetBatterySensor(dataStruct)
	return nil
}

func (c *CRSF) SetBatterySensor(data frames.BatterySensorData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.BatterySensor = data
}

func (c *CRSF) updateBarometer(data []byte) error {
	dataStruct, err := frames.UnmarshalBarometer(data)
	if err != nil {
		return err
	}
	c.SetBarometer(dataStruct)
	return nil
}

func (c *CRSF) SetBarometer(data frames.BarometerData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.Barometer = data
}

func (c *CRSF) updateLinkStats(data []byte) error {
	dataStruct, err := frames.UnmarshalLinkStats(data)
	if err != nil {
		return err
	}
	c.SetLinkStats(dataStruct)
	return nil
}

func (c *CRSF) SetLinkStats(data frames.LinkStatsData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.LinkStats = data
}

func (c *CRSF) updateLinkRx(data []byte) error {
	dataStruct, err := frames.UnmarshalLinkRx(data)
	if err != nil {
		return err
	}
	c.SetLinkRx(dataStruct)
	return nil
}

func (c *CRSF) SetLinkRx(data frames.LinkRxData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.LinkRx = data
}

func (c *CRSF) updateLinkTx(data []byte) error {
	dataStruct, err := frames.UnmarshalLinkTx(data)
	if err != nil {
		return err
	}
	c.SetLinkTx(dataStruct)
	return nil
}

func (c *CRSF) SetLinkTx(data frames.LinkTxData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.LinkTx = data
}

func (c *CRSF) updateAttitude(data []byte) error {
	dataStruct, err := frames.UnmarshalAttitude(data)
	if err != nil {
		return err
	}
	c.SetAttitude(dataStruct)
	return nil
}

func (c *CRSF) SetAttitude(data frames.AttitudeData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.Attitude = data
}

func (c *CRSF) updateFlightMode(data []byte) error {
	dataStruct, err := frames.UnmarshalFlightMode(data)
	if err != nil {
		return err
	}
	c.SetFlightMode(dataStruct)
	return nil
}

func (c *CRSF) SetFlightMode(data frames.FlightModeData) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data.FlightMode = data
}
