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

// Imports
import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
	"golang.org/x/net/context"
)

// Constants
const (
	pluginVendor  = "opsvision" // plugin vendor
	pluginName    = "signalfx"  // plugin name
	pluginVersion = 1           // plugin version
)

// SignalFx object
type SignalFx struct {
	initialized bool   // Initialization flag
	token       string // SignalFx API token
	hostname    string // Hostname
	namespace   string // Metric namespace
}

// New - Constructor
func New() *SignalFx {
	return new(SignalFx)
}

func (s *SignalFx) init(cfg plugin.Config) {
	if s.initialized {
		return
	}

	// Enable debugging
	s.configDebugging(cfg)

	// Set our SignalFx API token
	s.setToken(cfg)

	// Set the hostname
	s.setHostname(cfg)

	log.Println("SignalFx Plugin Initialized")
	s.initialized = true
}

// GetConfigPolicy - Returns the configPolicy for the plugin
func (s *SignalFx) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	// The SignalFx token
	policy.AddNewStringRule([]string{pluginVendor, pluginName},
		"token",
		true)

	// The hostname to use (defaults to local hostname)
	policy.AddNewStringRule([]string{pluginVendor, pluginName},
		"hostname",
		false)

	// The file name to use when debugging
	policy.AddNewStringRule([]string{pluginVendor, pluginName},
		"debug_file",
		false)

	return *policy, nil
}

// Publish - Publishes metrics to SignalFx using the TOKEN found in the config
func (s *SignalFx) Publish(mts []plugin.Metric, cfg plugin.Config) error {
	if len(mts) > 0 {
		s.init(cfg)
	}

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

// configDebugging will configure logging if the debug_file config
// setting is present in the task file
func (s *SignalFx) configDebugging(cfg plugin.Config) {
	fileName, err := cfg.GetString("debug_file")
	if err != nil {
		// No debug_file defined, moving on
		return
	}

	// Open the output file
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return
	}

	// Set logging output for debugging
	log.SetOutput(f)
}

// setToken will set the token required by the SignalFx API
func (s *SignalFx) setToken(cfg plugin.Config) {
	log.Println("Setting token from config file")

	// Fetch the token
	token, err := cfg.GetString("token")
	if err != nil {
		log.Panic(err)
	}
	s.token = token
}

// setHostname will set the hostname from the config file, or, if absent,
// will attempt to figure out the hostname. As a last resort, we default
// to using localhost.
func (s *SignalFx) setHostname(cfg plugin.Config) {
	log.Println("Determining hostname")

	hostname, err := cfg.GetString("hostname")
	if err != nil {
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
	}
	s.hostname = hostname

	log.Printf("Using %s\n", hostname)
}

// sendIntValue - Method for sending int64 values to SignalFx
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

// sendFloatValue - Method for sending float64 values to SignalFx
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
