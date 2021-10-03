package config

import (
	"io/fs"
)

type ConfigLoadRoots struct {
	CWDRoot fs.FS
	Home    fs.FS
}
