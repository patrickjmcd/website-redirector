package main

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"os"
)

var redirectTo = "https://plannedparenthood.com"

func main() {

	if v, ok := os.LookupEnv("REDIRECT_TO"); ok {
		redirectTo = v
	}

	r := gin.Default()

	counter := &ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "redirects_processed",
		Description: "number of redirects processed",
		Labels:      []string{"redirectTo"},
	}

	// get global Monitor object
	m := ginmetrics.GetMonitor()
	_ = ginmetrics.GetMonitor().AddMetric(counter)

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(r)

	r.NoRoute(func(c *gin.Context) {
		_ = ginmetrics.GetMonitor().GetMetric("redirects_processed").Inc([]string{})
		c.Redirect(301, redirectTo)
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080
}
