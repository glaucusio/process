package process

import (
	"context"
	"os"
	"os/signal"
)

var g = newSiqueue()

func Context() context.Context {
	return g.with(context.Background())
}

type siqueue struct {
	cancel  chan context.CancelFunc
	signals []os.Signal
}

func newSiqueue() siqueue {
	q := siqueue{
		cancel:  make(chan context.CancelFunc),
		signals: defaultSignals,
	}
	go q.process()
	return q
}

func (q siqueue) with(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	q.cancel <- cancel
	return ctx
}

func (q siqueue) process() {
	var queue []context.CancelFunc
	var sig = make(chan os.Signal, 1)

	signal.Notify(sig, q.signals...)

	for {
		select {
		case cancel, ok := <-q.cancel:
			if !ok {
				return
			}

			queue = append(queue, cancel)
		case <-sig:
			if len(queue) == 0 {
				close(q.cancel)
				return
			}

			var cancel func()
			cancel, queue = queue[0], queue[1:]
			cancel()
		}
	}
}
