//
// diskimage - handles ubuntu disk images
//
// Copyright (c) 2013 Canonical Ltd.
//
// Written by Sergio Schvezov <sergio.schvezov@canonical.com>
//
package diskimage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License version 3, as published
// by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranties of
// MERCHANTABILITY, SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR
// PURPOSE.  See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

var debugPrint bool

func init() {
	if debug := os.Getenv("DEBUG_DISK"); debug != "" {
		debugPrint = true
	}
}

type Image interface {
	Mount() error
	Unmount() error
	Format() error
	Partition() error
	Map() error
	Unmap() error
	BaseMount() string
}

type SystemImage interface {
	System() string
	Writable() string
}

type CoreImage interface {
	Image
	SystemImage
	SetupBoot(OemDescription) error
	FlashExtra(string) error
}

type HardwareDescription struct {
	Kernel          string `yaml:"kernel"`
	Dtbs            string `yaml:"dtbs"`
	Initrd          string `yaml:"initrd"`
	PartitionLayout string `yaml:"partition-layout,omitempty"`
	Bootloader      string `yaml:"bootloader"`
}

type OemDescription struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Hardware struct {
		Dtb string `yaml:"dtb,omitempty"`
	} `yaml:"hardware,omitempty"`
}

func (o OemDescription) InstallPath() string {
	return filepath.Join("/oem", o.Name, o.Version)
}

func sectorSize(dev string) (string, error) {
	out, err := exec.Command("blockdev", "--getss", dev).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to determine block size: %s", out)
	}

	return strings.TrimSpace(string(out)), err
}

func printOut(args ...interface{}) {
	if debugPrint {
		fmt.Println(args...)
	}
}
