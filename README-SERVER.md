
Glint Server
============

Copyright (C) 2017-2018 Index Data ApS.  This software is distributed under the
terms of the Apache License, Version 2.0.  See the file
[LICENSE](https://github.com/glintcore/glint-server/blob/master/LICENSE) for
more information.


##### Table of Contents  
Overview  
System requirements  
Installing the server  
Running the server


Overview
--------

[Glint](https://glintcore.net) is open source software for sharing data
sets.  The Glint server offers a service for sharing, describing, and
integrating data.  The [Glint
client](https://github.com/glintcore/glint-client) can be used to connect to
the server.


System requirements
-------------------

* Linux 2.6.24 or later
* PostgreSQL 9.2.22 or later
* [Go](https://golang.org) 1.10 or later

PostgreSQL has been used to store data in the server prototype.  It will be
removed as a dependency in future versions of Glint.

Go is needed in order to compile the server from source code.


Installing the server
---------------------

First ensure that the `GOPATH` environment variable specifies a path that
can serve as your Go workspace directory, the place where Glint and other Go
packages will be installed.  For example, to set it to `$HOME/go`:

```shell
$ export GOPATH=$HOME/go
```

Then to download and compile the Glint server:

```shell
$ go get -u -v github.com/glintcore/glint-server/...
```

The compiled executable file, `glintserver`, should appear in
`$GOPATH/bin/`.


Running the server
------------------

### Configuration file

The included file `glintserver.conf` can be used as a template for
configuring the server:

```ini
# Sample glintserver configuration file

[core]
# datadir is a directory where the server will store and manage data:
# (This is not yet fully implemented.)
datadir = /var/lib/glint/

[log]
# log is a file that the server will write access and error logs to:
file = /var/log/glint/glintserver.log

[http]
# host is the TCP host address to listen on:
host = glintcore.net
# port is the TCP port number to listen on:
port = 443
# tlscert is a file containing a certificate for the server:
tlscert = /etc/letsencrypt/live/glintcore.net/fullchain.pem
# tlskey is a file containing the matching private key for the server:
tlskey = /etc/letsencrypt/live/glintcore.net/privkey.pem

# The database section specifies connection parameters for PostgreSQL:
# (This is where data are currently stored.)
[database]
host = localhost
port = 5432
user = glint
password = password_goes_here
dbname = glint
```

The server looks for a configuration file like this one in a location
specified by the `GLINTSERVER_CONFIG_FILE` environment variable, which in
bash can be set with, for example:

```shell
$ export GLINTSERVER_CONFIG_FILE=/etc/glint/glintserver.conf
```

### Running the server on a privileged port

The default port for HTTPS URLs is 443, which makes this a good port to have
the server listen on, as in the sample configuration above.  However, since
443 is a privileged port, we need `setcap` to enable the server to use it.
In Debian Linux, `setcap` is installed with:

```shell
$ sudo aptitude install libcap2-bin
```

Then run `setcap` on the `glintserver` executable:

```shell
$ sudo setcap 'cap_net_bind_service=+ep' $GOPATH/bin/glintserver
```

Alternatively, glintserver can listen on a non-privileged port like 8443.
This only requires changing the HTTP port setting in the configuration file:

```ini
[http]
port = 8443
```

**Note: Do not run the server as root.**

### Starting the server

To start the server:

```shell
$ nohup glintserver run &
```

### Adding a user

```shell
$ glintserver adduser --user izzy --fullname 'Isaac Newton' --email 'izzy@indexdata.com'
Enter new password:
```

### Changing a user's password

```shell
$ glintserver passwd --user izzy
Enter new password:
```


