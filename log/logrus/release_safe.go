// +build safe

package logrus

import (
	"github.com/sirupsen/logrus"
)

func releaseEntry(log *logrus.Logger, ent *logrus.Entry) {}
