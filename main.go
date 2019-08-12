package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/gemfire/cloudcache-management-cf-plugin/cfservice"
	"os"
	"strings"
)


var username, password, endpoint, pccInUse, clusterCommand, serviceKey, region, jsonFile, group, id string
var hasGroup, isJSONOutput, usingPcc = false, false, false

var APICallStruct RestAPICall
var firstResponse SwaggerInfo
var availableEndpoints []IndividualEndpoint
var indivEndpoint IndividualEndpoint


func main() {
	plugin.Start(new(BasicPlugin))
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	cfClient := &cfservice.Cf{}
	if args[0] == "CLI-MESSAGE-UNINSTALL"{
		return
	}
	var err error
	err = getPCCInUseAndClusterCommand(args)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// first get credentials from environment
	username = os.Getenv("CFLOGIN")
	password = os.Getenv("CFPASSWORD")
	usingPcc = !(strings.Contains(pccInUse, "http://") || strings.Contains(pccInUse, "https://"))
	if !usingPcc {
		endpoint = pccInUse + "/management/experimental/api-docs"
	} else{
		// at this point, we have a valid clusterCommand
		serviceKey, err =  GetServiceKeyFromPCCInstance(cfClient)
		if err != nil{
			fmt.Printf(err.Error(), pccInUse, pccInUse)
			os.Exit(1)
		}

		serviceKeyUser, serviceKeyPswd, url, err := GetUsernamePasswordEndpoinFromServiceKey(cfClient)
		if err != nil{
			fmt.Println(GenericErrorMessage, err.Error())
			os.Exit(1)
		}

		endpoint = url

		// then get the credentials from the serviceKey
		if serviceKeyUser != "" && serviceKeyPswd != "" {
			username = serviceKeyUser
			password = serviceKeyPswd
		}
	}

	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	APICallStruct.parameters = make(map[string]string)

	// lastly get the credentials from the command line
	err = parseArguments(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if usingPcc {
		if username == "" && password == "" {
			fmt.Printf(NeedToProvideUsernamePassWordMessage, pccInUse, clusterCommand)
			os.Exit(1)
		} else if username != "" && password == "" {
			err = errors.New(ProvidedUsernameAndNotPassword)
			fmt.Printf(err.Error(), pccInUse, clusterCommand, username)
			os.Exit(1)
		} else if username == "" && password != "" {
			err = errors.New(ProvidedPasswordAndNotUsername)
			fmt.Printf(err.Error(), pccInUse, clusterCommand, password)
			os.Exit(1)
		}
	}


	err = executeFirstRequest()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	indivEndpoint, err = mapUserInputToAvailableEndpoint()
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = hasIDifNeeded()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = hasRegionIfNeeded()
	if err != nil {
		fmt.Printf(err.Error(), pccInUse)
		os.Exit(1)
	}
	urlResponse, err := executeSecondRequest()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jsonToBePrinted, err := GetJsonFromUrlResponse(urlResponse)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Println(jsonToBePrinted)

	//table, err := GetTableFromUrlResponse(APICallStruct.command, urlResponse)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	os.Exit(1)
	//}
	//fmt.Println(table)
	return
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "PCC_InDev",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "pcc",
				HelpText: "Commands to interact with geode cluster.\n",
				UsageDetails: plugin.Usage{
					Usage: "	cf  pcc  <*pcc_instance>  <action>  <data_type>  [*options]  (* = optional)\n\n" +
						"Supported commands:	" +
						printAvailableCommands()+
						"\nNote: pcc_instance can be saved at [$CFPCC], then omit <pcc_instance> from command ",
					Options: map[string]string{
						"h" : "this help screen\n",
						"u" : "followed by equals username (-u=<your_username>) [$CFLOGIN]\n",
						"p" : "followed by equals password (-p=<your_password>) [$CFPASSWORD]\n",
						"r" : "followed by equals region (-r=<your_region>)\n",
						"id" : "followed by an identifier required for any get command\n",
						"d" : "followed by @<json_file_path> OR single quoted JSON input \n" +
							"	     JSON required for creating/post commands\n",
					},
				},
			},
		},
	}
}
