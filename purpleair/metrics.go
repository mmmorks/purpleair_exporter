package purpleair

import "github.com/prometheus/client_golang/prometheus"

var (
	uptimeDesc      *prometheus.Desc
	rssiDesc        *prometheus.Desc
	httpSendsDesc   *prometheus.Desc
	temperatureDesc *prometheus.Desc
	dewPointDesc    *prometheus.Desc
	humidityDesc    *prometheus.Desc
	iaqDesc         *prometheus.Desc
	pressureDesc    *prometheus.Desc

	pm25AqiDesc       *prometheus.Desc
	massDesc          *prometheus.Desc
	particleCountDesc *prometheus.Desc
)

func init() {
	uptimeDesc = prometheus.NewDesc("purpleair_uptime_seconds_total", "The uptime of a sensor", []string{"id"}, nil)
	rssiDesc = prometheus.NewDesc("purpleair_rssi_dbm", "A measurement of wireless signal strength", []string{"id"}, nil)
	httpSendsDesc = prometheus.NewDesc("purpleair_http_sends_total", "A counter of outbound HTTP requests", []string{"id", "status"}, nil)

	temperatureDesc = prometheus.NewDesc("purpleair_temperature_c", "A temperature measurement, as determined by a Bosch BME680 sensor", []string{"id"}, nil)
	humidityDesc = prometheus.NewDesc("purpleair_humidity_percent", "A relative humidity measurement, as determined by a Bosch BME680 sensor", []string{"id"}, nil)
	dewPointDesc = prometheus.NewDesc("purpleair_dewpoint_c", "A dew point measurement, as determined by a Bosch BME680 sensor", []string{"id"}, nil)
	iaqDesc = prometheus.NewDesc("purpleair_iaq", "Index for Air Quality, as determined by a Bosch BME680 sensor", []string{"id"}, nil)
	pressureDesc = prometheus.NewDesc("purpleair_pressure_pa", "A barometric pressure measurement, as determined by a Bosch BME680 sensor", []string{"id"}, nil)

	pm25AqiDesc = prometheus.NewDesc("purpleair_pm25_aqi", "A PM2.5 Air Quality Index value", []string{"id", "channel", "period"}, nil)

	// CF=1, ATM; 1.0, 2.5
	massDesc = prometheus.NewDesc("purpleair_mass_ugm3", "A particulate mass measurement for a particular particle size", []string{"id", "variant", "size", "channel", "period"}, nil)

	// 0.3, 0.5, 1.0, 2.5, 5.0, 10.0
	particleCountDesc = prometheus.NewDesc("purpleair_particle_count", "A count of particles per deciliter of air for a particular particle size", []string{"id", "channel", "size", "period"}, nil)
}

func deviceDescs() []*prometheus.Desc {
	return []*prometheus.Desc{
		uptimeDesc,
		rssiDesc,
		httpSendsDesc,
		temperatureDesc,
		dewPointDesc,
		humidityDesc,
		iaqDesc,
		pressureDesc,
	}
}

func periodicDescs() []*prometheus.Desc {
	return []*prometheus.Desc{
		pm25AqiDesc,
		massDesc,
		particleCountDesc,
	}
}
