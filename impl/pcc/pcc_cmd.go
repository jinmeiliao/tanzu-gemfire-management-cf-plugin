package pcc

import "github.com/gemfire/cloudcache-management-cf-plugin/domain"

import "github.com/gemfire/cloudcache-management-cf-plugin/impl/common"

// Command is the receiver for the commands for PCC
type Command struct {
	commandData domain.CommandData
	comm        common.CommandProcessor
	connection  PluginConnection
}

// NewCommand constructor for PccCommand
func NewCommand(comm common.CommandProcessor, connection PluginConnection) (Command, error) {
	return Command{comm: comm, connection: connection}, nil
}

// Run is the main entry point for the standalone Geode command line interface
// It is run once for each command executed
func (pc *Command) Run(args []string) (err error) {
	err = nil

	return
}
