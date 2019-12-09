package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/cheynewallace/tabby"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/foundriesio/fioctl/client"
)

var targetShowCmd = &cobra.Command{
	Use:   "show <version>",
	Short: "Show details of a specific target.",
	Run:   doTargetsShow,
	Args:  cobra.ExactArgs(1),
}

func init() {
	targetsCmd.AddCommand(targetShowCmd)
}

func doTargetsShow(cmd *cobra.Command, args []string) {
	factory := viper.GetString("factory")
	logrus.Debugf("Showing target for %s %s", factory, args[0])

	targets, err := api.TargetsList(factory)
	if err != nil {
		fmt.Print("ERROR: ")
		fmt.Println(err)
		os.Exit(1)
	}

	hashes := make(map[string]string)
	var tags []string
	var apps map[string]client.DockerApp
	for _, target := range targets.Signed.Targets {
		custom, err := api.TargetCustom(target)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		if custom.Version != args[0] {
			continue
		}
		if custom.TargetFormat != "OSTREE" {
			logrus.Debugf("Skipping non-ostree target: %v", target)
			continue
		}
		for _, hwid := range custom.HardwareIds {
			hashes[hwid] = hex.EncodeToString(target.Hashes["sha256"])
		}
		apps = custom.DockerApps
		tags = custom.Tags
	}

	fmt.Printf("Tags:\t%s\n\n", strings.Join(tags, ","))

	t := tabby.New()
	t.AddHeader("HARDWARE ID", "OSTREE HASH - SHA256")
	for name, val := range hashes {
		t.AddLine(name, val)
	}
	t.Print()

	fmt.Println()

	t = tabby.New()
	t.AddHeader("DOCKER APP", "VERSION")
	for name, app := range apps {
		t.AddLine(name, app.FileName)
	}
	t.Print()
}
