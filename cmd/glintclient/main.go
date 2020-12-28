package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/glintdb/glintweb/api"
	"github.com/nassibnassar/goconfig/ini"
	"github.com/urfave/cli"
)

var glintconfig *ini.Config
var glintconfigfilename string

// fileModeRW is the umask "-rw-------".
const fileModeRW = 0600

func trimSlash(s string) string {
	return strings.TrimRight(s, "/")
}

func responseBody(resp *http.Response) (string, error) {
	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}

func responseBodyError(resp *http.Response) error {
	body, err := responseBody(resp)
	if err != nil {
		return errors.New(body + " [" + err.Error() + "]")
	}
	return errors.New(body)
}

func getUserPassword() (string, string, error) {
	user := glintconfig.Get("remote", "user")
	if user == "" {
		return "", "", errors.New("User not specified")
	}
	password := glintconfig.Get("remote", "password")
	if password == "" {
		var err error
		password, err = inputPassword(
			"Enter current password: ", false)
		if err != nil {
			return "", "", err
		}
	}
	return user, password, nil
}

func removeExtension(s string) string {
	i := strings.LastIndexByte(s, '.')
	if i == -1 {
		return s
	}
	return s[0:i]
}

func cliMd(c *cli.Context) error {
	user, password, err := getUserPassword()
	if err != nil {
		return err
	}

	fileAttr := c.Args().Get(0)
	if fileAttr == "" {
		return errors.New("Attribute not specified")
	}

	metadata := c.Args().Get(1)
	if metadata == "" {
		return errors.New("Metadata not specified")
	}

	var req api.MetadataRequest
	req.Metadata = metadata

	reqbody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}
	remote := trimSlash(glintconfig.Get("remote", "url"))
	url := remote + "/" + user + "/" + fileAttr
	httpreq, err := http.NewRequest(http.MethodPut, url,
		bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}
	httpreq.SetBasicAuth(user, password)
	httpreq.Header.Set("Content-Type", "application/json")

	httpresp, err := client.Do(httpreq)
	if err != nil {
		return err
	}

	if httpresp.StatusCode != http.StatusOK {
		if httpresp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Server at '" + remote +
				"' did not accept the username/password")
		}
		fmt.Println(httpresp.StatusCode)
		fmt.Println("(1)")
		return responseBodyError(httpresp)
	}

	// respbody, err := ioutil.ReadAll(httpresp.Body)
	// if err != nil {
	// 	fmt.Println("(2)")
	// 	return err
	// }

	// var resp api.PostResponse
	// err = json.Unmarshal(respbody, &resp)
	// if err != nil {
	// 	fmt.Println("(3)")
	// 	return err
	// }

	// fmt.Printf("%s\n", resp.Url)

	return nil
}

func cliPost(c *cli.Context) error {
	user, password, err := getUserPassword()
	if err != nil {
		return err
	}

	var req api.PostRequest
	dataFile := c.Args().Get(0)
	if dataFile == "" {
		return errors.New("Data file not specified")
	}
	fileinfo, err := os.Stat(dataFile)
	if err != nil {
		return err
	}

	// TODO Read line by line for conversion.
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}
	req.Data = strings.Replace(string(data), "\n", "\\n", -1)

	reqbody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}
	remote := trimSlash(glintconfig.Get("remote", "url"))
	fileName := removeExtension(fileinfo.Name())
	url := remote + "/" + user + "/" + fileName
	//fmt.Printf("url: [%s]\n", url)
	//fmt.Printf("req: [%v]\n", string(reqbody))
	httpreq, err := http.NewRequest(http.MethodPut, url,
		bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}
	httpreq.SetBasicAuth(user, password)
	httpreq.Header.Set("Content-Type", "application/json")

	httpresp, err := client.Do(httpreq)
	if err != nil {
		return err
	}

	if httpresp.StatusCode != http.StatusCreated {
		if httpresp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Server at '" + remote +
				"' did not accept the username/password")
		}
		fmt.Println(httpresp.StatusCode)
		return responseBodyError(httpresp)
	}

	respbody, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		return err
	}

	var resp api.PostResponse
	err = json.Unmarshal(respbody, &resp)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", resp.Url)

	return nil
}

func cliDelete(c *cli.Context) error {
	user, password, err := getUserPassword()
	if err != nil {
		return err
	}

	fileName := c.Args().Get(0)
	if fileName == "" {
		return errors.New("Data file not specified")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}
	remote := trimSlash(glintconfig.Get("remote", "url"))
	url := remote + "/" + user + "/" + fileName
	//fmt.Printf("url: [%s]\n", url)
	//fmt.Printf("req: [%v]\n", string(reqbody))
	httpreq, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	httpreq.SetBasicAuth(user, password)

	httpresp, err := client.Do(httpreq)
	if err != nil {
		return err
	}

	if httpresp.StatusCode != http.StatusNoContent {
		if httpresp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Server at '" + remote +
				"' did not accept the username/password")
		}
		fmt.Println(httpresp.StatusCode)
		return responseBodyError(httpresp)
	}

	fmt.Printf("OK\n")

	return nil
}

