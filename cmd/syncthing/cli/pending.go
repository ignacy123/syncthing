// Copyright (C) 2021 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package cli

import (
	"fmt"

	"encoding/json"

	"github.com/urfave/cli"

	"strconv"
)

var pendingCommand = cli.Command{
	Name:     "pending",
	HideHelp: true,
	Usage:    "Pending command group",
	Subcommands: []cli.Command{
		{
			Name:   "devices",
			Usage:  "Dump pending devices as json",
			Action: expects(0, indexDumpOutput("cluster/pending/devices")),
		},
		{
			Name:   "device",
			Usage:  "Print device-id of a pending device",
			ArgsUsage: "[number]",
			Action: expects(1, device),
		},
		{
			Name:   "folders",
			Usage:  "Dump pending folders as json",
			Action: expects(0, indexDumpOutput("cluster/pending/folders")),
		},
		{
			Name:   "folder-device",
			Usage:  "Dump pending folders shared by a given device as json",
			ArgsUsage: "[device-id]",
			Action: expects(1, folder_device),
		},
	},
}

func device(c *cli.Context) error {
    client, err := getClientFactory(c).getClient()
    if err != nil {
        return err
    }
    response, err := client.Get("cluster/pending/devices")
    if err != nil {
        return err
    }
	bytes, err := responseToBArray(response)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	i := 1
	number, err := strconv.Atoi(c.Args()[0])
	if err != nil{
        return err
    }
	for key := range data{
        if i == number {
            fmt.Println(key)
            return nil
        }
        i += 1
    }
    return fmt.Errorf("There is less than %d pending devices", number)
}

func folder_device(c *cli.Context) error {
    device_id := c.Args()[0]
    return indexDumpOutputWithParams("cluster/pending/folders", map[string]string{"device": device_id})(c)
}
