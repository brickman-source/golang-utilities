package retry

import (
	"github.com/brickman-source/golang-utilities/log"
	"time"
)

func Run(maxWait time.Duration, failAfter time.Duration, f func() error) (err error) {
	var lastStart time.Time

	loopWait := time.Millisecond * 100
	retryStart := time.Now()
	for retryStart.Add(failAfter).After(time.Now()) {
		lastStart = time.Now()
		if err = f(); err == nil {
			return nil
		}

		if lastStart.Add(maxWait * 2).Before(time.Now()) {
			retryStart = time.Now()
		}

		log.Errorf("run err:%s", err.Error())
		log.Debugf("Retrying in %f seconds...", loopWait.Seconds())

		time.Sleep(loopWait)

		loopWait = loopWait * time.Duration(int64(2))
		if loopWait > maxWait {
			loopWait = maxWait
		}
	}
	return err
}
