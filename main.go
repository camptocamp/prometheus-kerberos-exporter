package main

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/camptocamp/prometheus-kerberos-exporter/internal/exporter"
)

// Config : binary flags
type Config struct {
	Version        bool   `short:"v" long:"version" description:"Show version."`
	Username       string `short:"u" long:"username" description:"A username to use to connect to kerberos server." env:"KERBEROS_USER" required:"true"`
	Realm          string `short:"r" long:"realm" description:"A realm to use to connect to kerberos server." env:"KERBEROS_REALM" required:"true"`
	KeytabFile     string `short:"k" long:"keytab" description:"A keytab file to use to connect to kerberos server." env:"KERBEROS_KEYTAB_FILE" required:"true"`
	Servers        string `short:"s" long:"server" description:"A list of servers to connect to. (separated by commas)" env:"KERBEROS_SERVERS" required:"true"`
	ScrapeInterval string `long:"scrape-interval" description:"Duration between two scrapes." env:"KERBEROS_SCRAPE_INTERVAL" default:"5s"`
	ListenAddress  string `long:"listen-address" description:"Address to listen on for web interface and telemetry." env:"KERBEROS_LISTEN_ADDRESS" default:"0.0.0.0:9889"`
	MetricPath     string `long:"metric-path" description:"Path under which to expose metrics." env:"KERBEROS_METRIC_PATH" default:"/metrics"`
	Verbose        bool   `long:"verbose" description:"Enable debug mode" env:"KERBEROS_VERBOSE"`
}

var (
	// VERSION, BUILD_DATE, GIT_COMMIT are filled in by the build script
	version    = "<<< filled in by build >>>"
	buildDate  = "<<< filled in by build >>>"
	commitSha1 = "<<< filled in by build >>>"
)

func main() {
	var c Config
	parser := flags.NewParser(&c, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	log.Printf("Kerberos Metrics Exporter %s    build date: %s    sha1: %s    Go: %s",
		version, buildDate, commitSha1,
		runtime.Version(),
	)

	if c.Verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling debug output")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if c.Version {
		return
	}

	interval, err := time.ParseDuration(c.ScrapeInterval)
	if err != nil {
		log.Fatalf("Failed to parse scrape interval duration: %s", err)
	}

	servers := strings.Split(strings.ReplaceAll(c.Servers, " ", ""), ",")

	exp, err := exporter.NewKerberosExporter(c.Realm, c.Username, c.KeytabFile)
	//func NewKerberosExporter(realm string, username string, keytabFile string) (e *Exporter, err error) {
	if err != nil {
		log.Fatalf("Failed to initialize exporter: %s", err)
	}

	go exp.Scrape(servers, interval)

	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kerberos_exporter_build_info",
		Help: "Kerberos exporter build informations",
	}, []string{"version", "commit_sha", "build_date", "golang_version"})
	buildInfo.WithLabelValues(version, commitSha1, buildDate, runtime.Version()).Set(1)
	prometheus.MustRegister(buildInfo)

	http.Handle(c.MetricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`
<html>
<head><title>Prometheus Kerberos Exporter v` + version + `</title></head>
<body>
<h1>Prometheus Kerberos Exporter ` + version + `</h1>
<p><a href='` + c.MetricPath + `'>Metrics</a></p>
</body>
</html>
						`))
		if err != nil {
			log.Infof("An error occured while writing http response: %v", err)
		}
	})

	log.Infof("Providing metrics at %s%s", c.ListenAddress, c.MetricPath)
	log.Fatal(http.ListenAndServe(c.ListenAddress, nil))
}
