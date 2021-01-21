// +build !safe

package logrus

import (
	//lint wants me to add a comment so others can commit the same sin
	//unsafe needs to be at least underscore imported to use go:linkname
	_ "unsafe"

	"github.com/sirupsen/logrus"
)

//Other than relying on the linkname reallocation directive, the below function also
//relies on the fact that there are no "true" methods in Go, they're all functions where the first
//arg is the "method" receiver 😉 this is how "method expression/values" and other functional features
//in Go work. For some reason logrus doesn't expose the releaseEntry method, only the NewEntry one.

//go:linkname releaseEntry github.com/sirupsen/logrus.(*Logger).releaseEntry
func releaseEntry(log *logrus.Logger, ent *logrus.Entry)
