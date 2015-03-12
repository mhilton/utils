// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package shell

import (
	"os"
	"runtime"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/utils"
	"github.com/juju/utils/filepath"
)

// A PathRenderer generates paths that are appropriate for a given
// shell environment.
type PathRenderer interface {
	filepath.Renderer

	// ShQuote generates a new string with quotation marks and relevant
	// escape/control characters properly escaped. The resulting string
	// is wrapped in quotation marks such that it will be treated as a
	// single string by the shell.
	ShQuote(str string) string

	// ExeSuffix returns the filename suffix for executable files.
	ExeSuffix() string
}

type chmodder interface {
	// Chmod returns a shell command that sets the given file's
	// permissions. The result is equivalent to os.Chmod.
	Chmod(path string, perm os.FileMode) []string
}

type fileWriter interface {
	// WriteFile returns a shell command that writes the provided
	// content to a file. The command is functionally equivalent to
	// ioutil.WriteFile with permissions from the current umask.
	WriteFile(filename string, data []byte) []string
}

// Commands provides methods that may be used to generate shell
// commands for a variety of shell and filesystem operations.
type Commands interface {
	chmodder
	fileWriter

	// Mkdir returns a shell command for creating a directory. The
	// command is functionally equivalent to os.MkDir using permissions
	// appropriate for a directory.
	Mkdir(dirname string) []string

	// MkdirAll returns a shell command for creating a directory and
	// all missing parent directories. The command is functionally
	// equivalent to os.MkDirAll using permissions appropriate for
	// a directory.
	MkdirAll(dirname string) []string
}

// Renderer provides all the functionality needed to generate shell-
// compatible paths and commands.
type Renderer interface {
	PathRenderer
	Commands
}

// NewRenderer returns a Renderer for the given os.
func NewRenderer(os string) (Renderer, error) {
	if os == "" {
		os = runtime.GOOS
	}

	os = strings.ToLower(os)
	switch {
	case os == "windows":
		return &WindowsRenderer{}, nil
	case utils.OSIsUnix(os):
		return &UnixRenderer{}, nil
	case os == "ubuntu":
		return &UnixRenderer{}, nil
	default:
		return nil, errors.NotFoundf("renderer for %q", os)
	}
}
