package pcc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/implfakes"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/pcc"
)

var _ = Describe("Plugin", func() {

	var (
		cliConnection    *pluginfakes.FakeCliConnection
		pluginConnection *implfakes.FakeConnectionProvider
		commandData      domain.CommandData
		basicPlugin      BasicPlugin
		commandProcessor common.CommandProcessor
		requestHelper    *implfakes.FakeRequestHelper
	)

	BeforeEach(func() {
		cliConnection = new(pluginfakes.FakeCliConnection)
		pluginConnection = new(implfakes.FakeConnectionProvider)
		requestHelper = new(implfakes.FakeRequestHelper)

		commandProcessor, err := common.NewCommandProcessor(requestHelper)
		Expect(err).NotTo(HaveOccurred())

		basicPlugin, err = NewBasicPlugin(commandProcessor)
		Expect(err).NotTo(HaveOccurred())

		commandData = domain.CommandData{}
		commandData.UserCommand = domain.UserCommand{}
		commandData.UserCommand.Parameters = make(map[string]string)
		commandData.Target = "pcc1"
	})

	Context("Plugin has a functioning cli connection", func() {
		It("Exists immediately when arg[0] is 'CLI-MESSAGE-UNINSTALL'", func() {
			args := []string{"CLI-MESSAGE-UNINSTALL"}
			basicPlugin.Run(cliConnection, args)
		})

		It("Exits when user command is missing", func() {
			args := []string{}
			basicPlugin.Run(cliConnection, args)
		})

		It("Exits when it fails to obtain a plugin connection", func() {

			args := []string{"http://localhost:7070", "list", "members"}
			basicPlugin.Run(cliConnection, args)
		})

	})
})
