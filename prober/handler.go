package prober

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"time"
	"vmware_exporter/config"
	"vmware_exporter/vcenter"
)

func Handler(w http.ResponseWriter, r *http.Request, c config.Config, logger *slog.Logger) {

	start := time.Now()

	vmWareDatastoreCapacityGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "datastore_capacity_size",
			Help:      "Maximum capacity",
			Namespace: "vmware",
		},
		[]string{
			"dc_name",
			"ds_name",
			"ds_type",
		},
	)

	vmWareDatastoreFreeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "datastore_free_size",
			Help:      "Free space",
			Namespace: "vmware",
		},
		[]string{
			"dc_name",
			"ds_name",
			"ds_type",
		},
	)

	vmWareHostPowerStateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "host_power_state",
			Help:      "ESXi Host power status",
			Namespace: "vmware",
		},
		[]string{
			"dc_name",
			"host_name",
		},
	)

	vmWareHostConnectStateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "host_connection_state",
			Help:      "NetApp IOPS Metrics",
			Namespace: "vmware",
		},
		[]string{
			"dc_name",
			"host_name",
		},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(vmWareDatastoreCapacityGauge)
	registry.MustRegister(vmWareDatastoreFreeGauge)
	registry.MustRegister(vmWareHostPowerStateGauge)
	registry.MustRegister(vmWareHostConnectStateGauge)

	target := r.URL.Query().Get("target")
	if target == "" {
		logger.Error("VMWare Device Scrape", slog.String("request_duration_seconds", time.Since(start).String()), slog.String("err_msg", "Target parameter is missing"))
		http.Error(w, fmt.Sprintf("Target parameter is missing"), http.StatusBadRequest)
		return
	}

	dc := r.URL.Query().Get("dc")
	if dc == "" {
		logger.Error("VMWare Device Scrape", slog.String("request_duration_seconds", time.Since(start).String()), slog.String("err_msg", "DC parameter is missing"))

		http.Error(w, fmt.Sprintf("DC parameter is missing"), http.StatusBadRequest)
		return
	}

	vcenterApi := &vcenter.Model{
		User: c.VcenterUser,
		Pass: c.VcenterPass,
		DC:   dc,
		Host: target,
	}

	sessionId, err := vcenterApi.Authenticate()
	if err != nil {
		logger.Error("VMWare Device Scrape", slog.Float64("request_duration_seconds", time.Since(start).Seconds()), slog.Any("err_msg", err))
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	vcenterDC, err := vcenterApi.GetDatacenter(sessionId)
	if err != nil {
		logger.Error("VMWare Device Scrape", slog.Float64("request_duration_seconds", time.Since(start).Seconds()), slog.Any("err_msg", err))
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	datastoreMetrics, err := vcenterApi.GetDatastores(sessionId)
	if err != nil {
		logger.Error("VMWare Device Scrape", slog.Float64("request_duration_seconds", time.Since(start).Seconds()), slog.Any("err_msg", err))
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	for _, v := range datastoreMetrics {

		vmWareDatastoreCapacityGauge.WithLabelValues(vcenterDC[0].Name, v.Name, v.Type).Set(float64(v.Capacity))
		vmWareDatastoreFreeGauge.WithLabelValues(vcenterDC[0].Name, v.Name, v.Type).Set(float64(v.FreeSpace))

	}

	hostMetrics, err := vcenterApi.GetHosts(sessionId)
	if err != nil {
		logger.Error("VMWare Device Scrape", slog.Float64("request_duration_seconds", time.Since(start).Seconds()), slog.Any("err_msg", err))
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	for _, v := range hostMetrics {

		switch v.PowerState {
		case "POWERED_ON":
			vmWareHostPowerStateGauge.WithLabelValues(vcenterDC[0].Name, v.Name).Set(float64(1))

		default:
			vmWareHostPowerStateGauge.WithLabelValues(vcenterDC[0].Name, v.Name).Set(float64(0))
		}

		switch v.ConnectionState {
		case "CONNECTED":
			vmWareHostConnectStateGauge.WithLabelValues(vcenterDC[0].Name, v.Name).Set(float64(1))

		default:
			vmWareHostConnectStateGauge.WithLabelValues(vcenterDC[0].Name, v.Name).Set(float64(0))
		}
	}

	vcenterApi.LogOut(sessionId)
	logger.Info("VMWare Device Scrape", slog.Float64("request_duration_seconds", time.Since(start).Seconds()))
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
