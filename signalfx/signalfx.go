/*
 * http://www.apache.org/licenses/LICENSE-2.0.txt
 *
 * Copyright 2017 OpsVision Solutions
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package signalfx

import (
	"bytes"
	"fmt"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
)

const (
	NS_VENDOR = "opsvision"
	NS_PLUGIN = "signalfx"
	VERSION   = 1
)

var fileHandle *os.File

type SignalFx struct {
	token     string
	hostname  string
	namespace string
}

// Constructor
func New() *SignalFx {
	return new(SignalFx)
}

/**
 * Returns the configPolicy for the plugin
 */
func (s *SignalFx) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	// The SignalFx token
	policy.AddNewStringRule([]string{NS_VENDOR, NS_PLUGIN},
		"token",
		true)

	// The hostname to use (defaults to local hostname)
	policy.AddNewStringRule([]string{NS_VENDOR, NS_PLUGIN},
		"hostname",
		false)

	// The file name to use when debugging
	policy.AddNewStringRule([]string{NS_VENDOR, NS_PLUGIN},
		"debug-file",
		false)

	return *policy, nil
}

/**
 * Publish metrics to SignalFx using the TOKEN found in the config
 */
func (s *SignalFx) Publish(mts []plugin.Metric, cfg plugin.Config) error {
	// Enable debugging if the debug-file config property was set
	fileName, err := cfg.GetString("debug-file")
	if err != nil {
		// Open the output file
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		// Set logging output for debugging
		log.SetOutput(f)
	}

	// Fetch the token
	token, err := cfg.GetString("token")
	if err != nil {
		return err
	}
	s.token = token

	// Attempt to set the hostname
	hostname, err := cfg.GetString("hostname")
	if err != nil {
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
	}
	s.hostname = hostname

	// Iterate over the supplied metrics
	for _, m := range mts {
		var buffer bytes.Buffer

		// Convert the namespace to dot notation
		fmt.Fprintf(&buffer, "snap.%s", strings.Join(m.Namespace.Strings(), "."))
		s.namespace = buffer.String()

		// Do some type conversion and send the data
		switch v := m.Data.(type) {
		case uint:
			s.sendIntValue(int64(v))
		case uint32:
			s.sendIntValue(int64(v))
		case uint64:
			s.sendIntValue(int64(v))
		case int:
			s.sendIntValue(int64(v))
		case int32:
			s.sendIntValue(int64(v))
		case int64:
			s.sendIntValue(int64(v))
		case float32:
			s.sendFloatValue(float64(v))
		case float64:
			s.sendFloatValue(float64(v))
		default:
			log.Printf("Ignoring %T: %v\n", v, v)
			log.Printf("Contact the plugin author if you think this is an error")
		}
	}

	return nil
}

/**
 * Method for sending int64 values to SignalFx
 */
func (s *SignalFx) sendIntValue(value int64) {
	log.Printf("Sending [int64] %s -> %v", s.namespace, value)

	client := sfxclient.NewHTTPDatapointSink()
	client.AuthToken = s.token
	ctx := context.Background()
	client.AddDatapoints(ctx, []*datapoint.Datapoint{
		sfxclient.Gauge(s.namespace, map[string]string{
			"host": s.hostname,
		}, value),
	})
}

/**
 * Method for sending float64 values to SignalFx
 */
func (s *SignalFx) sendFloatValue(value float64) {
	log.Printf("Sending [float64] %s -> %v", s.namespace, value)

	client := sfxclient.NewHTTPDatapointSink()
	client.AuthToken = s.token
	ctx := context.Background()
	client.AddDatapoints(ctx, []*datapoint.Datapoint{
		sfxclient.GaugeF(s.namespace, map[string]string{
			"host": s.hostname,
		}, value),
	})
}
