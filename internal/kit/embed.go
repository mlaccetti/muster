package kit

import (
	"io/fs"

	rootkit "github.com/mlaccetti/muster/kit"
)

func SourceFS() fs.FS {
	return rootkit.Source
}
