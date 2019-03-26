package golog

import "sync"

const (
	ASYNC_MSG_KIND_WRITE = 1
	ASYNC_MSG_KIND_FLUSH = 2
	ASYNC_MSG_KIND_FREE  = 3
)

type asyncMsg struct {
	kind int
	msg  []byte
}

type asyncWriter struct {
	w IWriter

	msgCh chan *asyncMsg
	wg    *sync.WaitGroup
}

func NewAsyncWriter(w IWriter, queueSize int) *asyncWriter {
	a := &asyncWriter{
		w: w,

		msgCh: make(chan *asyncMsg, queueSize),
		wg:    new(sync.WaitGroup),
	}

	go a.asyncLogRoutine()
	a.wg.Add(1)

	return a
}

func (a *asyncWriter) asyncLogRoutine() {
	defer a.wg.Done()

	for {
		am := <-a.msgCh
		switch am.kind {
		case ASYNC_MSG_KIND_WRITE:
			a.w.Write(am.msg)
		case ASYNC_MSG_KIND_FLUSH:
			a.w.Flush()
		case ASYNC_MSG_KIND_FREE:
			a.w.Free()
			return
		}
	}
}

func (a *asyncWriter) Write(p []byte) (int, error) {
	a.msgCh <- &asyncMsg{ASYNC_MSG_KIND_WRITE, p}

	return len(p), nil
}

func (a *asyncWriter) Flush() error {
	a.msgCh <- &asyncMsg{ASYNC_MSG_KIND_FLUSH, nil}

	return nil
}

func (a *asyncWriter) Free() {
	a.msgCh <- &asyncMsg{ASYNC_MSG_KIND_FREE, nil}

	a.wg.Wait()
}
