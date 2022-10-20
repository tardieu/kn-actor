/*
Copyright IBM Corporation 2022

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke an actor instance",
	Long:  `Invoke an actor instance.`,
	RunE:  invoke,
}

type response struct {
	Value json.RawMessage `json:"value"`
	Error *string         `json:"error"`
}

var (
	invokeServiceOption   string
	invokeInstanceOption  string
	invokeMethodOption    string
	invokeNamespaceOption string
	invokeClusterOption   string
	invokeArgumentsOption []string
)

func init() {
	rootCmd.AddCommand(invokeCmd)

	invokeCmd.Flags().StringArrayVarP(&invokeArgumentsOption, "data", "d", []string{}, "Arguments")
	invokeCmd.Flags().StringVarP(&invokeServiceOption, "service", "s", "", "Target actor service (required)")
	invokeCmd.Flags().StringVarP(&invokeInstanceOption, "instance", "i", "", "Target actor instance (required)")
	invokeCmd.Flags().StringVarP(&invokeMethodOption, "method", "m", "", "Target actor method (required)")
	invokeCmd.Flags().StringVarP(&invokeNamespaceOption, "namespace", "n", "default", "Target namespace")
	invokeCmd.Flags().StringVar(&invokeClusterOption, "cluster", "127.0.0.1.sslip.io", "Target cluster")

	invokeCmd.MarkFlagRequired("service")
	invokeCmd.MarkFlagRequired("instance")
	invokeCmd.MarkFlagRequired("method")
}

func invoke(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	url := fmt.Sprintf("http://%s.%s.%s/actor/v1/invoke/%s/%s?",
		invokeServiceOption,
		invokeNamespaceOption,
		invokeClusterOption,
		invokeInstanceOption,
		invokeMethodOption)

	body := "[" + strings.Join(invokeArgumentsOption, ",") + "]"

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("K-Session", invokeInstanceOption)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp response
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return errors.New(string(*resp.Error))
	}
	fmt.Println(string(resp.Value))

	return nil
}
