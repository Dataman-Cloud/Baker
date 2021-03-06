package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/Dataman-Cloud/baker/client"
	"github.com/Dataman-Cloud/baker/util"
)

var DisConfPushCmd = cli.Command{
	Name:  "push",
	Usage: "push config files to disconfig",
	Action: func(c *cli.Context) {
		if err := disConfPush(c); err != nil {
			logrus.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "application name",
		},
		cli.StringFlag{
			Name:  "label",
			Usage: "application label",
		},
		cli.StringFlag{
			Name:  "ymlfile",
			Usage: "ymlfile for uploading config files",
			Value: "push.yml",
		},
	},
}

type Upload struct {
	Files map[string]FileParam `yaml:"uploads"`
}

type FileParam struct {
	ContainerPath string `yaml:"container_path"` // container path.
}

func disConfPush(c *cli.Context) error {
	// validation.
	appName := c.String("name")
	if appName == "" {
		logrus.Fatal("no name in input.")
		return errors.New("no name in input.")
	}
	label := c.String("label")
	if label == "" {
		logrus.Fatal("no label in input.")
		return errors.New("no label input.")
	}
	ymlfile := c.String("ymlfile")
	buf, err := ioutil.ReadFile(ymlfile)
	if err != nil {
		logrus.Fatalf("error open ymlfile: %s", err)
		return err
	}
	upload := &Upload{}
	err = yaml.Unmarshal(buf, upload)
	if err != nil {
		logrus.Fatalf("error unmarshal ymlfile: %s", err)
		return err
	}

	// login baker server
	baseUri := c.GlobalString("server")
	client, err := client.NewHttpClient(baseUri, c.GlobalString("username"), c.GlobalString("password"))
	if err != nil {
		logrus.Fatalf("erro login baker server", err)
		return err
	}

	// disconf push
	logrus.Infof("push config files into disconf.")
	logrus.Infof("upload files size:%d", len(upload.Files))
	for filename, fileparam := range upload.Files {
		path, _ := os.Getwd()
		path += "/" + filename
		extraParams := map[string]string{
			"app-name":       appName,
			"label":          label,
			"timestamp":      strconv.FormatInt(time.Now().Unix(), 10),
			"container-path": fileparam.ContainerPath,
		}
		req, err := util.FileUploadRequest("http://"+baseUri+"/api/disconf/push", client.Token, "uploadfile", path, extraParams)
		if err != nil {
			logrus.Fatalf("error create disconf push request: %s ", err)
			return err
		}
		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			logrus.Fatalf("error disconf push request: %s", err)
			return err
		}
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatalf("error read response in disconf push: %s", err)
			return err
		}
		logrus.Info(resp.Status)
		logrus.Infof(string(respBody))
	}
	return nil
}
