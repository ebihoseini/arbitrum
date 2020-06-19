/*
* Copyright 2020, Offchain Labs, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package utils

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type RPCFlags struct {
	certFile *string
	keyFile  *string
}

func AddRPCFlags(fs *flag.FlagSet) RPCFlags {
	certFile := fs.String("cert", "", "path to certificate file (if using ssl)")
	privkeyFile := fs.String("privkey", "", "path to private key file (if using ssl)")

	return RPCFlags{
		certFile: certFile,
		keyFile:  privkeyFile,
	}
}

func LaunchRPC(handler http.Handler, port string, flags RPCFlags) error {
	r := mux.NewRouter()
	r.Handle("/", handler).Methods("GET", "POST", "OPTIONS")

	headersOk := handlers.AllowedHeaders(
		[]string{"X-Requested-With", "Content-Type", "Authorization"},
	)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods(
		[]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	)
	h := handlers.CORS(headersOk, originsOk, methodsOk)(r)

	if *flags.certFile != "" && *flags.keyFile != "" {
		log.Println("Launching rpc server over https with cert", *flags.certFile, "and key", *flags.keyFile)
		return http.ListenAndServeTLS(
			":"+port,
			*flags.certFile,
			*flags.keyFile,
			h,
		)
	} else {
		log.Println("Launching rpc server over http")
		return http.ListenAndServe(
			":"+port,
			h,
		)
	}

}
