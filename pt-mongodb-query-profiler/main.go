package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/howeyc/gopass"
	"github.com/montanaflynn/stats"
	"github.com/pborman/getopt"
	"github.com/percona/toolkit-go/mongolib/proto"
	"github.com/y0ssar1an/q"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type iter interface {
	All(result interface{}) error
	Close() error
	Err() error
	For(result interface{}, f func() error) (err error)
	Next(result interface{}) bool
	Timeout() bool
}

type options struct {
	Help     bool
	Host     string
	User     string
	Password string
	Database string
	AuthDB   string
	Debug    bool
}

const (
	MAX_DEPTH_LEVEL = 10
)

type statsArray []stat

func (a statsArray) Len() int           { return len(a) }
func (a statsArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a statsArray) Less(i, j int) bool { return a[i].Count < a[j].Count }

type times []time.Time

func (a times) Len() int           { return len(a) }
func (a times) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a times) Less(i, j int) bool { return a[i].Before(a[j]) }

type stat struct {
	ID          string
	Fingerprint string
	Namespace   string
	Query       map[string]interface{}
	Count       int
	TableScan   bool
	NScanned    []float64
	NReturned   []float64
	QueryTime   []float64 // in milliseconds
	LockTime    times
	BlockedTime times
	FirstSeen   time.Time
	LastSeen    time.Time
}

type groupKey struct {
	Fingerprint string
	Namespace   string
}

type statistics struct {
	Pct    float64
	Total  float64
	Min    float64
	Max    float64
	Avg    float64
	Pct95  float64
	StdDev float64
	Median float64
}

type queryInfo struct {
	Rank        int
	ID          string
	Count       int
	Ratio       float64
	Fingerprint string
	Namespace   string
	Scanned     statistics
	Returned    statistics
	QueryTime   statistics
	FirstSeen   time.Time
	LastSeen    time.Time
}

func main() {

	opts, err := getOptions()
	if err != nil {
		log.Printf("error processing commad line arguments: %s", err)
		os.Exit(1)
	}
	if opts.Help {
		getopt.Usage()
		return
	}

	di := getDialInfo(opts)
	if di.Database == "" {
		log.Printf("must indicate a database")
		getopt.PrintUsage(os.Stderr)
		os.Exit(2)
	}

	session, err := mgo.DialWithInfo(di)
	if err != nil {
		log.Printf("error connecting to the db %s", err)
		os.Exit(3)
	}

	i := session.DB(di.Database).C("system.profile").Find(bson.M{"op": bson.M{"$nin": []string{"getmore", "delete"}}}).Sort("-$natural").Iter()
	queries := getData(i)

	queryStats := calcQueryStats(queries)

	t := template.Must(template.New("oplogInfo").Parse(getTemplate()))
	for _, qs := range queryStats {
		q.Q(qs)
		t.Execute(os.Stdout, qs)
	}

}

func calcQueryStats(queries []stat) []queryInfo {
	queryStats := []queryInfo{}
	_, totalScanned, totalReturned, totalQueryTime := calcTotals(queries)
	for rank, query := range queries {
		qi := queryInfo{
			Rank:        rank,
			Count:       query.Count,
			ID:          query.ID,
			Fingerprint: query.Fingerprint,
			Scanned:     calcStats(query.NScanned),
			Returned:    calcStats(query.NReturned),
			QueryTime:   calcStats(query.QueryTime),
			FirstSeen:   query.FirstSeen,
			LastSeen:    query.LastSeen,
			Namespace:   query.Namespace,
		}
		if totalScanned > 0 {
			qi.Scanned.Pct = qi.Scanned.Total * 100 / totalScanned
		}
		if totalReturned > 0 {
			qi.Returned.Pct = qi.Returned.Total * 100 / totalReturned
		}
		if totalQueryTime > 0 {
			qi.QueryTime.Pct = qi.QueryTime.Total * 100 / totalQueryTime
		}
		if qi.Returned.Total > 0 {
			qi.Ratio = qi.Scanned.Total / qi.Returned.Total
		}
		queryStats = append(queryStats, qi)
	}
	return queryStats
}

func getTotals(queries []stat) stat {

	qt := stat{}
	for _, query := range queries {
		qt.NScanned = append(qt.NScanned, query.NScanned...)
		qt.NReturned = append(qt.NReturned, query.NReturned...)
		qt.QueryTime = append(qt.QueryTime, query.QueryTime...)
	}
	return qt

}

func calcTotals(queries []stat) (totalCount int, totalScanned, totalReturned, totalQueryTime float64) {

	for _, query := range queries {
		totalCount += query.Count

		scanned, _ := stats.Sum(query.NScanned)
		totalScanned += scanned

		returned, _ := stats.Sum(query.NReturned)
		totalReturned += returned

		queryTime, _ := stats.Sum(query.QueryTime)
		totalQueryTime += queryTime
	}
	return
}

func calcStats(samples []float64) statistics {
	var s statistics
	s.Total, _ = stats.Sum(samples)
	s.Min, _ = stats.Min(samples)
	s.Max, _ = stats.Max(samples)
	s.Avg, _ = stats.Mean(samples)
	s.Pct95, _ = stats.PercentileNearestRank(samples, 95)
	s.StdDev, _ = stats.StandardDeviation(samples)
	s.Median, _ = stats.Median(samples)
	return s
}

