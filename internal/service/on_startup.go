package service

import "context"

type Startable interface {
	startup(ctx context.Context)
}

func OnStartup(ss []Startable) func (c context.Context) {
	return func (c context.Context) {
		for _, s := range ss {
			s.startup(c)
		}
	}
}