package main

import "fmt"
import "io/ioutil"
import "encoding/json"
import "runtime"
import "os"
import "os/exec"
import "path/filepath"

type Config struct {
	Packages []Package `json: "packages"`
}

type Package struct {
	Name    string `json: "name"`
	Version string `json: "version"`
	Option OptionPackage `json: "option"`
}

type OptionPackage struct{
	Windows WindowsOption `json: "windows"`
	Darwin DarwinOption `json: "darwin"`
	Linux LinuxOption `json: "linux"`
}

type WindowsOption struct {
	Type string `json: "type"`
}

type DarwinOption struct {
	Type string `json: "type"`
	Tap string `json: "tap"`
}

type LinuxOption struct {
	Type string `json: "type"`
}

// ディレクトリか確認
func confirmDir(d string) (isDir bool) {
	fInfo, _ := os.Stat(d)
	return fInfo.IsDir()
}

func main() {
	// 引数
	args := os.Args

	var configFile string
	var err error

	// config.json
	if len(args) == 1 {
		var curDir, _ = os.Getwd()
		configFile = curDir + "/config.json"
	} else {
		if confirmDir(args[1]) {
			configFile = args[1] + "/config.json"
		} else {
			fmt.Printf("Error: args don't directory path.")
		}
	}

	// ファイル読み込み
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	// JSON解析
	var data Config
	err = json.Unmarshal([]byte(string(file)), &data)
	if err != nil {
		panic(err)
	}

	// settingsファイルの置き場所
	var settingsFile string
	var versionManageDir string
	switch runtime.GOOS {
	case "windows":
		versionManageDir = runtime.GOOS + "/chocolatey"
		settingsFile = "packages.config"
	case "darwin":
		versionManageDir = runtime.GOOS + "/homebrew"
		settingsFile = "Brewfile"
	case "linux":
		versionManageDir = runtime.GOOS + "/linuxbrew"
		settingsFile = "Brewfile"
	default:
		panic(err)
	}

	// それ用のフォルダを作成
	settingsDir := filepath.Dir(configFile)
	err = os.MkdirAll(settingsDir+"/"+versionManageDir, 0777)
	if err != nil {
		panic(err)
	}

	// ファイル操作
	settingsDirPath := settingsDir + "/" + versionManageDir
	settingsFilePath := settingsDirPath + "/" + settingsFile
	_, err = os.Stat(settingsFilePath)
	if err != nil {
		err = ioutil.WriteFile(settingsFilePath, []byte(""), 0644)
		if err != nil {
			panic(err)
		}
	}
	file, err = ioutil.ReadFile(settingsFilePath)
	if err != nil {
		panic(err)
	}

	// ファイル書き込み
	if runtime.GOOS == "windows"{
		content := "<?xml version=\"1.0\"?>"
		content += "<packages>"
		for _, v := range data.Packages {
			content += "<package id=\""+v.Name+"\" />\n"
		}
		content += "</packages>"
		ioutil.WriteFile(settingsFilePath, []byte(content), os.ModePerm)

	}else if runtime.GOOS == "darwin"{
		content := "cask_args appdir: '/Applications'\n"
		for _, v := range data.Packages {
			if v.Option.Darwin.Tap != ""{
				content += "tap '" + v.Option.Darwin.Tap + "'\n"
			}

			content += v.Option.Darwin.Type + " '"+v.Name + "'\n"
		}
		ioutil.WriteFile(settingsFilePath, []byte(content), os.ModePerm)
		_, err := exec.Command("ruby", "-v").Output()
		if err != nil {
            panic(err)
        }
		_, err = exec.Command("brew", "-v").Output()
		if err != nil {
			_, err = exec.Command("ruby", "-e", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)").Output()
			if err != nil {
	            panic(err)
	        }
        }
		_, err = exec.Command("brew", "tap", "Homebrew/bundle").Output()
		if err != nil {
			panic(err)
		}
		os.Chdir(settingsDirPath)
		out, err := exec.Command("brew", "bundle").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(out))
	}else if runtime.GOOS == "linux"{
		content := "cask_args appdir: '/Applications'\n"
		for _, v := range data.Packages {
			if v.Option.Darwin.Tap != ""{
				content += "tap '" + v.Option.Darwin.Tap + "'\n"
			}

			content += v.Option.Darwin.Type + " '"+v.Name + "'\n"
		}
		ioutil.WriteFile(settingsFilePath, []byte(content), os.ModePerm)
		_, err := exec.Command("ruby", "-v").Output()
		if err != nil {
            panic(err)
        }
		_, err = exec.Command("brew", "-v").Output()
		if err != nil {
			_, err = exec.Command("ruby", "-e", "$(curl -fsSL https://raw.githubusercontent.com/Linuxbrew/linuxbrew/go/install)").Output()
			if err != nil {
	            panic(err)
	        }
        }
		_, err = exec.Command("brew", "tap", "Homebrew/bundle").Output()
		if err != nil {
			panic(err)
		}
		os.Chdir(settingsDirPath)
		out, err := exec.Command("brew", "bundle").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(out))
	}
}
