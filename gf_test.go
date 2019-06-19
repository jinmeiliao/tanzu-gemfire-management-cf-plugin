package main

import (
	"code.cloudfoundry.org/cli/cf/errors"
	cloudcachemanagementcfpluginfakes "github.com/gemfire/cloudcache-management-cf-plugin/cloudcache-management-cf-pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cf gf plugin", func() {
	Context("Retrieving Username, Password, and Endpoint", func() {
		It("Returns correct information", func() {
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			pccService := "jjack"
			key := "mykey"
			keyInfo := `Getting key mykey for service instance jjack as admin...

{
 "distributed_system_id": "0",
 "gfsh_login_string": "connect --url=https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/gemfire/v1 --user=cluster_operator_ygTWCaBfqtFHuTWxdaOMQ --password=W97ghWi4p2YF5MsfRCu6Eg --skip-ssl-validation",
 "locators": [
  "10.0.8.6[55221]"
 ],
 "urls": {
  "gfsh": "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/gemfire/v1",
  "pulse": "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/pulse"
 },
 "users": [
  {
   "password": "W97ghWi4p2YF5MsfRCu6Eg",
   "roles": [
    "cluster_operator"
   ],
   "username": "cluster_operator_ygTWCaBfqtFHuTWxdaOMQ"
  },
  {
   "password": "vcM942IBtpZrL3MxWyyi6Q",
   "roles": [
    "developer"
   ],
   "username": "developer_T2ONcuzffmoQ3Zv1HIraGQ"
  }
 ],
 "wan": {
  "sender_credentials": {
   "active": {
    "password": "nOVwPCVF25SeVfWgVTgCKA",
    "username": "gateway_sender_J3wGFvXCzO1hGH7ESx7EA"
   }
  }
 }
}
`
			expectedUsername :="cluster_operator_ygTWCaBfqtFHuTWxdaOMQ"
			expectedPassword := "W97ghWi4p2YF5MsfRCu6Eg"
			expectedEndpoint := "https://cloudcache-7fe65c41-cca5-43c2-afaa-019ef452c6a1.sys.mammothlakes.cf-app.com/geode-management/v2"
			fakeCf.CmdReturns(keyInfo, nil)
			username, password, endpoint, err := GetUsernamePasswordEndpoint(fakeCf, pccService, key)
			Expect(username).To(Equal(expectedUsername))
			Expect(password).To(Equal(expectedPassword))
			Expect(endpoint).To(Equal(expectedEndpoint))
			Expect(err).To(BeNil())
		})
		It("Returns an error.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("", errors.New("CF Command Error"))
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf, "", "")
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incorrect JSON.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			fakeCf.CmdReturns("{", nil)
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf, "", "")
			Expect(err).To(Not(BeNil()))
		})
		It("Resolving incomplete JSON.", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			pccService := "jjack"
			key := "mykey"
			keyInfo := `Getting key mykey for service instance jjack as admin...

{
 
`
			fakeCf.CmdReturns(keyInfo, nil)
			_, _, _, err := GetUsernamePasswordEndpoint(fakeCf, pccService, key)
			Expect(err).To(Not(BeNil()))
		})
	})
	Context("Retrieving Service Key from PCC Instance", func() {
		It("Returns correct information", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `Getting keys for service instance jjack as admin...

name
mykey

`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			response, err := GetServiceKeyFromPCCInstance(fakeCf, "jjack")
			Expect(err).To(BeNil())
			expectedResponse := "mykey"
			Expect(response).To(Equal(expectedResponse))
		})
		It("Handling a no service instance found", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `FAILED
Service instance jjackk not found
`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			_, err := GetServiceKeyFromPCCInstance(fakeCf, "jjack")
			Expect(err).To(Not(BeNil()))
		})
		It("Handling no service key available", func(){
			fakeCf := &cloudcachemanagementcfpluginfakes.FakeCfService{}
			resultFromCFServiceKeys := `Getting keys for service instance oowen as admin...
No service key for service instance oowen
`
			fakeCf.CmdReturns(resultFromCFServiceKeys, nil)
			_, err := GetServiceKeyFromPCCInstance(fakeCf, "oowen")
			Expect(err).To(Not(BeNil()))
		})
	})
	Context("Editing Response on Group", func() {
		It("Editing on known group", func(){
			urlResponse := `{"statusCode":"OK","result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","groups":["group3","group2"],"regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjack","type":"REPLICATE","entryCount":0,"uri":"/regions/jjack"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","groups":["group2"],"regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjoris","type":"REPLICATE","entryCount":0,"uri":"/regions/jjoris"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjens","type":"REPLICATE","entryCount":0,"uri":"/regions/jjens"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"jjinmei","type":"REPLICATE","entryCount":0,"uri":"/regions/jjinmei"}]}`
			groups := []string{"group2", "group3"}
			clusterCommand := "list-regions"
			response, _ := EditResponseOnGroup(urlResponse, groups, clusterCommand)
			expectedResponse := `{"statusCode":"OK","statusMessage":"","memberStatus":null,"result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group3","group2"],"name":"jjack","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjack"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group2"],"name":"jjoris","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjoris"}]}`
			Expect(response).To(Equal(expectedResponse))
		})
		It("Editing on unknown gorup", func(){
			urlResponse :=`{"statusCode":"OK","result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","groups":["group1"],"regionAttributes":{"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK","concurrencyChecksEnabled":true},"name":"name1","type":"REPLICATE","entryCount":-1,"uri":"/regions/name1"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","regionAttributes":{"partitionAttributes":{"redundantCopies":"1"},"evictionAttributes":{"lruHeapPercentage":{"action":"LOCAL_DESTROY"}},"dataPolicy":"PARTITION","concurrencyChecksEnabled":true},"name":"example_partition_region","type":"PARTITION_REDUNDANT_HEAP_LRU","entryCount":0,"uri":"/regions/example_partition_region"}]}`
			groups := []string{"group4", "group5"}
			clusterCommand := "list-regions"
			response, _ := EditResponseOnGroup(urlResponse, groups, clusterCommand)
			expectedResponse := `{"statusCode":"OK","statusMessage":"","memberStatus":null,"result":null}`
			Expect(response).To(Equal(expectedResponse))
		})
	})
	Context("Safekeeping tests", func(){
		It("Validate PCC instance", func(){
			ourPCCInstance := "pcc1"
			pccInstancesAvailable := []string{"pcc1", "pcc2", "pcc3"}
			err := ValidatePCCInstance(ourPCCInstance, pccInstancesAvailable)
			Expect(err).To(BeNil())
		})
		It("Validate answer in table form", func(){
			clusterCommand := "list-members"
			urlResponse := `{"statusCode":"OK","result":[{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","groups":["group3"],"id":"server3","host":"10.118.19.35","status":"online","pid":63340,"cacheServers":[{"port":40406,"maxConnections":800,"maxThreads":0}],"maxHeap":4096,"initialHeap":256,"usedHeap":55,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/server3/server3.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/server3","clientConnections":0,"locator":false,"coordinator":false,"uri":"/members/server3"},{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","groups":["group-1"],"id":"server1","host":"10.118.19.35","status":"online","pid":39383,"cacheServers":[{"port":53105,"maxConnections":800,"maxThreads":0}],"maxHeap":4096,"initialHeap":256,"usedHeap":37,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/server1/server1.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/server1","clientConnections":0,"locator":false,"coordinator":false,"uri":"/members/server1"},{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","groups":["group2"],"id":"server2","host":"10.118.19.35","status":"online","pid":62028,"cacheServers":[{"port":40404,"maxConnections":800,"maxThreads":0}],"maxHeap":4096,"initialHeap":256,"usedHeap":57,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/server2/server2.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/server2","clientConnections":0,"locator":false,"coordinator":false,"uri":"/members/server2"},{"class":"org.apache.geode.management.configuration.RuntimeMemberConfig","id":"locator1","host":"10.118.19.35","status":"online","pid":39363,"port":10334,"maxHeap":4096,"initialHeap":256,"usedHeap":132,"logFile":"/Users/jackweissburg/projects/cf-plugin-demo/locator1/locator1.log","workingDirectory":"/Users/jackweissburg/projects/cf-plugin-demo/locator1","clientConnections":0,"locator":true,"coordinator":true,"uri":"/members/locator1"}]}`
			response, _ := GetAnswerFromUrlResponse(clusterCommand, urlResponse)
			expectedResponse := `Status Code: OK

 id                 | host               | status             | pid                |
 ------------------------------------------------------------------------------------
 server3            | 10.118.19.35       | online             | 63340              |
 server1            | 10.118.19.35       | online             | 39383              |
 server2            | 10.118.19.35       | online             | 62028              |
 locator1           | 10.118.19.35       | online             | 39363              |

Number of Members: 4`
			Expect(response).To(Equal(expectedResponse))
		})
		It("Validate JSON from URL response", func(){
			urlResponse := `{"statusCode":"OK","statusMessage":"","memberStatus":null,"result":[{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group3","group2"],"name":"jjack","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjack"},{"class":"org.apache.geode.management.configuration.RuntimeRegionConfig","entryCount":0,"groups":["group2"],"name":"jjoris","regionAttributes":{"concurrencyChecksEnabled":true,"dataPolicy":"REPLICATE","scope":"DISTRIBUTED_ACK"},"type":"REPLICATE","uri":"/regions/jjoris"}]}`
			expectedResponse := `{
  "statusCode": "OK",
  "statusMessage": "",
  "memberStatus": null,
  "result": [
    {
      "class": "org.apache.geode.management.configuration.RuntimeRegionConfig",
      "entryCount": 0,
      "groups": [
        "group3",
        "group2"
      ],
      "name": "jjack",
      "regionAttributes": {
        "concurrencyChecksEnabled": true,
        "dataPolicy": "REPLICATE",
        "scope": "DISTRIBUTED_ACK"
      },
      "type": "REPLICATE",
      "uri": "/regions/jjack"
    },
    {
      "class": "org.apache.geode.management.configuration.RuntimeRegionConfig",
      "entryCount": 0,
      "groups": [
        "group2"
      ],
      "name": "jjoris",
      "regionAttributes": {
        "concurrencyChecksEnabled": true,
        "dataPolicy": "REPLICATE",
        "scope": "DISTRIBUTED_ACK"
      },
      "type": "REPLICATE",
      "uri": "/regions/jjoris"
    }
  ]
}`
			response, _:= GetJsonFromUrlResponse(urlResponse)
			Expect(response).To(Equal(expectedResponse))
		})
		It("Validate table filling", func(){
			columnSize := 20
			value := "some string"
			filler := "-"
			response := Fill(columnSize, value, filler)
			expectedResponse := " some string--------"
			Expect(response).To(Equal(expectedResponse))
		})
	})
})