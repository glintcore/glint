package server

import (
	"bytes"
	"net/http"
	"strings"
)

// ParseBasic makes a very rudimentary parse of a thump command sequence.
func thumpParseBasic(r *http.Request) map[string][]string {
	var cmdlist []string = strings.Split(r.URL.RawQuery, ")")
	var m = make(map[string][]string)
	var x int
	for x = range cmdlist {
		if cmdlist[x] != "" {
			var cmd_args []string = strings.Split(cmdlist[x], "(")
			m[cmd_args[0]] = strings.Split(cmd_args[1], ",")
		}
	}
	return m
}

// ShowBasic writes columnar data, selecting only columns specified by show.
func thumpShowBasic(data string, show []string) string {
	var out bytes.Buffer
	var show_n []int
	var rows []string = strings.Split(data, "\n")
	var x int
	for x = range rows {
		if rows[x] == "" {
			continue
		}
		var d []string = strings.Split(rows[x], ",")
		if x == 0 {
			var first bool = true
			var y int
			for y = range d {
				var z int
				for z = range show {
					if d[y] == show[z] {
						if first {
							first = false
						} else {
							out.WriteString(",")
						}
						out.WriteString(show[z])
						show_n = append(show_n, y)
					}
				}
			}
		} else {
			var first bool = true
			var y int
			for y = range d {
				var z int
				for z = range show_n {
					if y == show_n[z] {
						if first {
							first = false
						} else {
							out.WriteString(",")
						}
						out.WriteString(d[y])
					}
				}
			}
		}
		out.WriteString("\n")
	}
	return out.String()
}
