// Copyright © 2016 EMC Corporation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codedellemc/gorackhd/client/skus"
	"github.com/codedellemc/gorackhd/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// nodeslistCmd represents the nodeslist command
var nodeslistCmd = &cobra.Command{
	Use:   "list",
	Short: "List Nodes in RackHD",
	Long:  "List Nodes in RackHD",
	Run:   listNodes,
}

var nodeSku string
var shortList bool
var withtags string
var withouttags string

func init() {
	nodesCmd.AddCommand(nodeslistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeslistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	nodeslistCmd.Flags().StringVar(&nodeSku, "sku", "", "SKU id")
	nodeslistCmd.Flags().BoolVarP(&shortList, "quiet", "q", false, "list only Node IDs")
	nodeslistCmd.Flags().StringVar(&withtags, "with-tags", "", "only show nodes that have at least ONE of the given tags (comma separated)")
	nodeslistCmd.Flags().StringVar(&withouttags, "without-tags", "", "only show nodes that do not have ANY of the given tags (comma separated)")
}

func listNodes(cmd *cobra.Command, args []string) {
	var payload []*models.Node
	if nodeSku != "" {
		skuParams := skus.GetSkusIdentifierNodesParams{}
		skuParams.WithIdentifier(nodeSku)
		resp, err := clients.rackMonorailClient.Skus.GetSkusIdentifierNodes(&skuParams, nil)
		if err != nil {
			log.Fatal(err)
		}
		payload = resp.Payload
	} else {
		resp, err := clients.rackMonorailClient.Nodes.GetNodes(nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		payload = resp.Payload
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID", "type", "SKU", "Tags"})

	withTagsSlice := strings.Split(withtags, ",")
	withoutTagsSlice := strings.Split(withouttags, ",")

	for _, node := range payload {
		tags := getTags(&node.Tags)
		if withtags != "" {
			found := false
			for _, a_tag := range tags {
				if stringInSlice(a_tag, withTagsSlice) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if withouttags != "" {
			found := false
			for _, a_tag := range tags {
				if stringInSlice(a_tag, withoutTagsSlice) {
					found = true
					break
				}
			}
			if found {
				continue
			}
		}
		table.Append([]string{*(node.Name), node.ID, node.Type, node.Sku, strings.Join(tags, ",")})
		if shortList {
			fmt.Println(node.ID)
		}
		//fmt.Printf("%s %s %s\n", *(n.Name), n.ID, n.Type)
		//fmt.Printf("%#v\n\n", node)
	}
	if !shortList {
		table.Render()
	}
}

func getTags(input *[]interface{}) []string {
	tags := make([]string, len(*input))
	for i, tag := range *input {
		tags[i] = tag.(string)
	}
	return tags
}
