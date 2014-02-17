#Raft Leader Election

This is a `go` implementation of Raft leader election process for distributed systems. With multiple instances initiated, the library takes care of leader election process and makes sure that there is at most one leader at any instance of time.

#Installation and Test

```
go get github.com/sk4x0r/raft-leader-election
go test github.com/sk4x0r/raft-leader-election
```

#Dependency
Library depends upon ZeroMQ 4, which can be installed from github.com/pebbe/zmq4


#Usage
New instance of server can be created by method `New()`.
```
s1:=raft.New(serverId, configFile)
```
`New()` method takes two parameter, and returns object of type `Server`.

| Parameter		| Type		| Description  
| -------------|:---------:| -----------
| serverId		| `int` 	| unuque id assigned to each server
| configFile	| `string`  | path of the file containing configuration of all the servers

For example of configuration file, see _config.json_ file in source.

Running instance of a server can be killed using `StopServer()` method.
```
s1.StopServer()
```

# License

The package is available under GNU General Public License. See the _LICENSE_ file.
