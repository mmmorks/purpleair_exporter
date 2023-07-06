package purpleair

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Scraper struct {
	targetIP  net.IP
	oneSecond bool
	twoMinute bool
}

func NewScraper(targetIP net.IP, oneSecond bool, twoMinute bool) *Scraper {
	s := &Scraper{
		targetIP:  targetIP,
		oneSecond: oneSecond,
		twoMinute: twoMinute,
	}

	return s
}

func (s Scraper) Describe(descs chan<- *prometheus.Desc) {
	for _, desc := range deviceDescs() {
		descs <- desc
	}
	for _, desc := range periodicDescs() {
		descs <- desc
	}
}

func (s Scraper) Collect(metrics chan<- prometheus.Metric) {
	ctx := context.Background()

	if s.oneSecond {
		if err := s.get(ctx, "/json?live=true", "1s", true, metrics); err != nil {
			outputErr(metrics, err, deviceDescs())
			outputErr(metrics, err, periodicDescs())
		}
	}

	if s.twoMinute {
		if err := s.get(ctx, "/json?live=false", "2m", !s.oneSecond, metrics); err != nil {
			if !s.oneSecond {
				outputErr(metrics, err, deviceDescs())
			}
			outputErr(metrics, err, periodicDescs())
		}
	}
}

func outputErr(metrics chan<- prometheus.Metric, err error, descs []*prometheus.Desc) {
	for _, desc := range descs {
		metrics <- prometheus.NewInvalidMetric(desc, err)
	}
}

