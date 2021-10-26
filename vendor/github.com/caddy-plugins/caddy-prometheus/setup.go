package metrics

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/admpub/caddy"
	"github.com/admpub/caddy/caddyhttp/httpserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	caddy.RegisterPlugin("prometheus", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

const (
	defaultPath = "/metrics"
	defaultAddr = "localhost:9180"
)

var once sync.Once

// Metrics holds the prometheus configuration.
type Metrics struct {
	next           httpserver.Handler
	addr           string // where to we listen
	useCaddyAddr   bool
	hostname       string
	path           string
	extraLabels    []extraLabel
	latencyBuckets []float64
	sizeBuckets    []float64
	// subsystem?
	once sync.Once

	handler http.Handler
}

type extraLabel struct {
	name  string
	value string
}

// NewMetrics -
func NewMetrics() *Metrics {
	return &Metrics{
		path:        defaultPath,
		addr:        defaultAddr,
		extraLabels: []extraLabel{},
	}
}

func (m *Metrics) start() error {
	m.once.Do(func() {
		m.define("")

		prometheus.MustRegister(requestCount)
		prometheus.MustRegister(requestDuration)
		prometheus.MustRegister(responseLatency)
		prometheus.MustRegister(responseSize)
		prometheus.MustRegister(responseStatus)

		if !m.useCaddyAddr {
			http.Handle(m.path, m.handler)
			go func() {
				err := http.ListenAndServe(m.addr, nil)
				if err != nil {
					log.Printf("[ERROR] Starting handler: %v", err)
				}
			}()
		}
	})
	return nil
}

func (m *Metrics) extraLabelNames() []string {
	names := make([]string, 0, len(m.extraLabels))

	for _, label := range m.extraLabels {
		names = append(names, label.name)
	}

	return names
}

func setup(c *caddy.Controller) error {
	metrics, err := parse(c)
	if err != nil {
		return err
	}

	metrics.handler = promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		ErrorHandling: promhttp.HTTPErrorOnError,
		ErrorLog:      log.New(os.Stderr, "", log.LstdFlags),
	})

	once.Do(func() {
		c.OnStartup(metrics.start)
	})

	cfg := httpserver.GetConfig(c)
	if metrics.useCaddyAddr {
		cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
			return httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
				if r.URL.Path == metrics.path {
					metrics.handler.ServeHTTP(w, r)
					return 0, nil
				}
				return next.ServeHTTP(w, r)
			})
		})
	}
	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		metrics.next = next
		return metrics
	})
	return nil
}

// prometheus {
//	address localhost:9180
// }
// Or just: prometheus localhost:9180
func parse(c *caddy.Controller) (*Metrics, error) {
	var (
		metrics *Metrics
		err     error
	)

	for c.Next() {
		if metrics != nil {
			return nil, c.Err("prometheus: can only have one metrics module per server")
		}
		metrics = NewMetrics()
		args := c.RemainingArgs()

		switch len(args) {
		case 0:
		case 1:
			metrics.addr = args[0]
		default:
			return nil, c.ArgErr()
		}
		addrSet := false
		for c.NextBlock() {
			switch c.Val() {
			case "path":
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				metrics.path = args[0]
			case "address":
				if metrics.useCaddyAddr {
					return nil, c.Err("prometheus: address and use_caddy_addr options may not be used together")
				}
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				metrics.addr = args[0]
				addrSet = true
			case "hostname":
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				metrics.hostname = args[0]
			case "use_caddy_addr":
				if addrSet {
					return nil, c.Err("prometheus: address and use_caddy_addr options may not be used together")
				}
				metrics.useCaddyAddr = true
			case "label":
				args = c.RemainingArgs()
				if len(args) != 2 {
					return nil, c.ArgErr()
				}

				labelName := strings.TrimSpace(args[0])
				labelValuePlaceholder := args[1]

				metrics.extraLabels = append(metrics.extraLabels, extraLabel{name: labelName, value: labelValuePlaceholder})
			case "latency_buckets":
				args = c.RemainingArgs()
				if len(args) < 1 {
					return nil, c.Err("prometheus: must specify 1 or more latency buckets")
				}
				metrics.latencyBuckets = make([]float64, len(args))
				for i, v := range args {
					b, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, c.Errf("prometheus: invalid bucket %q - must be a number", v)
					}
					metrics.latencyBuckets[i] = b
				}
			case "size_buckets":
				args = c.RemainingArgs()
				if len(args) < 1 {
					return nil, c.Err("prometheus: must specify 1 or more size buckets")
				}
				metrics.sizeBuckets = make([]float64, len(args))
				for i, v := range args {
					b, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, c.Errf("prometheus: invalid bucket %q - must be a number", v)
					}
					metrics.sizeBuckets[i] = b
				}
			default:
				return nil, c.Errf("prometheus: unknown item: %s", c.Val())
			}
		}
	}
	return metrics, err
}
