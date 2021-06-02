package system

import "github.com/shirou/gopsutil/v3/host"

func SensorsTemperatures() ([]host.TemperatureStat, error) {
	r := []host.TemperatureStat{}
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return r, err
	}
	for _, temp := range temps {
		if temp.Temperature == 0 {
			continue
		}
		temp.SensorKey = keyFormatter(temp.SensorKey)
		r = append(r, temp)
	}
	return r, nil
}
