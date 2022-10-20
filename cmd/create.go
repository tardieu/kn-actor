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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an actor project",
	Long:  `Create an actor project.`,
	RunE:  create,
}

var (
	createRuntimes      = []string{"node"}
	createRuntimeOption string
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createRuntimeOption, "runtime", "r", "", "Actor runtime (required)")

	createCmd.MarkFlagRequired("runtime")
}

func contains(options []string, option string) bool {
	for _, o := range options {
		if option == o {
			return true
		}
	}
	return false
}

func create(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	if !contains(createRuntimes, createRuntimeOption) {
		return fmt.Errorf("the actor runtime '%s' is not supported", createRuntimeOption)
	}

	packageJson := `{
  "dependencies": {
    "express": "^4.18.1",
    "node-fetch": "^2.6.7"
  }
}
`
	fmt.Println("Writing package.json.")
	err := os.WriteFile("package.json", []byte(packageJson), 0644)
	if err != nil {
		return err
	}

	indexJs := `class Actor {
  set (v) { this.v = v; return 'OK' }
  get () { return this.v }
  ip () { return require('os').networkInterfaces().eth0[0].address }
}

// DO NOT MODIFY CODE BELOW THIS POINT

const fetch = require('node-fetch')
const actor = {
  invoke: async function (service, instance, method, ...args) {
    const fqn = service.includes('.') ? service : service + '.default'
    const resp = await fetch(` + "`" + `http://${fqn}/actor/v1/invoke/${instance}/${method}` + "`" + `, {
      method: 'post', body: JSON.stringify(args), headers: { 'Content-Type': 'application/json', 'K-Session': instance }
    })
    const body = await resp.json()
    if (body.error !== undefined) throw new Error(body.error); else return body.value
  }
}

const express = require('express')
const app = express()
const actors = {}
app.use(express.json({ type: '*/*' }))
app.post('/actor/v1/invoke/:id/:method', async (req, res) => {
  try {
    const id = req.params.id
    if (id === undefined) throw new Error('invalid actor id')
    if (!actors[id]) actors[id] = new Actor()
    if (typeof actors[id][req.params.method] !== 'function') throw new Error('undefined method')
    const value = await actors[id][req.params.method](...req.body)
    const body = JSON.stringify({ value })
    res.status(200).type('application/json').send(body)
  } catch (err) {
    const error = typeof err.message === 'string' ? err.message : typeof err === 'string' ? err : 'internal error'
    const body = JSON.stringify({ error })
    res.status(400).type('application/json').send(body)
  }
})
app.delete('/actor/v1/deactivate/:id', async (req, res) => {
  delete actors[req.params.id]
  res.status(200).send('OK')
})
process.on('SIGTERM', () => { console.log('Actor runtime exiting'); process.exit(0) })
app.listen(8080, () => console.log('Actor runtime started'))
`

	fmt.Println("Writing index.js.")
	err = os.WriteFile("index.js", []byte(indexJs), 0644)
	if err != nil {
		return err
	}

	dockerfile := `FROM node:16-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install --omit=dev
COPY index.js ./
CMD [ "node" , "index.js" ]
`

	fmt.Println("Writing Dockerfile.")
	err = os.WriteFile("Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Actor project created.")
	return nil
}