func cliLogin(c *cli.Context) error {
	user, password, err := getUserPassword()
	if err != nil {
		return err
	}

	var req api.LoginRequest

	//req.Data = strings.Replace(string(data), "\n", "\\n", -1)

	reqbody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}
	remote := trimSlash(glintconfig.Get("remote", "url"))
	url := remote + "/login"
	//fmt.Printf("url: [%s]\n", url)
	//fmt.Printf("req: [%v]\n", string(reqbody))
	httpreq, err := http.NewRequest(http.MethodPost, url,
		bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}
	httpreq.SetBasicAuth(user, password)
	httpreq.Header.Set("Content-Type", "application/json")

	httpresp, err := client.Do(httpreq)
	if err != nil {
		return err
	}

	if httpresp.StatusCode != http.StatusCreated {
		if httpresp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Server at '" + remote +
				"' did not accept the username/password")
		}
		fmt.Println(httpresp.StatusCode)
		return responseBodyError(httpresp)
	}

	respbody, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		return err
	}

	var resp api.LoginResponse
	err = json.Unmarshal(respbody, &resp)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", resp.SessionId)

	return nil
}

func cliPasswd(c *cli.Context) error {
	user, password, err := getUserPassword()
	if err != nil {
		return err
	}

	var req api.AccountPasswordRequest
	req.Password, err = inputPassword("Enter new password: ", true)
	if err != nil {
		return err
	}

	reqbody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}
	remote := trimSlash(glintconfig.Get("remote", "url"))
	httpreq, err := http.NewRequest("POST", remote+"/account/password",
		bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}
	httpreq.SetBasicAuth(user, password)
	httpreq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpreq)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Server at '" + remote +
				"' did not accept the username/password")
		}
		return responseBodyError(resp)
	}

	fmt.Printf("Updated password on server\n")

	glintconfig.Set("remote", "password", req.Password)
	fmt.Printf("Updated password in configuration file %s\n", glintconfigfilename)

	return nil
}

func writeConfigFile(filename string) error {

	// Delete file so that we can set the proper umask.
	_ = os.Remove(filename)

	var f *os.File
	var err error
	f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		fileModeRW)
	if err != nil {
		return err
	}
	defer f.Close()

	err = glintconfig.WriteTo(f)
	return err
}

func cliConfig(c *cli.Context) error {

	key := c.Args().Get(0)
	if key == "" {
		return errors.New("Missing key and value")
	}

	var err error
	value := c.Args().Get(1)
	if value == "" {
		if key != "remote.password" {
			return errors.New("Missing value")
		}
		// Passwords are a special case; we can use the
		// terminal for input.
		value, err = inputPassword(
			"Enter new password: ", false)
		if err != nil {
			return err
		}
		if value == "" {
			return errors.New("No password specified; " +
				"password unchanged")
		}
	}

	var sk []string = strings.Split(key, ".")
	glintconfig.Set(sk[0], sk[1], value)
	err = writeConfigFile(glintconfigfilename)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Read configuration file.
	glintconfigfilename = os.Getenv("HOME") + "/" + ".glintconfig"
	glintconfig = readConfig(glintconfigfilename)
	// Run commands specified on the command line.
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "print the Glint version",
	}
	app := cli.NewApp()
	app.Name = "glint"
	app.Version = "0"
	app.HideVersion = true
	app.HelpName = "glint"
	app.Usage = "Glint client for sharing and integrating data sets"
	app.UsageText = "glint [command] [arguments]"
	app.EnableBashCompletion = true
	/*
		app.Flags = []cli.Flag{
			// Verbose flag not currently implemented.
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "enable verbose output",
			},
		}
	*/
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "config",
			Usage:     "Sets configuration options",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				err := cliConfig(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "passwd",
			Usage:     "Changes the user's password",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				err := cliPasswd(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "login",
			Hidden:    true,
			Usage:     "Authenticates with server and retrieves a session id",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				err := cliLogin(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "post",
			Usage:     "Publishes data on the server",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				// TODO Implement --type and --no-header flags.
				cli.StringFlag{
					Name:  "type",
					Usage: "file format",
				},
				cli.StringFlag{
					Name: "no-header",
					Usage: "data file does not include " +
						"column names",
				},
			},
			Action: func(c *cli.Context) error {
				err := cliPost(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "delete",
			Usage:     "Deletes data from the server",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				err := cliDelete(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		cli.Command{
			Name:      "md",
			Usage:     "Adds metadata to an attribute",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				err := cliMd(c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
