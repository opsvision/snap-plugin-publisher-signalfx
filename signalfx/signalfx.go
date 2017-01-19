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
	//"log"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
	"golang.org/x/net/context"
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
	initialized bool
}

// Constructor
func New() *SignalFx {
	return new(SignalFx)
}

func (s *SignalFx) init() error {

	s.initialized = true

	return nil
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

	return *policy, nil
}

/**
 * Publish metrics to SignalFx using the TOKEN found in the config
 */
func (s *SignalFx) Publish(mts []plugin.Metric, cfg plugin.Config) error {
	// Make sure we've initialized
	if !s.initialized {
		s.init()
	}

	/*
		// Set the output file
		f, err := os.OpenFile("/tmp/signalfx.debug", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		log.SetOutput(f)
		log.Printf("Inside publisher")
	*/

	// Fetch the token
	token, err := cfg.GetString("token")
	if err != nil {
		return err
	}

	// Attempt to set the hostname
	hostname, err := cfg.GetString("hostname")
	if err != nil {
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
	}

	// Iterate over the supplied metrics
	for _, m := range mts {
		var buffer bytes.Buffer

		buffer.WriteString("snap.")
		buffer.WriteString(strings.Join(m.Namespace.Strings(), "."))

		client := sfxclient.NewHTTPDatapointSink()
		client.AuthToken = token
		ctx := context.Background()
		client.AddDatapoints(ctx, []*datapoint.Datapoint{
			sfxclient.GaugeF("snap.testing", map[string]string{
				"host": hostname,
			}, float64(m.Data.(float64))),
		})
	}

	return nil
}
