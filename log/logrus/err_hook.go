package logrus

import (
	"github.com/sirupsen/logrus"

	"github.com/secureworks/taegis-sdk-go/log"
	"github.com/secureworks/taegis-sdk-go/log/internal/common"
)

type errHook struct{}

func (errHook) Levels() []logrus.Level { return logrus.AllLevels }

func (errHook) Fire(event *logrus.Entry) (err error) {
	//already has a stack?
	if _, ok := event.Data[log.StackField]; ok {
		return
	}

	//doesn't have a StackTracer?
	st, ok := event.Data[logrus.ErrorKey].(common.StackTracer)
	if !ok {
		return
	}

	event.Data[log.StackField] = st.StackTrace()
	return
}
