#pt-mongodb-summary
pt-mongodb-summary collects information about a MongoDB cluster.

##Usage
pt-mongodb-summary [options] [host:[port]]

Default host:port is `localhost:27017`. 
For better results, host must be a **mongos** server.

###Paramters

|Short|Long|default||
|---|---|---|---|
|u|user|empty|user name to use when connecting if DB auth is enabled|
|p|password|empty|password to use when connecting if DB auth is enabled|
|a|auth-db|admin|database used to establish credentials and privileges with a MongoDB server|


###Output example
```
# Instances ####################################################################################
ID    Host                         Type                                 ReplSet  Engine Status 
  0 localhost:17001                PRIMARY                                r1 
  1 localhost:17002                SECONDARY                              r1 
  2 localhost:17003                SECONDARY                              r1 
  0 localhost:18001                PRIMARY                                r2 
  1 localhost:18002                SECONDARY                              r2 
  2 localhost:18003                SECONDARY                              r2

# This host
# Mongo Executable #############################################################################
       Path to executable | /home/karl/tmp/MongoDB32Labs/3.0/bin/mongos
              Has symbols | No
# Report On 0 ########################################
                     User | karl
                PID Owner | mongos
                     Time | 2016-10-30 00:18:49 -0300 ART
                 Hostname | karl-HP-ENVY
                  Version | 3.0.11
                 Built On | Linux x86_64
                  Started | 2016-10-30 00:18:49 -0300 ART
                Databases | 0
              Collections | 0
                  Datadir | /data/db
                Processes | 0
             Process Type | mongos

# Running Ops ##################################################################################

Type         Min        Max        Avg
Insert           0          0          0/5s
Query            0          0          0/5s
Update           0          0          0/5s
Delete           0          0          0/5s
GetMore          0          0          0/5s
Command          0         22         16/5s

# Security #####################################################################################
Users 0
Roles 0
Auth  disabled
SSL   disabled


# Oplog ########################################################################################
Oplog Size     18660 Mb
Oplog Used     55 Mb
Oplog Length   0.91 hours
Last Election  2016-10-30 00:18:44 -0300 ART


# Cluster wide #################################################################################
            Databases: 3
          Collections: 17
  Sharded Collections: 1
Unsharded Collections: 16
    Sharded Data Size: 68 GB
  Unsharded Data Size: 0 KB
# Balancer (per day)
              Success: 6
               Failed: 0
               Splits: 0
                Drops: 0
```
