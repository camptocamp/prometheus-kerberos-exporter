package exporter

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	username string
	realm    string
	keytab   *keytab.Keytab
	metrics  map[string]*prometheus.GaugeVec
}

var (
	metricMap = map[string]string{
		"kerberos_status_available": "kerberos_status_available",
	}
)

func NewKerberosExporter(realm string, username string, keytabFile string) (e *Exporter, err error) {
	e = &Exporter{
		username: username,
		realm:    realm,
	}

	e.keytab, err = keytab.Load(keytabFile)
	if err != nil {
		log.Fatalf("Failed to create new keytab: %s", err)
		return
	}

	e.initGauges()

	return
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrics {
		m.Describe(ch)
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, m := range e.metrics {
		m.Collect(ch)
	}
}

func (e *Exporter) Scrape(kdc []string, interval time.Duration) {
	for {
		var state int

		for _, server := range kdc {

			cfg, err := config.NewConfigFromString(`
[libdefaults]
  default_realm = ` + e.realm + `

[realms]
  ` + e.realm + ` = {
    kdc = ` + server + `:88
  }
`)
			if err != nil {
				log.Fatalf("Failed to create new config for %s: %s", server, err)
				return
			}

			cl := client.NewClientWithKeytab(e.username, e.realm, e.keytab, cfg)
			if err != nil {
				log.Fatalf("Failed to create new client for %s: %s", server, err)
				return
			}
			err = cl.Login()

			if err != nil {
				log.Errorf("Failed to connect to kerberos %s: %s", server, err)
				state = 0
			} else {
				//log.Println("Connection to kerberos was successfull !")
				state = 1
			}

			cl.Destroy()

			e.metrics["kerberos_status_available"].With(prometheus.Labels{"kdc": server}).Set(float64(state))

		}

		time.Sleep(interval)
	}

}

func (e *Exporter) initGauges() {
	e.metrics = map[string]*prometheus.GaugeVec{}

	e.metrics["kerberos_status_available"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kerberos_status_available",
		Help: "Kerberos server availability",
	}, []string{"kdc"})

	for _, m := range e.metrics {
		prometheus.MustRegister(m)
	}
}
