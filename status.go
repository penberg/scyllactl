package main

import (
	"encoding/json"
	"fmt"
	"github.com/mkideal/cli"
	"net/http"
	"strings"
)

var _ = app.Register(&cli.Command{
	Name: "status",
	Desc: "Print cluster information (state, load, IDs, ...)",
	Argv: func() interface{} { return new(statusT) },
	Fn:   status,
})

type statusT struct {
	cli.Helper
	Host string `cli:"host" usage:"Node hostname or IP address" dft:"localhost"`
	Port string `cli:"port" usage:"API server port number" dft:"10000"`
}

type FloatEntry struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type StringEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Node struct {
	Status  string
	State   string
	Address string
	Load    string
	Tokens  string
	Owns    string
	HostID  string
	Rack    string
}

type Datacenter struct {
	Name  string
	Nodes []Node
}

func (dc *Datacenter) AddNode(node Node) {
	dc.Nodes = append(dc.Nodes, node)
}

func getJson(baseURL string, url string, target interface{}) error {
	r, err := http.Get(baseURL + url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func status(ctx *cli.Context) error {
	argv := ctx.Argv().(*statusT)
	baseURL := fmt.Sprintf("http://%s:%s", argv.Host, argv.Port)
	loadMap := new([]FloatEntry)
	if err := getJson(baseURL, "/storage_service/load_map", loadMap); err != nil {
		return err
	}
	ownershipMap := new([]StringEntry)
	if err := getJson(baseURL, "/storage_service/ownership", ownershipMap); err != nil {
		return err
	}
	hostIDMap := new([]StringEntry)
	if err := getJson(baseURL, "/storage_service/host_id", hostIDMap); err != nil {
		return err
	}
	datacenters := make(map[string]Datacenter)
	liveNodes := new([]string)
	if err := getJson(baseURL, "/gossiper/endpoint/live", liveNodes); err != nil {
		return err
	}
	for _, nodeAddr := range *liveNodes {
		datacenterName := new(string)
		if err := getJson(baseURL, "/snitch/datacenter?host="+nodeAddr, datacenterName); err != nil {
			return err
		}
		rack := new(string)
		if err := getJson(baseURL, "/snitch/rack?host="+nodeAddr, rack); err != nil {
			return err
		}
		status := "U"
		state := "N"
		load := "?"
		for _, entry := range *loadMap {
			if entry.Key == nodeAddr {
				load = fmt.Sprintf("%f", entry.Value)
				break
			}
		}
		ownership := "?"
		for _, entry := range *ownershipMap {
			if entry.Key == nodeAddr {
				ownership = entry.Value
				break
			}
		}
		tokens := new([]string)
		if err := getJson(baseURL, "/storage_service/tokens/"+nodeAddr, tokens); err != nil {
			return err
		}
		tokenCount := fmt.Sprintf("%d", len(*tokens))
		hostID := "?"
		for _, entry := range *hostIDMap {
			if entry.Key == nodeAddr {
				hostID = entry.Value
				break
			}
		}
		node := Node{
			Status:  status,
			State:   state,
			Address: nodeAddr,
			Load:    load,
			Tokens:  tokenCount,
			Owns:    ownership,
			HostID:  hostID,
			Rack:    *rack,
		}
		datacenter, ok := datacenters[*datacenterName]
		if !ok {
			datacenter = Datacenter{
				Name: *datacenterName,
			}
		}
		datacenter.AddNode(node)
		datacenters[*datacenterName] = datacenter
	}
	for name, datacenter := range datacenters {
		title := fmt.Sprintf("Datacenter: %s", name)
		ctx.String("%s\n", title)
		ctx.String("%s\n", strings.Repeat("=", len(title)))
		ctx.String("Status=Up/Down\n")
		ctx.String("|/ State=Normal/Leaving/Joining/Moving\n")
		ctx.String("%-2s  %-15s  %-10s  %-6s  %-16s  %-36s  %s\n", "--", "Address", "Load", "Tokens", "Owns (effective)", "Host ID", "Rack")
		for _, node := range datacenter.Nodes {
			ctx.String("%s%s  %-15s  %-10s  %-6s  %-16s  %-36s  %s\n", node.Status, node.State, node.Address, node.Load, node.Tokens, node.Owns, node.HostID, node.Rack)
		}
	}
	return nil
}
