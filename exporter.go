package exporter

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

// Exporter handles serving the metrics
type Exporter struct {
	addr                   string
	gearmanAddr            string
	ignoredEndpointRegexes regexp.Regexp
	logger                 *zap.Logger
}

// OptionsFunc is a function passed to new for setting options on a new Exporter.
type OptionsFunc func(*Exporter) error

// New creates an exporter.
func New(options ...OptionsFunc) (*Exporter, error) {
	e := &Exporter{
		addr:        "127.0.0.1:9418",
		gearmanAddr: "127.0.0.1:4730",
	}

	for _, f := range options {
		if err := f(e); err != nil {
			return nil, errors.Wrap(err, "failed to set options")
		}
	}

	if e.logger == nil {
		l, err := NewLogger()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create logger")
		}
		e.logger = l
	}

	return e, nil
}

// SetLogger creates a function that will set the logger.
// Generally only used when creating a new Exporter.
func SetLogger(l *zap.Logger) func(*Exporter) error {
	return func(e *Exporter) error {
		e.logger = l
		return nil
	}
}

// SetAddress creates a function that will set the listening address.
// Generally only used when creating a new Exporter.
func SetAddress(addr string) func(*Exporter) error {
	return func(e *Exporter) error {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return errors.Wrapf(err, "invalid address")
		}
		e.addr = net.JoinHostPort(host, port)
		return nil
	}
}

// SetGearmanAddress creates a function that will set the address to contact gearman.
// Generally only used when creating a new Exporter.
func SetGearmanAddress(addr string) func(*Exporter) error {
	return func(e *Exporter) error {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return errors.Wrapf(err, "invalid address")
		}
		e.gearmanAddr = net.JoinHostPort(host, port)
		return nil
	}
}

// SetIgnoredGearmanEndpointRegex creates a function that will set the regex
// used to ignore gearman endpoints.
// Generally only used when creating a new Exporter.
func SetIgnoredGearmanEndpointRegex(regex string) func(*Exporter) error {
	return func(e *Exporter) error {
		r, err := regexp.Compile(regex)
		if err != nil {
			return err
		}
		e.ignoredEndpointRegexes = *r
		return nil
	}
}

var healthzOK = []byte("ok\n")

func (e *Exporter) healthz(w http.ResponseWriter, r *http.Request) {
	// TODO: check if we can contact gearman?
	_, _ = w.Write(healthzOK)
}

// Run starts the http server and collecting metrics. It generally does not return.
func (e *Exporter) Run() error {

	c := e.newCollector(newGearman(e.gearmanAddr, e.ignoredEndpointRegexes))
	if err := prometheus.Register(c); err != nil {
		return errors.Wrap(err, "failed to register metrics")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<a href="/metrics">Gearman exporter</a>`)
	})

	http.HandleFunc("/healthz", e.healthz)
	http.Handle("/metrics", promhttp.Handler())
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{Addr: e.addr}
	var g errgroup.Group

	g.Go(func() error {
		// TODO: allow TLS
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-stopChan
		// XXX: should shutdown time be configurable?
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
		return nil
	})

	if err := g.Wait(); err != http.ErrServerClosed {
		return errors.Wrap(err, "failed to run server")
	}

	return nil
}
