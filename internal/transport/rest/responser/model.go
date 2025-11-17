package responser

import "go.uber.org/zap"

type Responser struct {
	log *zap.Logger
}

func NewResponser(log *zap.Logger) *Responser {
	return &Responser{log: log}
}
