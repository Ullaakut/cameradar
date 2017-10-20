// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/EtixLabs/cameradar/server/actor/pubsub"
	"github.com/EtixLabs/cameradar/server/adaptor/websocket"
	graceful "gopkg.in/tylerb/graceful.v1"
)

func main() {
	webSocketFactory := websocket.NewGorillaFactory()
	pubsub := pubsub.NewWebSocket(webSocketFactory)

	// create and setup the http server
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", pubsub.Accept)

	httpServer := &graceful.Server{
		NoSignalHandling: true,
		Server: &http.Server{
			Addr:    fmt.Sprintf("%v:%v", "0.0.0.0", 7000),
			Handler: serverMux,
		},
	}

	fmt.Printf("cameradar server listening on %v\n", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("could not start server: %v\n", err)
		os.Exit(1)
	}
}
