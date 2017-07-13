package commands

import "github.com/earlcherry/gouter"

type Command interface {
	Run(args gouter.RouteArgs) error
}
