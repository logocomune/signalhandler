package signalhandler

import (
	"errors"
	"os"
	"os/signal"
)

type NotifyHandler interface {
	Register(notifyFnc NotifyFnc, signals ...os.Signal) error
	UnRegister(signals ...os.Signal)
}

type NotifyFnc func(sig os.Signal)

type signalListMap map[os.Signal]chan os.Signal

type sHandle struct {
	signalList map[os.Signal]chan os.Signal
}

type sProperty struct {
	channel chan os.Signal
}

// NewNotifyHandler get a new handler
func NewNotifyHandler() NotifyHandler {
	s := sHandle{
		signalList: make(signalListMap),
	}
	return &s
}

// Register a new os.Signal
func (s *sHandle) Register(notifyFnc NotifyFnc, signals ...os.Signal) error {
	if len(signals) == 0 {
		return errors.New("At least a signal must be passed")
	}
	signal.Reset(signals...)
	for _, sgn := range signals {
		sigChannel := make(chan os.Signal, 1)

		cleanHandlerIfExists(sgn, s)

		s.signalList[sgn] = sigChannel

		signal.Notify(sigChannel, sgn)

		go func() {
			for s := range sigChannel {
				notifyFnc(s)
			}
		}()
	}
	return nil
}

// UnRegister remove handler
// If no signals are provided, all signal handlers will be reset.
func (s *sHandle) UnRegister(signals ...os.Signal) {
	signal.Reset(signals...)

	if len(signals) == 0 {
		for sgn := range s.signalList {
			cleanHandlerIfExists(sgn, s)
		}
		return
	}
	for _, sgn := range signals {
		cleanHandlerIfExists(sgn, s)
	}
}

func cleanHandlerIfExists(key os.Signal, sH *sHandle) {
	if v, ok := sH.signalList[key]; ok {

		close(v)
		delete(sH.signalList, key)
	}

}
