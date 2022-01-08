// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package cmdutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/syncthing/syncthing/lib/locations"
)

type SourceDestinationPair struct {
	source      string
	destination string
}

func CheckPathArgumentsConflict(homeDir, confDir, dataDir string) error {
	homeSet := homeDir != ""
	confSet := confDir != ""
	dataSet := dataDir != ""
	if dataSet != confSet {
		return errors.New("either both or none of --config and --data must be given, use --home to set both at once")
	}
	if homeSet && dataSet {
		return errors.New("--home must not be used together with --config and --data")
	}
	return nil
}

func OverrideCredentials(homeDirSrc, confDirSrc, dataDirSrc string) error {
	if err := CheckPathArgumentsConflict(homeDirSrc, confDirSrc, dataDirSrc); err != nil {
		return err
	}
	if homeDirSrc != "" {
		confDirSrc = homeDirSrc
		dataDirSrc = homeDirSrc
	}
	if dataDirSrc != "" {
		confPath := locations.Get(locations.ConfigFile)
		certPath := locations.Get(locations.CertFile)
		keyPath := locations.Get(locations.KeyFile)
		srcDestPairs := [3]SourceDestinationPair{
			SourceDestinationPair{filepath.Join(dataDirSrc, filepath.Base(keyPath)), keyPath},
			SourceDestinationPair{filepath.Join(dataDirSrc, filepath.Base(certPath)), certPath},
			SourceDestinationPair{filepath.Join(confDirSrc, filepath.Base(confPath)), confPath},
		}
		for _, p := range srcDestPairs {
			source, err := os.Open(p.source)
			if err != nil {
				return err
			}
			defer source.Close()
			destination, err := os.OpenFile(p.destination, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return err
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SetConfigDataLocationsFromFlags(homeDir, confDir, dataDir string) error {
	if err := CheckPathArgumentsConflict(homeDir, confDir, dataDir); err != nil {
		return err
	}
	if homeDir != "" {
		confDir = homeDir
		dataDir = homeDir
	}
	if dataDir != "" {
		if err := locations.SetBaseDir(locations.ConfigBaseDir, confDir); err != nil {
			return err
		}
		return locations.SetBaseDir(locations.DataBaseDir, dataDir)
	}
	return nil
}
