package main

import "fmt"
import "io/ioutil"
import "encoding/json"
import "runtime"
import "os"
import "os/exec"
import "path"
import "errors"

type Config struct {
	Packages []Package `json: "packages"`
}

type Package struct {
	Name    string        `json: "name"`
	Version string        `json: "version"`
	Option  OptionPackage `json: "option"`
}

type OptionPackage struct {
	Windows WindowsOption `json: "windows"`
	Darwin  DarwinOption  `json: "darwin"`
	Linux   LinuxOption   `json: "linux"`
}

type WindowsOption struct {
	Type string `json: "type"`
}

type DarwinOption struct {
	Type string `json: "type"`
	Tap  string `json: "tap"`
}

type LinuxOption struct {
	Type string `json: "type"`
}

// os error
func osError(err error, exitCode int) {
	if err != nil {
		fmt.Println(err)
		os.Exit(exitCode)
	}
	return
}

// ディレクトリか確認
func isDir(d string) (bool, error) {
	var _isDir bool
	var err error
	fInfo, _ := os.Stat(d)
	_isDir = fInfo.IsDir()
	if !_isDir {
		err = errors.New("args don't directory path.")
	}
	return _isDir, err
}

func main() {
	fmt.Println("installing...")

	// 引数
	args := os.Args

	var configFile string
	var err error

	// config.json
	if len(args) == 1 {
		var curDir, _ = os.Getwd()
		configFile = path.Join(curDir, "config.json")
	} else {
		isDir, err := isDir(args[1])
		if isDir {
			configFile = path.Join(args[1], "config.json")
		} else {
			osError(err, 400)
		}
	}

	// ファイル読み込み
	file, err := ioutil.ReadFile(configFile)
	osError(err, 401)

	// JSON解析
	var data Config
	err = json.Unmarshal([]byte(string(file)), &data)
	osError(err, 402)

	// settingsファイルの置き場所
	var settingsFile string
	var versionManageDir string
	switch runtime.GOOS {
	case "windows":
		versionManageDir = path.Join(runtime.GOOS, "chocolatey")
		settingsFile = "packages.config"
	case "darwin":
		versionManageDir = path.Join(runtime.GOOS, "homebrew")
		settingsFile = "Brewfile"
	case "linux":
		versionManageDir = path.Join(runtime.GOOS, "linuxbrew")
		settingsFile = "Brewfile"
	default:
		osError(err, 403)
	}

	// それ用のフォルダを作成
	settingsDir := path.Dir(configFile)
	settingsDirPath := path.Join(settingsDir, versionManageDir)
	err = os.MkdirAll(settingsDirPath, 0777)
	osError(err, 404)

	// ファイル操作
	settingsFilePath := path.Join(settingsDirPath, settingsFile)
	_, err = os.Stat(settingsFilePath)
	if err != nil {
		err = ioutil.WriteFile(settingsFilePath, []byte(""), 0644)
		osError(err, 405)
	}
	file, err = ioutil.ReadFile(settingsFilePath)
	osError(err, 406)

	// ファイル書き込み
	if runtime.GOOS == "windows" {
		content := "<?xml version=\"1.0\"?>\n"
		content += "<packages>\n"
		for _, v := range data.Packages {
			content += "<package id=\"" + v.Name + "\" />\n"
		}
		content += "</packages>"
		ioutil.WriteFile(settingsFilePath, []byte(content), os.ModePerm)

		command := `@powershell 
					-NoProfile 
					-ExecutionPolicy 
					Bypass 
					-Command 
					"iex ((new-object net.webclient).DownloadString('https://chocolatey.org/install.ps1'))"
					&& 
					SET PATH=%PATH%;
					%ALLUSERSPROFILE%\chocolatey\bin`
		out, err := exec.Command("cmd", "/C", command).Output()
		osError(err, 407)

		var curDir, _ = os.Getwd()
		os.Chdir(settingsDirPath)
		out, err = exec.Command("cmd", "/C", "cinst -y "+settingsFilePath).Output()
		osError(err, 408)
		fmt.Println(string(out))
		os.Chdir(curDir)
		os.Exit(0)

	} else if runtime.GOOS == "darwin" {
		content := "cask_args appdir: '/Applications'\n"
		for _, v := range data.Packages {
			if v.Option.Darwin.Tap != "" {
				content += "tap '" + v.Option.Darwin.Tap + "'\n"
			}

			content += v.Option.Darwin.Type + " '" + v.Name + "'\n"
		}
		ioutil.WriteFile(settingsFilePath, []byte(content), os.ModePerm)
		_, err := exec.Command("ruby", "-v").Output()
		osError(err, 407)
		_, err = exec.Command("brew", "-v").Output()
		if err != nil {
			_, err = exec.Command("ruby", "-e", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)").Output()
			osError(err, 408)
		}
		_, err = exec.Command("brew", "tap", "Homebrew/bundle").Output()
		osError(err, 409)
		os.Chdir(settingsDirPath)
		out, err := exec.Command("brew", "bundle").Output()
		osError(err, 410)
		fmt.Printf(string(out))
	} else if runtime.GOOS == "linux" {
		content := "cask_args appdir: '/Applications'\n"
		for _, v := range data.Packages {
			if v.Option.Darwin.Tap != "" {
				content += "tap '" + v.Option.Darwin.Tap + "'\n"
			}

			content += v.Option.Darwin.Type + " '" + v.Name + "'\n"
		}
		ioutil.WriteFile(settingsFilePath, []byte(content), os.ModePerm)
		_, err := exec.Command("ruby", "-v").Output()
		osError(err, 407)
		_, err = exec.Command("brew", "-v").Output()
		if err != nil {
			_, err = exec.Command("ruby", "-e", "$(curl -fsSL https://raw.githubusercontent.com/Linuxbrew/linuxbrew/go/install)").Output()
			osError(err, 408)
		}
		_, err = exec.Command("brew", "tap", "Homebrew/bundle").Output()
		osError(err, 409)
		os.Chdir(settingsDirPath)
		out, err := exec.Command("brew", "bundle").Output()
		osError(err, 410)
		fmt.Printf(string(out))
	}
	fmt.Println("finished")
}
