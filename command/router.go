// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bitbucket.org/jonforums/uru/env"
)

type HandlerFunc func(*env.Context)

type CommandRouter struct {
	handlers   map[string]HandlerFunc
	defHandler HandlerFunc
}

// Returns a newly configured, ready-to-use command router. Provide a default
// handler function that will be called when no predefined commands can be
// matched against the user specified command string. Most often used when
// creating a top-level command router in which arbitrary tokens are used to
// activate a particular ruby runtime.
//
// If the default handler function is nil, the function will never be called.
func NewRouter(handler HandlerFunc) *CommandRouter {

	return &CommandRouter{
		handlers:   make(map[string]HandlerFunc),
		defHandler: handler,
	}
}

// Handle binds a function to a set of user CLI command alias strings. The bound
// function is be executed whenever a user specifies one of the command aliases.
func (r *CommandRouter) Handle(cmds []string, handler HandlerFunc) {
	for _, c := range cmds {
		r.handlers[c] = handler
	}
}

// Dispatch calls a previously bound function corresponding to the user
// specified command string, passing a context as the only arg. If the
// command string is not a recognized command, and the CommandRouter instance
// has been created with a non-nil default handler, the default handler will
// be invoked with a context as the only arg.
func (r *CommandRouter) Dispatch(ctx *env.Context, cmd string) {
	if f, ok := r.handlers[cmd]; ok {
		f(ctx)
	} else {
		if r.defHandler != nil {
			r.defHandler(ctx)
		}
	}
}
