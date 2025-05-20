package log

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}
