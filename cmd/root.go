/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type AIValue map[string]interface{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "appicon",
	Short: "iOS AppIcon 生成器",
	Long:  `通过一张尺寸为 1024*1024 的图片制作iOS应用的 AppIcon`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		if path == "" {
			cmd.PrintErrln(errors.New("no such file or directory"))
			os.Exit(1)
			return
		}
		file, err := os.Open(path)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		img, err := png.Decode(file)
		file.Close()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		bounds := img.Bounds()
		if (bounds.Max != image.Point{1024, 1024}) {
			cmd.PrintErrln(errors.New("The size of the image must be 1024 * 1024"))
			os.Exit(1)
			return
		}
		data, err := ioutil.ReadFile("cmd/Contents.json")
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		var contents map[string]interface{}
		if err := json.Unmarshal(data, &contents); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		images := contents["images"].([]interface{})
		for _, image := range images {
			item := image.(map[string]interface{})
			filename := item["filename"].(string)
			fmt.Printf(filename)
			size := item["size"].(string)
			width, _ := strconv.ParseFloat(strings.Split(size, "x")[0], 64)
			fmt.Printf(" %f ", width)
			scale := item["scale"].(string)
			multiple, _ := strconv.Atoi(scale[:len(scale) - 1])
			fmt.Println(multiple)
			realSize := uint(width * float64(multiple))
			fmt.Println(realSize)
		}
		// resize.Resize()
		

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appicon.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "version information")
	rootCmd.Flags().StringP("path", "p", "", "path of image")
	rootCmd.Flags().StringP("output", "o", "", "output path of AppIcon.appiconset  (default current directory)")
}
