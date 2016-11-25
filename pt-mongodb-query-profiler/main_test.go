package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/percona/toolkit-go/mongolib/proto"

	"gopkg.in/mgo.v2/dbtest"
)

var Server dbtest.DBServer

func TestMain(m *testing.M) {
	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	os.Setenv("CHECK_SESSIONS", "0")
	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	dat, err := ioutil.ReadFile("test/sample/system.profile.json")
	if err != nil {
		fmt.Printf("cannot load fixtures: %s", err)
		os.Exit(1)
	}

	var docs []proto.SystemProfile
	err = json.Unmarshal(dat, &docs)
	c := Server.Session().DB("samples").C("system_profile")
	for _, doc := range docs {
		c.Insert(doc)
	}

	retCode := m.Run()

	Server.Session().Close()
	Server.Session().DB("samples").DropDatabase()

	// Stop shuts down the temporary server and removes data on disk.
	Server.Stop()

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestGetData(t *testing.T) {
	var docs []interface{}
	it := Server.Session().DB("samples").C("system_profile").Find(nil).Iter()
	err := Server.Session().DB("samples").C("system_profile").Find(nil).All(&docs)
	if err != nil {
		t.Errorf("cannot read docs: %s", err)
	}
	tests := []struct {
		name string
		i    iter
		want []stat
	}{
		{
			name: "test 1",
			i:    it,
			want: []stat{
				stat{
					ID:          "ea170e2cafb1337755c8b3d5ae4437f4",
					Fingerprint: "find",
					Query:       map[string]interface{}{"find": "col1"},
					Count:       10,
					TableScan:   false,
					NScanned:    []float64{71, 72, 73, 74, 75, 76, 77, 78, 79, 80},
					NReturned:   []float64{71, 72, 73, 74, 75, 76, 77, 78, 79, 80},
					QueryTime:   []float64{19, 20, 21, 22, 23, 24, 25, 26, 27, 28},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getData(tt.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getData() = %#v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFingerprint(t *testing.T) {
	tests := []struct {
		name  string
		query map[string]interface{}
		want  string
	}{
		{
			query: map[string]interface{}{"query": map[string]interface{}{}, "orderby": map[string]interface{}{"ts": -1}},
			want:  "orderby,query,ts",
		},
		{
			query: map[string]interface{}{"find": "system.profile", "filter": map[string]interface{}{}, "sort": map[string]interface{}{"$natural": 1}},
			want:  "$natural,filter,find,sort",
		},
		{

			query: map[string]interface{}{"collection": "system.profile", "batchSize": 0, "getMore": 18531768265},
			want:  "batchSize,collection,getMore",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fingerprint(tt.query); got != tt.want {
				t.Errorf("fingerprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimesLen(t *testing.T) {
	tests := []struct {
		name string
		a    times
		want int
	}{
		{
			name: "Times.Len",
			a:    []time.Time{time.Now()},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Len(); got != tt.want {
				t.Errorf("times.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimesSwap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	t1 := time.Now()
	t2 := t1.Add(1 * time.Minute)
	tests := []struct {
		name string
		a    times
		args args
	}{
		{
			name: "Times.Swap",
			a:    times{t1, t2},
			args: args{i: 0, j: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.Swap(tt.args.i, tt.args.j)
			if tt.a[0] != t2 || tt.a[1] != t1 {
				t.Errorf("%s has (%v, %v) want (%v, %v)", tt.name, tt.a[0], tt.a[1], t2, t1)
			}
		})
	}
}

func TestTimesLess(t *testing.T) {
	type args struct {
		i int
		j int
	}
	t1 := time.Now()
	t2 := t1.Add(1 * time.Minute)
	tests := []struct {
		name string
		a    times
		args args
		want bool
	}{
		{
			name: "Times.Swap",
			a:    times{t1, t2},
			args: args{i: 0, j: 1},
			want: true,
		},
		{
			name: "Times.Swap",
			a:    times{t2, t1},
			args: args{i: 0, j: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("times.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}
