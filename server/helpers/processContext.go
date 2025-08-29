package helpers

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"server/config"
)

var pc *ProcessContext

type ProcessContext struct {
	done           chan bool
	interrupt      chan os.Signal
	FatalErrorChan chan error

	CTX    context.Context
	cancel context.CancelFunc

	interruptOnce *sync.Once

	WG *sync.WaitGroup
}

func GetProcessContext() context.Context {
	return pc.CTX
}

func GetMainProcess() *ProcessContext {
	return pc
}

func InitProcessContext() {
	pc = new(ProcessContext)
	pc.CTX, pc.cancel = context.WithCancel(context.Background())
	pc.WG = new(sync.WaitGroup)
	pc.done = make(chan bool)
	pc.FatalErrorChan = make(chan error)
	pc.interrupt = make(chan os.Signal)
	pc.interruptOnce = new(sync.Once)
}

func (ctx *ProcessContext) AddWorker(n int) {
	ctx.WG.Add(n)
}

func (ctx *ProcessContext) CompleteOneWorker() {
	ctx.WG.Done()
}

func (ctx *ProcessContext) WaitForFinish() {
	go ctx.handleInterrupt()
	go ctx.waitGroupDone()
	go ctx.watchError()
	ctx.gracefullyExit()
}

func (ctx *ProcessContext) handleInterrupt() {
	// system inturrpt signal or terminate signal will be passed on interrupt channel
	signal.Notify(ctx.interrupt, syscall.SIGINT, syscall.SIGTERM)

	for range ctx.interrupt {
		ctx.interruptOnce.Do(func() {
			ctx.cancel()
			config.GetLogger().Println("First interrupt signal captured.")
		})
	}
}

func (ctx *ProcessContext) watchError() {
	err := <-ctx.FatalErrorChan
	config.GetLogger().Println("Fatal error captured :: ", err)
	ctx.cancel()
}

func (ctx *ProcessContext) waitGroupDone() {
	ctx.WG.Wait()
	ctx.done <- true
}

func (ctx *ProcessContext) gracefullyExit() {
	<-ctx.done
	config.GetLogger().Println("Gracefully exit")
	os.Exit(0)
}

func (ctx *ProcessContext) SetContextValue(key config.ContextKey, value any) {
	ctx.CTX = context.WithValue(ctx.CTX, key, value)
}

func (ctx *ProcessContext) GetContextValue(key config.ContextKey) any {
	return ctx.CTX.Value(key)
}
