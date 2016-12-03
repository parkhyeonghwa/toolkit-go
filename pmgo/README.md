#Percona MongoDB Go driver

This is just a collection of interfaces around the structures in mgo, ([Rich MongoDB driver for Go](https://labix.org/mgo)) to be able to mock methods in the driver.  

The motivation for this package is that there are certain things, like errors, that cannot be tested/reproduced using a real db connection.  
Also, for some of our tests we need very specific MongoDB configuration. Tests for some parts of our code need 2 replicas, config and mongo**s** servers and that's not easily reproducible in all CI environments.

##How to use it

This package is almost a drop-in replacement with the exception that you need to use the `Dialer` interface.

```
package main

import (
    "github.com/percona/toolkit-go/pmgo"
    "gopkg.in/mgo.v2/bson"
)

type testT struct {
    ID   int    `bson:"id"`
    Name string `bson:"name"`
}

func main() {
    dialer := pmgo.NewDialer()
    session, err := dialer.Dial("localhost")
    if err != nil {
        print(err)
        return
    }

    var test testT
    _ = session.DB("test").C("testc").Find(bson.M{"id": 1}).One(&test)

}
```

