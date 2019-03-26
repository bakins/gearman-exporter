package exporter

import (
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// see http://gearman.org/protocol/ "Administrative Protocol"

// not concurrency safe, so caller should ensure is only in use by one goroutine
type gearman struct {
	addr string
	conn *textproto.Conn
}

func newGearman(addr string) *gearman {
	return &gearman{
		addr: addr,
	}
}

func (g *gearman) connect() error {
	// XXX: configurable timeout?
	c, err := net.DialTimeout("tcp", g.addr, 10*time.Second)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to gearman")
	}
	g.conn = textproto.NewConn(c)
	return nil
}

func (g *gearman) close() {
	_ = g.conn.Close()
	g.conn = nil
}

func (g *gearman) getConnection() (*textproto.Conn, error) {
	if g.conn != nil {
		return g.conn, nil
	}

	if err := g.connect(); err != nil {
		return nil, err
	}
	return g.conn, nil
}

type funcStatus struct {
	total   int
	running int
	workers int
}

func (g *gearman) getStatus() (map[string]*funcStatus, error) {
	c, err := g.getConnection()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get gearman connection")
	}

	id, err := c.Cmd("status")
	if err != nil {
		g.close()
		return nil, errors.Wrap(err, "failed to send status command")
	}
	c.StartResponse(id)
	defer c.EndResponse(id)

	metrics := make(map[string]*funcStatus)
	for {
		data, err := c.ReadLine()
		if err != nil {
			g.close()
			return nil, errors.Wrap(err, "failed to read status response")
		}
		if data == "." {
			break
		}
		// FUNCTION\tTOTAL\tRUNNING\tAVAILABLE_WORKERS
		parts := strings.SplitN(data, "\t", 4)
		if len(parts) != 4 {
			return nil, errors.Wrap(err, "invalid status response")
		}
		// parse each field. this is a bit brute force, but easy to understand
		s := &funcStatus{}
		s.total, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse total field in status response")
		}
		s.running, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse running field in status response")
		}
		s.workers, err = strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse workers field in status response")
		}
		metrics[parts[0]] = s
	}

	return metrics, nil
}

func (g *gearman) getVersion() (string, error) {
	c, err := g.getConnection()
	if err != nil {
		return "", errors.Wrap(err, "failed to get gearman connection")
	}

	id, err := c.Cmd("version")
	if err != nil {
		g.close()
		return "", errors.Wrap(err, "failed to send version command")
	}
	c.StartResponse(id)
	defer c.EndResponse(id)

	data, err := c.ReadLine()
	if err != nil {
		g.close()
		return "", errors.Wrap(err, "failed to read version response")
	}

	parts := strings.SplitN(data, " ", 2)
	if len(parts) != 2 {
		return "", errors.Wrap(err, "invalid version response")
	}

	if parts[0] != "OK" {
		return "", errors.Wrapf(err, "unexpected status: %s", parts[0])
	}

	return parts[1], nil
}
