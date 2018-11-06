package app

import "github.com/stakater/ProxyInjector/internal/pkg/cmd"

// Run runs the command
func Run() error {
	cmd := cmd.NewProxyInjectorCommand()
	return cmd.Execute()
}
