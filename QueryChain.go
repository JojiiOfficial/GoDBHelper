package godbhelper

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

//QueryChain a list of SQL queries over time
type QueryChain struct {
	Name    string     `json:"name"`
	Order   int        `json:"order"`
	Queries []SQLQuery `json:"queries"`
}

//SQLQuery a query
type SQLQuery struct {
	VersionAdded float32  `json:"vs"`
	QueryString  string   `json:"query"`
	Params       []string `json:"params"`
	FqueryString string   `json:"queryf"`
	Fparams      []string `json:"fparams"`
}

//NewQueryChain QueryChain constructor
func NewQueryChain(name string, order int) *QueryChain {
	return &QueryChain{
		Name:  name,
		Order: order,
	}
}

//RestoreQueryChain loads an exported queryChain from file
func RestoreQueryChain(file string) (*QueryChain, error) {
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

//ExportQueryChain saves/exports a queryChain to a file
func (queryChain *QueryChain) ExportQueryChain(file string, perm os.FileMode) error {
	data, err := json.Marshal(queryChain)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, perm)
}

//LoadQueries loads queries from a .sql file and executes the statements (row for row).
//The SQLQuery Version of the statements are 0.
//This is intended to initialize the database-schema
func LoadQueries(name, file string, chainOrder int) (*QueryChain, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(f)
	queryChain := QueryChain{}

	for fileScanner.Scan() {
		sql := fileScanner.Text()
		queryChain.Queries = append(queryChain.Queries, SQLQuery{
			VersionAdded: 0,
			QueryString:  sql,
		})
	}
	queryChain.Name = name
	queryChain.Order = chainOrder
	return &queryChain, nil
}

//LoadQueries loads queries from a .sql file and executes the statements (row for row).
//The SQLQuery Version of the statements are 0.
//This is intended to initialize the database-schema
func (dbhelper *DBhelper) LoadQueries(name, file string, chainOrder int) error {
	queries, err := LoadQueries(name, file, chainOrder)
	if err != nil {
		return dbhelper.handleErrHook(err)
	}
	dbhelper.AddQueryChain(*queries)
	return nil
}

//InitSQL init sql obj
type InitSQL struct {
	Query   string
	Params  string
	FParams []string
}

//CreateInitVersionSQL creates SQLQuery[] for init version
func CreateInitVersionSQL(arg ...InitSQL) []SQLQuery {
	var queries []SQLQuery

	for _, query := range arg {
		queries = append(queries, SQLQuery{
			VersionAdded: 0,
			Fparams:      query.FParams,
			FqueryString: query.Query,
			QueryString:  query.Query,
		})
	}

	return queries
}
