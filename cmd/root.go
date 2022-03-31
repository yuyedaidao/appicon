/*
Copyright © 2022 yuyedaidao <wyqpadding@gmail.com>

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

	"github.com/nfnt/resize"
	"github.com/spf13/cobra"
)
const appiconsetContents string = `
{
    "images" : [
      {
        "filename" : "icon-20@2x.png",
        "idiom" : "iphone",
        "scale" : "2x",
        "size" : "20x20"
      },
      {
        "filename" : "icon-20@3x.png",
        "idiom" : "iphone",
        "scale" : "3x",
        "size" : "20x20"
      },
      {
        "filename" : "icon-29.png",
        "idiom" : "iphone",
        "scale" : "1x",
        "size" : "29x29"
      },
      {
        "filename" : "icon-29@2x.png",
        "idiom" : "iphone",
        "scale" : "2x",
        "size" : "29x29"
      },
      {
        "filename" : "icon-29@3x.png",
        "idiom" : "iphone",
        "scale" : "3x",
        "size" : "29x29"
      },
      {
        "filename" : "icon-40@2x.png",
        "idiom" : "iphone",
        "scale" : "2x",
        "size" : "40x40"
      },
      {
        "filename" : "icon-40@3x.png",
        "idiom" : "iphone",
        "scale" : "3x",
        "size" : "40x40"
      },
      {
        "filename" : "icon-60@2x.png",
        "idiom" : "iphone",
        "scale" : "2x",
        "size" : "60x60"
      },
      {
        "filename" : "icon-60@3x.png",
        "idiom" : "iphone",
        "scale" : "3x",
        "size" : "60x60"
      },
      {
        "filename" : "icon-20-ipad.png",
        "idiom" : "ipad",
        "scale" : "1x",
        "size" : "20x20"
      },
      {
        "filename" : "icon-20@2x-ipad.png",
        "idiom" : "ipad",
        "scale" : "2x",
        "size" : "20x20"
      },
      {
        "filename" : "icon-29-ipad.png",
        "idiom" : "ipad",
        "scale" : "1x",
        "size" : "29x29"
      },
      {
        "filename" : "icon-29@2x-ipad.png",
        "idiom" : "ipad",
        "scale" : "2x",
        "size" : "29x29"
      },
      {
        "filename" : "icon-40.png",
        "idiom" : "ipad",
        "scale" : "1x",
        "size" : "40x40"
      },
      {
        "filename" : "icon-40@2x.png",
        "idiom" : "ipad",
        "scale" : "2x",
        "size" : "40x40"
      },
      {
        "filename" : "icon-76.png",
        "idiom" : "ipad",
        "scale" : "1x",
        "size" : "76x76"
      },
      {
        "filename" : "icon-76@2x.png",
        "idiom" : "ipad",
        "scale" : "2x",
        "size" : "76x76"
      },
      {
        "filename" : "icon-83.5@2x.png",
        "idiom" : "ipad",
        "scale" : "2x",
        "size" : "83.5x83.5"
      },
      {
        "filename" : "icon-1024.png",
        "idiom" : "ios-marketing",
        "scale" : "1x",
        "size" : "1024x1024"
      }
    ],
    "info" : {
      "author" : "appicon",
      "version" : 1
    }
  }
`
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
		// 删除已存在的AppIcon.xcassets
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			output = "./"
		}
		if output == "" {
			output = "./"
		}
		if !strings.HasSuffix(output, "/") {
			output = output + "/"
		}
		appiconset := output + "AppIcon.appiconset"
		os.RemoveAll(appiconset)
		fmt.Printf("make directory: %s\n", appiconset)
		if err := os.Mkdir(appiconset, os.ModePerm); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
		data := []byte(appiconsetContents)
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
			size := item["size"].(string)
			width, _ := strconv.ParseFloat(strings.Split(size, "x")[0], 64)
			scale := item["scale"].(string)
			multiple, _ := strconv.Atoi(scale[:len(scale) - 1])
			realSize := uint(width * float64(multiple))
			newImage := resize.Resize(realSize, realSize, img, resize.Lanczos3)
			filePath := appiconset + "/" + filename
			fmt.Printf("create image %s\n", filename)
			out, err := os.Create(filePath)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
				return
			}
			if err := png.Encode(out, newImage); err != nil {
				out.Close()
				cmd.PrintErrln(err)
				os.Exit(1)
				return
			}
			out.Close()
		}
		fmt.Println("copy Contents.json to AppIcon.appiconset")
		jsonPath := output + "AppIcon.appiconset/Contents.json"
		if err := ioutil.WriteFile(jsonPath, data, 0644); err != nil {
			cmd.PrintErrln(err)
				os.Exit(1)
				return
		}
		fmt.Println("created appiconset successfuly")
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
