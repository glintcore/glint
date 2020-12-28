
Glint Client
============

Copyright (C) 2017-2018 Index Data ApS.  This software is distributed under the
terms of the Apache License, Version 2.0.  See the file
[LICENSE](https://github.com/glintcore/glint-client/blob/master/LICENSE) for
more information.


##### Table of Contents  
Overview  
System requirements  
Installing the client  
Running the client  


Overview
--------

[Glint](https://glintcore.net) is open source software for sharing data
sets.  The Glint client is a command line client that connects to a [Glint
server](https://github.com/glintcore/glint-server) to post, describe, or
integrate data.


System requirements
-------------------

* Linux or macOS
* [Go](https://golang.org) 1.10 or later

Go is needed in order to compile Glint from source code.  On macOS running
[Homebrew](https://brew.sh/), Go can be installed with `brew install go`.


Installing the client
---------------------

First ensure that the `GOPATH` environment variable specifies a path
that can serve as your Go workspace directory, the place where Glint and
other Go packages will be installed.  For example, to set it to
`$HOME/go`:

```shell
$ export GOPATH=$HOME/go
```

Then to download and compile the Glint client:

```shell
$ go get -u -v github.com/glintcore/glint-client/...
```

The compiled executable file, `glint`, should appear in `$GOPATH/bin/`.  The
next section assumes that `glint` has been added to the path.


Running the client
------------------

### Configuring glint for the first time

Before the Glint client can connect to the server, it has to be
configured with details such as the server URL and your user name and
password for accessing the server.  This can be done using the
`config` command, which stores settings in the file `.glintconfig` in
your home directory:

```shell
$ glint config remote.url https://glintcore.net
$ glint config remote.user izzy
$ glint config remote.password
Enter new password:
```

The first time you connect to a server, change your password with the
`passwd` command, which will update the password both on the server
and in the local configuration file `.glintconfig`.


```shell
$ glint passwd
Enter new password:
(Confirming) Enter new password:
```

### Posting data on the server

A basic function of Glint is to share data by posting it on a server.
By default a data file is expected to be in CSV (comma-separated
values) format with a header line containing column names, although
other formats also can be used.  Posting data looks something like:

```shell
$ glint post ocean.csv
https://glintcore.net/izzy/ocean
```

Glint responds with a URL to the newly posted data set.  This URL can be
used to share the data set with others.  In a web browser the data appear as
a formatted table.

In other contexts the URL provides the data in a CSV or tab-delimited
form, which is easy for software or services to parse and is natively
supported by existing software such as R and Excel.

For example in R:

```r
> ocean <- read.csv("https://glintcore.net/izzy/ocean")
> ocean
  id                   t record site_id air_temp_avg baro_press_avg rel_hum_avg
1  1 2016-12-19 17:04:00   8109       1           NA          792.5       171.4
2  2 2016-12-19 17:34:00   8110       1           NA          789.0       163.7
3  3 2016-12-19 18:04:00   8111       1           NA          790.4       169.7
4  4 2016-12-19 18:34:00   8112       1        12.64         1012.0        92.7
5  5 2016-12-19 19:04:00   8113       1        13.26         1011.0        92.5
  dew_pt_avg vpr_press_avg wind_speed wind_dir stdev wind_gust wtr_lvl_avgreal
1         NA            NA      0.443    26.72 0.048     0.443        1.238093
2         NA            NA      0.443    26.72 0.048     0.443        1.237691
3         NA            NA      0.000     0.00 0.000     0.000        1.238556
4      11.50         1.355      0.000     0.00 0.000     0.000        1.237252
5      12.08         1.408      0.000     0.00 0.000     0.000        1.236872
```

Or with curl:

```shell
$ curl -o - https://glintcore.net/izzy/ocean
id,t,record,site_id,air_temp_avg,baro_press_avg,rel_hum_avg,dew_pt_avg,vpr_press_avg,wind_speed,wind_dir,stdev,wind_gust,wtr_lvl_avgreal
1,2016-12-19 17:04:00,8109,1,,792.5,171.399993896484,,,0.442999988794327,26.7199993133545,0.0480000004172325,0.442999988794327,1.2380930185318
2,2016-12-19 17:34:00,8110,1,,789,163.699996948242,,,0.442999988794327,26.7199993133545,0.0480000004172325,0.442999988794327,1.23769104480743
3,2016-12-19 18:04:00,8111,1,,790.400024414062,169.699996948242,,,0,0,0,0,1.23855602741241
4,2016-12-19 18:34:00,8112,1,12.6400003433228,1012,92.6999969482422,11.5,1.35500001907349,0,0,0,0,1.23725199699402
5,2016-12-19 19:04:00,8113,1,13.2600002288818,1011,92.5,12.0799999237061,1.40799999237061,0,0,0,0,1.23687195777893
```

### Changing how data are retrieved

Glint interprets commands added to the end of data set URLs as changing how the
data should be retrieved.  (The syntax is roughly based on
[THUMP](https://tools.ietf.org/html/draft-kunze-thump-03).) For example:

```shell
$ curl -o - https://glintcore.net/izzy/ocean?show(t,wind_dir)as(tsv)
```

The commands, `show()` and `as()`, have been added to the end of the
URL created in the previous example: `show()` asks Glint to select a
subset of the columns to retrieve, and `as()` sets the format of the
retrieved data, in this case TSV (tab-separated values).


### Adding metadata

Another feature of Glint is the ability to add metadata to data sets
easily, especially standard metadata on columns.  These examples
demonstrate associating columns with [Dublin
Core](http://dublincore.org) and
[YAMZ](https://github.com/nassar/yamz) metadata:

```shell
$ glint md ocean.t dc:date

$ glint md ocean.wind_speed yamz:h3846
```

The first example tags the column `t` with the metadata element
`dc:date` which identifies it as a date/time field.  The second
example tags the column `wind_speed` with `yamz:h3846`, a YAMZ
identifier for wind speed data.


### Integrating data with services

Glint provides a basic plotting service for time series data at
`/plot-time-series`, as a demonstration of integrating data with
services.  The service accepts any data set having a column that has
been tagged with the metadata elements, `dc:date` or `yamz:h1317`,
both representing a date/time.  It accepts a data set as input in the
form of a Glint URL that refers to the data.  Glint can include
metadata tags in the header line of a data set, in a format that is
easy for the service to parse, e.g.:

```
t{dc:date},air_temp_avg,wind_speed{yamz:h3846},wind_dir
```

The time series plotting service asks the Glint server to include these
metadata tags by adding `md()` to the URL that was provided as input.
For example, suppose that this URL is given as input to the service:

```http
https://glintcore.net/izzy/ocean?show(t,air_temp_avg,wind_speed)
```

The service adds `md()` to the end of the URL before using it to
retrieve the data:

```http
https://glintcore.net/izzy/ocean?show(t,air_temp_avg,wind_speed)md()
```

This allows the service to identify the time column to use for the
x-axis and plot other columns on the y-axis.


