// This is an NetApp exporter for getting data from NetApp SAN
// Author: Thomas Elsgaard <thomas.elsgaard@trucecommerce.com>
package main

import (
	"log/slog"
	"maragu.dev/env"
	"net"
	"net/http"
	"os"
	"strconv"
	"vmware_exporter/config"
	"vmware_exporter/prober"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	_ = env.Load()

	host := env.GetStringOrDefault("HOST", "0.0.0.0")
	port := env.GetIntOrDefault("PORT", 9141)

	address := net.JoinHostPort(host, strconv.Itoa(port))

	envConfig := config.Config{
		VcenterUser: env.GetStringOrDefault("VCENTER_USER", ""),
		VcenterPass: env.GetStringOrDefault("VCENTER_PASS", ""),
		VcenterDC:   env.GetStringOrDefault("VCENTER_DC", ""),
	}

	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		prober.Handler(w, r, envConfig, logger)
	})

	logger.Info("VMWare Exporter Starting", "binding_address", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		slog.Error("Error starting server")
	}

}
