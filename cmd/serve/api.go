// Copyright 2019 Grendel Authors. All rights reserved.
//
// This file is part of Grendel.
//
// Grendel is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Grendel is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Grendel. If not, see <https://www.gnu.org/licenses/>.

package serve

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ubccr/grendel/api"
)

func init() {
	apiCmd.PersistentFlags().String("api-listen", fmt.Sprintf("0.0.0.0:%d", api.DefaultPort), "address to listen on")
	viper.BindPFlag("api.listen", apiCmd.PersistentFlags().Lookup("api-listen"))
	apiCmd.PersistentFlags().String("api-socket", "", "path to unix socket")
	viper.BindPFlag("api.socket_path", apiCmd.PersistentFlags().Lookup("api-socket"))
	apiCmd.PersistentFlags().String("api-cert", "", "path to ssl cert")
	viper.BindPFlag("api.cert", apiCmd.PersistentFlags().Lookup("api-cert"))
	apiCmd.PersistentFlags().String("api-key", "", "path to ssl key")
	viper.BindPFlag("api.key", apiCmd.PersistentFlags().Lookup("api-key"))

	serveCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run API server",
		Long:  `Run API server`,
		RunE: func(command *cobra.Command, args []string) error {
			ctx, _ := NewInterruptContext()
			return serveAPI(ctx)
		},
	}
)

func serveAPI(ctx context.Context) error {
	apiServer, err := api.NewServer(DB, viper.GetString("api.socket_path"), viper.GetString("api.listen"))
	if err != nil {
		return err
	}

	apiServer.KeyFile = viper.GetString("api.key")
	apiServer.CertFile = viper.GetString("api.cert")

	if err := apiServer.Serve(ctx); err != nil {
		return err
	}

	return nil
}