func getData(i iter) []stat {
	var doc proto.SystemProfile
	stats := make(map[groupKey]*stat)

	for i.Next(&doc) && i.Err() == nil {

		if len(doc.Query) > 0 {
			query := doc.Query
			if squery, ok := doc.Query["$query"]; ok {
				if ssquery, ok := squery.(map[string]interface{}); ok {
					query = ssquery
				}
			}
			fp := fingerprint(query)
			var s *stat
			var ok bool
			key := groupKey{
				Fingerprint: fp,
				Namespace:   doc.Ns,
			}
			if s, ok = stats[key]; !ok {
				s = &stat{
					ID:          fmt.Sprintf("%x", md5.Sum([]byte(fp+doc.Ns))),
					Fingerprint: fp,
					Namespace:   doc.Ns,
					TableScan:   false,
					Query:       query,
				}
				stats[key] = s
			}
			s.Count++
			s.NScanned = append(s.NScanned, float64(doc.DocsExamined))
			s.NReturned = append(s.NReturned, float64(doc.Nreturned))
			s.QueryTime = append(s.QueryTime, float64(doc.Millis))
			var zeroTime time.Time
			if s.FirstSeen == zeroTime || s.FirstSeen.After(doc.Ts) {
				s.FirstSeen = doc.Ts
			}
			if s.LastSeen == zeroTime || s.LastSeen.Before(doc.Ts) {
				s.LastSeen = doc.Ts
			}
		}
	}

	// We need to sort the data but a hash cannot be sorted so, convert the hash having
	// the results to a slice
	sa := statsArray{}
	for _, s := range stats {
		sa = append(sa, *s)
	}

	// Sort by count, descending order
	sort.Sort(sort.Reverse(sa))
	return sa
}

// TODO REMOVE. Used for debug.
func format(title string, templateData interface{}) string {
	txt, _ := json.MarshalIndent(templateData, "", "    ")
	return title + "\n" + string(txt)
}

// TODO REMOVE. Used for debug.
func write(title string, templateData interface{}) {
	txt, _ := json.MarshalIndent(templateData, "", "    ")
	f, _ := os.Create("test/sample/" + title + ".json")
	f.Write(txt)
	f.Close()
}

func getOptions() (*options, error) {
	opts := &options{Host: "localhost:27017"}
	getopt.BoolVarLong(&opts.Help, "help", '?', "Show help")
	getopt.StringVarLong(&opts.User, "user", 'u', "", "username")
	getopt.StringVarLong(&opts.Password, "password", 'p', "", "password").SetOptional()
	getopt.StringVarLong(&opts.AuthDB, "auth-db", 'a', "admin", "database used to establish credentials and privileges with a MongoDB server")
	getopt.StringVarLong(&opts.Database, "database", 'd', "", "database to profile")
	getopt.SetParameters("host[:port][/database]")

	getopt.Parse()
	if opts.Help {
		return opts, nil
	}

	args := getopt.Args() // host is a positional arg
	if len(args) > 0 {
		opts.Host = args[0]
	}

	if getopt.IsSet("password") && opts.Password == "" {
		print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return nil, err
		}
		opts.Password = string(pass)
	}

	return opts, nil
}

func getDialInfo(opts *options) *mgo.DialInfo {
	di, _ := mgo.ParseURL(opts.Host)
	di.FailFast = true

	if getopt.IsSet("user") {
		di.Username = opts.User
	}
	if getopt.IsSet("password") {
		di.Password = opts.Password
	}
	if getopt.IsSet("auth-db") {
		di.Source = opts.AuthDB
	}

	if getopt.IsSet("database") {
		di.Database = opts.Database
	}

	return di
}

func fingerprint(query map[string]interface{}) string {
	return strings.Join(keys(query, 0), ",")
}

func keys(query map[string]interface{}, level int) []string {
	ks := []string{}
	for key, value := range query {
		ks = append(ks, key)
		if m, ok := value.(map[string]interface{}); ok {
			level++
			if level <= MAX_DEPTH_LEVEL {
				ks = append(ks, keys(m, level)...)
			}
		}
	}
	sort.Strings(ks)
	return ks
}

func getTemplate() string {

	t := `
# Query {{.Rank}}: x.xx QPS, ID {{.ID}}
# Ratio {{.Ratio}}
# Time range: {{.FirstSeen}} to {{.LastSeen}}
# Attribute       pct   total     min     max     avg     95%  stddev  median
# =============== === ======= ======= ======= ======= ======= ======= =======
# Count               {{printf "% 7d" .Count}}
# Exec Time ms    {{printf "% 3.0f" .QueryTime.Pct}} {{printf "% 7.0f" .QueryTime.Total}} {{printf "% 7.0f" .QueryTime.Min}} {{printf "% 7.0f" .QueryTime.Max}} {{printf "% 7.0f" .QueryTime.Avg}} {{printf "% 7.0f" .QueryTime.Pct95}} {{printf "% 7.0f" .QueryTime.StdDev}} {{printf "% 7.0f" .QueryTime.Median}}
# Lock time         
# Docs Scanned    {{printf "% 3.0f" .Scanned.Pct}} {{printf "% 7.0f" .Scanned.Total}} {{printf "% 7.0f" .Scanned.Min}} {{printf "% 7.0f" .Scanned.Max}} {{printf "% 7.2f" .Scanned.Avg}} {{printf "% 7.2f" .Scanned.Pct95}} {{printf "% 7.2f" .Scanned.StdDev}} {{printf "% 7.2f" .Scanned.Median}}
# Docs Returned   {{printf "% 3.0f" .Returned.Pct}} {{printf "% 7.0f" .Returned.Total}} {{printf "% 7.0f" .Returned.Min}} {{printf "% 7.0f" .Returned.Max}} {{printf "% 7.2f" .Returned.Avg}} {{printf "% 7.2f" .Returned.Pct95}} {{printf "% 7.2f" .Returned.StdDev}} {{printf "% 7.2f" .Returned.Median}}
# Bytes sent        
# Query size   
# Boolean:
# Full scan    
# String:
# Namespaces      {{.Namespace}}
# Fingerprint     {{.Fingerprint}}
`
	return t
}
