package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	U "github.com/therealfakemoot/go-unidecode"

	"github.com/therealfakemoot/pom/poe"
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

type HistogramSet struct {
	HistogramVec *prometheus.HistogramVec
	Histograms   map[poe.HistoKey]prometheus.Observer
}

func (hs HistogramSet) Add(i poe.Item) {
	hk := i.Key()
	_, ok := hs.Histograms[hk]
	if !ok {
		hs.Histograms[hk] = hs.HistogramVec.With(i.Labels())
	}
}
