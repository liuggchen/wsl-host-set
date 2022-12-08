package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const hostPath = "C:/Windows/System32/drivers/etc/hosts"
const domainFile = "wsl_domain.conf"

func init() {
	if l := len(os.Args); l > 1 && isHelpArg(os.Args[1]) {
		fmt.Printf(`设置wsl的IP到windows宿主机的hosts文件中的小工具

使用方法：
	方法1：将需要配置的域名写入到同目录的 %s 文件中，运行将自动按行读取
	
	方法2：将域名追加到运行命令后
	例如： ./wsl-host-set.exe local-website.com buyaoqiao.tech liuggchen.com

`, domainFile)

		os.Exit(0)
	}
}
func isHelpArg(arg string) bool {
	helpArgs := []string{"help", "--help", "-help", "--h", "-h"}
	for _, helpArg := range helpArgs {
		if strings.EqualFold(helpArg, arg) {
			return true
		}
	}
	return false
}

func main() {
	var logFile = InitLogger()
	defer logFile.Close()

	var domains = ParseDomains()
	if len(domains) == 0 {
		return
	}
	log.Printf("添加 %v %d\n", domains, len(domains))
	var wslIp = GetWslIp()
	log.Printf("wsl ip => %s\n", wslIp)
	originHost := CleanHost(domains)
	log.Printf("\n---------原始数据---------\n%s\n---------------------", originHost)
	newHost := AppendHost(originHost, domains, wslIp)
	log.Printf("\n---------新数据---------\n%s\n---------------------", newHost)

	WriteHost(newHost)
}

func ParseDomains() []string {
	var m = make(map[string]struct{}, 0)
	for i, arg := range os.Args {
		if i > 0 {
			m[strings.ToLower(arg)] = struct{}{}
		}
	}
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	bts, err := os.ReadFile(path.Join(dir, domainFile))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("读取文件 %s 失败 %v", domainFile, err)
	} else if err == nil {
		for _, domain := range strings.Split(string(bts), "\n") {
			if domain = strings.ToLower(strings.TrimSpace(domain)); domain != "" {
				m[domain] = struct{}{}
			}
		}
	}
	var domains = make([]string, 0)
	for s, _ := range m {
		domains = append(domains, s)
	}
	return domains
}

func WriteHost(content string) {
	f, err := os.OpenFile(hostPath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("写host文件时，打开文件失败")
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalf("写host文件时，写入失败")
	}
}

func AppendHost(originHost string, domains []string, ip string) string {
	var c strings.Builder
	c.WriteString(originHost)
	for _, domain := range domains {
		c.WriteString(fmt.Sprintf("%s %s\n", ip, domain))
	}
	return c.String()
}

func CleanHost(domain []string) string {
	f, err := os.Open(hostPath)
	if err != nil {
		log.Fatalf("打开hosts文件失败")
	}
	defer f.Close()
	var r = bufio.NewReader(f)
	var fileContent strings.Builder
	for {
		lineBts, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("按行读取hosts文件失败")
		}
		if s := string(lineBts); !isDomainLine(domain, s) {
			fileContent.Write(lineBts)
			fileContent.WriteByte('\n')
		}
	}
	return fileContent.String()
}

func isDomainLine(domains []string, lineStr string) bool {
	lineStr = strings.ToLower(strings.TrimSpace(lineStr))
	for _, domain := range domains {
		if strings.HasSuffix(lineStr, strings.ToLower(domain)) {
			return true
		}
	}
	return false
}

func GetWslIp() string {
	cmd := exec.Command("wsl", "hostname", "-I")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("获取wsl ip 失败")
	}
	return strings.TrimSpace(string(out))
}

func InitLogger() *os.File {
	logFile, _ := os.OpenFile("wsl-host-set.log", os.O_APPEND|os.O_CREATE, 0600)
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	return logFile
}
