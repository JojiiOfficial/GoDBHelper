package godbhelper

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//QueryChain a list of SQL queries over time
type QueryChain struct {
	Name    string     `json:"name"`
	Queries []SQLQuery `json:"queries"`
}

//SQLQuery a query
type SQLQuery struct {
	VersionAdded float32  `json:"vs"`
	QueryString  string   `json:"query"`
	Params       []string `json:"params"`
}

//NewQueryChain QueryChain constructor
func NewQueryChain(name string) *QueryChain {
	return &QueryChain{
		Name: name,
	}
}

//NewQueryChainFromfile loads querychain from file
func NewQueryChainFromfile(file string) (*QueryChain, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var chain QueryChain
	err = json.Unmarshal(b, &chain)
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

//SaveQueryChain saves querychain to file
func (queryChain *QueryChain) SaveQueryChain(file string, perm os.FileMode) error {
	data, err := json.Marshal(queryChain)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, perm)
}