func (s Scraper) get(ctx context.Context, path string, period string, includeDevice bool, metrics chan<- prometheus.Metric) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+s.targetIP.String()+path, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	var data struct {
		SensorID            string  `json:"SensorId"`
		DateTime            string  `json:"DateTime"`
		Geo                 string  `json:"Geo"`
		Mem                 int     `json:"Mem"`
		Memfrag             int     `json:"memfrag"`
		Memfb               int     `json:"memfb"`
		Memcs               int     `json:"memcs"`
		ID                  int     `json:"Id"`
		Lat                 float64 `json:"lat"`
		Lon                 float64 `json:"lon"`
		Adc                 float64 `json:"Adc"`
		Loggingrate         int     `json:"loggingrate"`
		Place               string  `json:"place"`
		Version             string  `json:"version"`
		Uptime              int     `json:"uptime"`
		Rssi                int     `json:"rssi"`
		Period              int     `json:"period"`
		Httpsuccess         int     `json:"httpsuccess"`
		Httpsends           int     `json:"httpsends"`
		Hardwareversion     string  `json:"hardwareversion"`
		Hardwarediscovered  string  `json:"hardwarediscovered"`
		CurrentTempF        int     `json:"current_temp_f"`
		CurrentHumidity     int     `json:"current_humidity"`
		CurrentDewpointF    int     `json:"current_dewpoint_f"`
		Pressure            float64 `json:"pressure"`
		CurrentTempF680     float64 `json:"current_temp_f_680"`
		CurrentHumidity680  float64 `json:"current_humidity_680"`
		CurrentDewpointF680 float64 `json:"current_dewpoint_f_680"`
		Pressure680         float64 `json:"pressure_680"`
		Gas680              float64 `json:"gas_680"`
		P25Aqic             string  `json:"p25aqic"`
		P25AqicB            string  `json:"p25aqic_b"`
		P03Um               float64 `json:"p_0_3_um"`
		P03UmB              float64 `json:"p_0_3_um_b"`
		P05Um               float64 `json:"p_0_5_um"`
		P05UmB              float64 `json:"p_0_5_um_b"`
		P100Um              float64 `json:"p_10_0_um"`
		P100UmB             float64 `json:"p_10_0_um_b"`
		P10Um               float64 `json:"p_1_0_um"`
		P10UmB              float64 `json:"p_1_0_um_b"`
		P25Um               float64 `json:"p_2_5_um"`
		P25UmB              float64 `json:"p_2_5_um_b"`
		P50Um               float64 `json:"p_5_0_um"`
		P50UmB              float64 `json:"p_5_0_um_b"`
		PaLatency           int     `json:"pa_latency"`
		Pm100Atm            float64 `json:"pm10_0_atm"`
		Pm100AtmB           float64 `json:"pm10_0_atm_b"`
		Pm100Cf1            float64 `json:"pm10_0_cf_1"`
		Pm100Cf1B           float64 `json:"pm10_0_cf_1_b"`
		Pm10Atm             float64 `json:"pm1_0_atm"`
		Pm10AtmB            float64 `json:"pm1_0_atm_b"`
		Pm10Cf1             float64 `json:"pm1_0_cf_1"`
		Pm10Cf1B            float64 `json:"pm1_0_cf_1_b"`
		Pm25Aqi             float64 `json:"pm2.5_aqi"`
		Pm25AqiB            float64 `json:"pm2.5_aqi_b"`
		Pm25Atm             float64 `json:"pm2_5_atm"`
		Pm25AtmB            float64 `json:"pm2_5_atm_b"`
		Pm25Cf1             float64 `json:"pm2_5_cf_1"`
		Pm25Cf1B            float64 `json:"pm2_5_cf_1_b"`
		Response            int     `json:"response"`
		ResponseDate        int     `json:"response_date"`
		Latency             int     `json:"latency"`
		Wlstate             string  `json:"wlstate"`
		Status0             int     `json:"status_0"`
		Status1             int     `json:"status_1"`
		Status2             int     `json:"status_2"`
		Status3             int     `json:"status_3"`
		Status4             int     `json:"status_4"`
		Status5             int     `json:"status_5"`
		Status7             int     `json:"status_7"`
		Status8             int     `json:"status_8"`
		Status9             int     `json:"status_9"`
		Ssid                string  `json:"ssid"`
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return err
	}

	/*
		timestamp, err := time.Parse("2006/01/02T15:04:05z", data.DateTime)
		if err != nil {
			return err
		}
		tx := func(m prometheus.Metric) {
			metrics <- prometheus.NewMetricWithTimestamp(timestamp, m)
		}
	*/

	tx := func(m prometheus.Metric) {
		metrics <- m
	}

	if includeDevice {
		tx(prometheus.MustNewConstMetric(uptimeDesc, prometheus.CounterValue, float64(data.Uptime), data.SensorID))
		tx(prometheus.MustNewConstMetric(rssiDesc, prometheus.GaugeValue, float64(data.Rssi), data.SensorID))
		tx(prometheus.MustNewConstMetric(httpSendsDesc, prometheus.CounterValue, float64(data.Httpsends), data.SensorID, "total"))
		tx(prometheus.MustNewConstMetric(httpSendsDesc, prometheus.CounterValue, float64(data.Httpsuccess), data.SensorID, "success"))
		tx(prometheus.MustNewConstMetric(httpSendsDesc, prometheus.CounterValue, float64(data.Httpsends-data.Httpsuccess), data.SensorID, "failure"))

		tx(prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, ftoc(data.CurrentTempF680), data.SensorID))
		tx(prometheus.MustNewConstMetric(humidityDesc, prometheus.GaugeValue, data.CurrentHumidity680, data.SensorID))
		tx(prometheus.MustNewConstMetric(dewPointDesc, prometheus.GaugeValue, ftoc(data.CurrentDewpointF680), data.SensorID))
		tx(prometheus.MustNewConstMetric(iaqDesc, prometheus.GaugeValue, data.Gas680, data.SensorID))
		tx(prometheus.MustNewConstMetric(pressureDesc, prometheus.GaugeValue, data.Pressure680*100, data.SensorID))
	}

	tx(prometheus.MustNewConstMetric(pm25AqiDesc, prometheus.GaugeValue, data.Pm25Aqi, data.SensorID, "A", period))
	tx(prometheus.MustNewConstMetric(pm25AqiDesc, prometheus.GaugeValue, data.Pm25AqiB, data.SensorID, "B", period))

	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm10Cf1, data.SensorID, "CF1", "1.0", "A", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm10Cf1B, data.SensorID, "CF1", "1.0", "B", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm10Atm, data.SensorID, "ATM", "1.0", "A", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm10AtmB, data.SensorID, "ATM", "1.0", "B", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm25Cf1, data.SensorID, "CF1", "2.5", "A", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm25Cf1B, data.SensorID, "CF1", "2.5", "B", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm25Atm, data.SensorID, "ATM", "2.5", "A", period))
	tx(prometheus.MustNewConstMetric(massDesc, prometheus.GaugeValue, data.Pm25AtmB, data.SensorID, "ATM", "2.5", "B", period))

	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P03Um, data.SensorID, "A", "0.3", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P03UmB, data.SensorID, "B", "0.3", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P05Um, data.SensorID, "A", "0.5", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P05UmB, data.SensorID, "B", "0.5", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P10Um, data.SensorID, "A", "1.0", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P10UmB, data.SensorID, "B", "1.0", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P25Um, data.SensorID, "A", "2.5", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P25UmB, data.SensorID, "B", "2.5", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P50Um, data.SensorID, "A", "5.0", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P50UmB, data.SensorID, "B", "5.0", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P100Um, data.SensorID, "A", "10.0", period))
	tx(prometheus.MustNewConstMetric(particleCountDesc, prometheus.GaugeValue, data.P100UmB, data.SensorID, "B", "10.0", period))

	return nil
}

func ftoc(f float64) float64 {
	return (f - 32) * 5 / 9
}
