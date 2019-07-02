package pom

import (
	// "log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	U "github.com/therealfakemoot/go-unidecode"
)

func SanitizeName(name string) string {
	s := U.Unidecode(strings.ReplaceAll(name, `"`, ""))
	s = strings.ReplaceAll(s, `'`, "")
	s = strings.ReplaceAll(s, `,`, "")
	s = strings.ReplaceAll(s, `-`, "_")
	s = strings.ReplaceAll(s, ` `, "_")
	s = s + "_price_chaos"

	return strings.ToLower(s)
}

func NewGauge(i Item) prometheus.Gauge {
	sanitized := SanitizeName(i.TypeLine)
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "market",
		Name:      sanitized,
	})
}
