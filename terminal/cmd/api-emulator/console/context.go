package console

import "context"

// Context provides console context to handler
type Context struct {
	context.Context
	cancelContext func()
	console       *Console
	prompt        string
	activity      Activity
}

// Printf prints!
func (cc Context) Printf(f string, args ...interface{}) {
	cc.console.Printf(cc.prefix()+f, args...)
}

// Debugf prints!
func (cc Context) Debugf(f string, args ...interface{}) {
	cc.console.Debugf(cc.prefix()+f, args...)
}

// Infof prints!
func (cc Context) Infof(f string, args ...interface{}) {
	cc.console.Infof(cc.prefix()+f, args...)
}

// Warningf prints!
func (cc Context) Warningf(f string, args ...interface{}) {
	cc.console.Warningf(cc.prefix()+f, args...)
}

// Errorf prints!
func (cc Context) Errorf(f string, args ...interface{}) {
	cc.console.Errorf(cc.prefix()+f, args...)
}

// Enter console
func (cc Context) Enter(ctx context.Context, h Activity, prompt string) bool {
	return cc.console.push(ctx, h, prompt)
}

// Leave console
func (cc Context) Leave() {
	cc.console.pop()
}

func (cc Context) prefix() string {
	// if cc.prompt != "" {
	// 	return cc.prompt + ": "
	// }
	return ""
}
