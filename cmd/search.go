package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/karust/openserp/baidu"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/google"
	"github.com/karust/openserp/yandex"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var searchCMD = &cobra.Command{
	Use:     "search",
	Aliases: []string{"find"},
	Short:   "Search results using chosen web search engine (google, yandex, baidu)",
	Args:    cobra.MatchAll(cobra.OnlyValidArgs, cobra.ExactArgs(2)),
	Run:     search,
}

func search(cmd *cobra.Command, args []string) {
	var err error
	engineType := args[0]
	results := []core.SearchResult{}

	query := core.Query{
		Text:  args[1],
		Limit: 10,
	}

	if appConf.IsRawRequests {
		results, err = searchRaw(engineType, query)
	} else {
		results, err = searchBrowser(engineType, query)
	}

	if err != nil {
		logrus.Error(err)
		return
	}

	b, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		logrus.Error(err)
		return
	}

	fmt.Println(string(b))
}
func searchBrowser(engineType string, query core.Query) ([]core.SearchResult, error) {
	var engine core.SearchEngine

	opts := core.BrowserOpts{
		IsHeadless: !appConf.IsBrowserHead, // Disable headless if browser head mode is set
		IsLeakless: appConf.IsLeakless,
		Timeout:    time.Second * time.Duration(appConf.Timeout),
	}

	//if appConf.IsDebug {
	//	opts.IsHeadless = false
	//}

	browser, err := core.NewBrowser(opts)
	if err != nil {
		logrus.Error(err)
	}

	switch strings.ToLower(engineType) {
	case "yandex":
		engine = yandex.New(*browser)
	case "google":
		engine = google.New(*browser)
	case "baidu":
		engine = baidu.New(*browser)
	default:
		logrus.Infof("No `%s` search engine found", engineType)
	}

	return engine.Search(query)
}

func searchRaw(engineType string, query core.Query) ([]core.SearchResult, error) {
	logrus.Warn("Browserless results are very inconsistent or may not even work!")

	switch strings.ToLower(engineType) {
	case "yandex":
		return yandex.Search(query)
	case "google":
		return google.Search(query)
	case "baidu":
		return baidu.Search(query)
	default:
		logrus.Infof("No `%s` search engine found", engineType)
	}
	return nil, nil
}

func init() {
	RootCmd.AddCommand(searchCMD)
}