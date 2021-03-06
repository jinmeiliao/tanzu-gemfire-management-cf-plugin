/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for additional information regarding
 * copyright ownership. The ASF licenses this file to You under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License. You may obtain a
 * copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package format_test

import (
	"strings"

	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common"
	"github.com/gemfire/cloudcache-management-cf-plugin/impl/common/filter"
	. "github.com/gemfire/cloudcache-management-cf-plugin/impl/common/format"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting", func() {

	Describe("Fill tests", func() {
		It("Replaces invalid filler with 'space' characters", func() {
			columnSize := 5
			value := "id"
			filler := "test"
			response := Fill(columnSize, value, filler)
			expectedResponse := " id  "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})

		It("small column size", func() {
			columnSize := 4
			value := "id"
			filler := " "
			response := Fill(columnSize, value, filler)
			expectedResponse := " id "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})

		It("small column size", func() {
			columnSize := 3
			value := "i"
			filler := " "
			response := Fill(columnSize, value, filler)
			expectedResponse := " i "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})

		It("small column size", func() {
			columnSize := 2
			value := "idle"
			filler := " "
			response := Fill(columnSize, value, filler)
			expectedResponse := " ... "
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(5))
		})

		It("Fills the table with filler characters", func() {
			columnSize := 20
			value := "some string"
			filler := "-"
			response := Fill(columnSize, value, filler)
			expectedResponse := "-some string--------"
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})

		It("Truncates the value and adds Ellipsis at the end of the value", func() {
			columnSize := 20
			value := "some strings that is longer than 20 characters"
			filler := "-"
			response := Fill(columnSize, value, filler)
			expectedResponse := "-some strings th...-"
			Expect(response).To(Equal(expectedResponse))
			Expect(len(response)).To(Equal(columnSize))
		})
	})

	Context("FormatResponse tests", func() {
		var (
			formatter common.Formatter
			err       error
		)

		BeforeEach(func() {
			formatter, err = New(filter.GOJQFilter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Returns the input as an indented string", func() {
			inputString := `{"name": "value"}`
			expectedString := "{\n  \"name\": \"value\"\n}"
			output, err := formatter.FormatResponse(inputString, "", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(expectedString))
		})

		It("Returns the input 'as-is'", func() {
			inputString := "foobar"
			output, err := formatter.FormatResponse(inputString, "", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(inputString))
		})

		It("Returns an error when faulty json query string is used", func() {
			inputString := `[{"name": "value"}]`
			output, err := formatter.FormatResponse(inputString, ".[], | {name:.name}", true)
			Expect(output).To(BeEmpty())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("json query failed"))
		})

		It("Returns filtered input with correct userFilter appended", func() {
			inputString := `[{"name": "value"}]`
			output, err := formatter.FormatResponse(inputString, ".[] | {name:.name}", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(" name  \n-------\n value \n\nJQFilter: .[] | {name:.name}\n"))
		})

		It("Returns no results with userFilter that yields empty array", func() {
			inputString := `{"result": []}`
			filter := `.result[]`
			output, err := formatter.FormatResponse(inputString, filter, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal("\nJQFilter: " + filter + "\n"))
		})

		It("Returns a result with default filter that yields empty array", func() {
			inputString := `{"result": []}`
			filter := `.result[]`
			output, err := formatter.FormatResponse(inputString, filter, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(" result \n--------\n []     \n\nJQFilter: .\n"))
		})

		It("Returns all results with default . filter", func() {
			inputString := `{"name": "value"}`
			output, err := formatter.FormatResponse(inputString, ".", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(" name  \n-------\n value \n\nJQFilter: .\n"))
		})

		It("Returns list results with filter", func() {
			inputString := `{"result": [{"name": "value1"}, {"name":"value2"}, {"name":"value3"}, {"name":"value4"}]}`
			filter := `.result[]`
			output, err := formatter.FormatResponse(inputString, filter, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(Equal(` name   
--------
 value1 
 value2 
 value3 
 value4 

JQFilter: .result[]
`))
		})
	})

	Context("tabular output", func() {
		It("Returns the input as table format", func() {
			json := `[{
				"id": "server",
				"status": "online"},
			  {"id": "locator",
				"status": "online"}]`
			output, _ := Tabular(json)
			splitOutput := strings.Split(output, "\n")
			Expect(splitOutput[0]).To(ContainSubstring("id"))
			Expect(splitOutput[0]).To(ContainSubstring("status"))
			Expect(splitOutput[1]).To(Equal("------------------"))
			Expect(splitOutput[2] + splitOutput[3]).To(ContainSubstring("server"))
			Expect(splitOutput[2] + splitOutput[3]).To(ContainSubstring("online"))
		})
		It("different attributes", func() {
			json := `[{"id": "server"},{"status": "online"}]`
			output, _ := Tabular(json)
			splitOutput := strings.Split(output, "\n")
			Expect(splitOutput[0]).To(ContainSubstring("id"))
			Expect(splitOutput[0]).To(ContainSubstring("status"))
			Expect(splitOutput[1]).To(Equal("-----------------"))
			Expect(splitOutput[2] + splitOutput[3]).To(ContainSubstring("server"))
			Expect(splitOutput[2] + splitOutput[3]).To(ContainSubstring("online"))
		})

		It("empty json array", func() {
			output, _ := Tabular("[]")
			Expect(output).To(Equal(""))
		})
		It("invalid json string", func() {
			_, err := Tabular(`{"name":"test"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to parse:"))
		})
	})

	Context("DescribeEndpoint", func() {
		var (
			endPoint  domain.RestEndPoint
			formatter common.Formatter
			err       error
		)

		BeforeEach(func() {
			formatter, err = New(filter.GOJQFilter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Shows the expected command name in the output", func() {
			endPoint = domain.RestEndPoint{}
			endPoint.CommandName = "testcommand"

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("testcommand"))
		})

		It("Shows all expected parameters in the output", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramTwo := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramTwo.Name = "paramTwo"
			endPoint.Parameters = []domain.RestAPIParam{paramOne, paramTwo}

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("--paramOne"))
			Expect(result).To(ContainSubstring("--paramTwo"))
		})

		It("Shows if parameters are optional in the output", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramOne.Description = "first parameter"
			paramOne.Required = false
			endPoint.Parameters = []domain.RestAPIParam{paramOne}

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("[--paramOne <first parameter>]"))
		})

		It("Shows if parameters are required in the output", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramOne.Required = true
			endPoint.Parameters = []domain.RestAPIParam{paramOne}

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(BeEmpty())
			Expect(result).NotTo(ContainSubstring("["))
			Expect(result).NotTo(ContainSubstring("]"))
			Expect(result).To(ContainSubstring("--paramOne"))
		})

		It("Shows GeneralOptions when showDetails flag set to true", func() {
			endPoint = domain.RestEndPoint{}

			result := formatter.DescribeEndpoint(endPoint, true)
			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring(GeneralOptions))
		})

		It("Hides GeneralOptions when showDetails flag set to false", func() {
			endPoint = domain.RestEndPoint{}

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).NotTo(ContainSubstring(GeneralOptions))
		})

		It("Correctly display expected body format if 'body' parameter present and showDetails flag set to true", func() {
			endPoint = domain.RestEndPoint{}
			paramOne := domain.RestAPIParam{}
			paramOne.Name = "paramOne"
			paramOne.Required = true
			paramOne.In = "body"
			bodyDefinition := make(map[string]interface{})

			bodyDefinition["propDetail1"] = "string-value"
			bodyDefinition["propDetail2"] = 42
			bodyDefinition["propDetail3"] = []int{21, 22}
			bodyDefinition["propDetail4"] = true
			bodyDefinition["propDetail5"] = map[string]interface{}{"1someString": "stringValue", "2someNumber": 23, "3someBool": false}

			paramOne.BodyDefinition = bodyDefinition
			endPoint.Parameters = []domain.RestAPIParam{paramOne}
			expectedOutput := `--paramOne format:
		{
		  "propDetail1": "string-value",
		  "propDetail2": 42,
		  "propDetail3": [
		    21,
		    22
		  ],
		  "propDetail4": true,
		  "propDetail5": {
		    "1someString": "stringValue",
		    "2someNumber": 23,
		    "3someBool": false
		  }
		}`

			result := formatter.DescribeEndpoint(endPoint, true)
			Expect(result).To(ContainSubstring(expectedOutput))
			Expect(result).To(ContainSubstring(GeneralOptions))
		})

		It("describe the rest end point without body param", func() {
			var endPoint domain.RestEndPoint
			endPoint.CommandName = "test"
			endPoint.Parameters = make([]domain.RestAPIParam, 2)

			var param1, param2 domain.RestAPIParam
			param1.In = "query"
			param1.Name = "id"
			param1.Description = "id"
			param1.Required = true

			param2.In = "query"
			param2.Name = "group"
			param2.Description = "group"
			param2.Required = false
			endPoint.Parameters[0] = param2
			endPoint.Parameters[1] = param1

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).To(Equal("test --id <id> [--group <group>]"))
		})

		It("describe the rest end point with body param", func() {
			var endPoint domain.RestEndPoint
			endPoint.CommandName = "test"
			endPoint.Parameters = make([]domain.RestAPIParam, 2)

			var param1, param2 domain.RestAPIParam
			param1.In = "body"
			param1.Name = "config"
			param1.Description = "config"
			param1.Required = true

			param2.In = "query"
			param2.Name = "group"
			param2.Description = "group"
			param2.Required = false
			endPoint.Parameters[0] = param2
			endPoint.Parameters[1] = param1

			result := formatter.DescribeEndpoint(endPoint, false)
			Expect(result).To(Equal("test --config <json or @json_file_path> [--group <group>]"))
		})
	})
})
