package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"regexp"
	"strconv"
)

type Data struct {
	Name string `json:"name"`
	Ip1  string `json:"ip1"`
	Ip2  string `json:"ip2"`
}

type App struct {
	dnsMap map[string][]string
	active string
}

func NewApp() *App {
	a := &App{}
	a.dnsMap = a.getDNSList()
	return a
}

func (a *App) GetDNSList() map[string][]string {
	return a.dnsMap
}

func (a *App) GetActiveDNS() string {
	return a.active
}

func (a *App) SetDNS(name string) string {
	dns, ok := a.dnsMap[name]
	if !ok {
		return "DNS name not found"
	}
	a.clearCache()
	iface := getInterfaceName()
	if iface == "" {
		return "No suitable network interface found"
	}

	cmd := exec.Command("netsh", "interface", "ip", "set", "dns", fmt.Sprintf(`name="%s"`, iface), "static", dns[0])
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("❌ Primary DNS error: %v\n%s", err, string(out))
	}

	cmd2 := exec.Command("netsh", "interface", "ip", "add", "dns", fmt.Sprintf(`name="%s"`, iface), dns[1], "index=2")
	cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("❌ Secondary DNS error: %v\n%s", err, string(out2))
	}

	a.active = name
	return "✅ DNS updated successfully"
}

func (a *App) AddDNS(name, ip1, ip2 string) string {
	if name == "" || ip1 == "" || ip2 == "" {
		return "All fields are required"
	}

	a.dnsMap[name] = []string{ip1, ip2}

	// Append to dnslist.jsonl
	file, err := os.OpenFile("dnslist.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "Failed to open dnslist.jsonl"
	}
	defer file.Close()

	data := Data{Name: name, Ip1: ip1, Ip2: ip2}
	jsonbytes, err := json.Marshal(data)
	if err != nil {
		return "Failed to marshal new DNS entry"
	}
	_, err = file.Write(append(jsonbytes, '\n'))
	if err != nil {
		return "Failed to write DNS to file"
	}

	return "✅ New DNS added"
}

func (a *App) getDNSList() map[string][]string {
	file, err := os.Open("dnslist.jsonl")
	if err != nil {
		return a.writeNewList()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	myIPS := make(map[string][]string)
	var data Data

	for scanner.Scan() {
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &data); err == nil {
			myIPS[data.Name] = []string{data.Ip1, data.Ip2}
		}
	}

	if len(myIPS) == 0 {
		return a.writeNewList()
	}
	return myIPS
}


func (a *App) PingDNS() string {
	cmd := exec.Command("ping", "-n", "4", "4.2.2.4")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Ping error"
	}

	output := string(out)
	// دنبال میانگین latency در خط آخر می‌گردیم
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Average") {
			re := regexp.MustCompile(`Average = (\d+)ms`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 2 {
				avg, _ := strconv.Atoi(matches[1])
				return "Average ping: " + strconv.Itoa(avg) + " ms"
			}
		}
	}

	return "Could not read ping"
}

func (a *App) writeNewList() map[string][]string {
	file, err := os.Create("dnslist.jsonl")
	if err != nil {
		log.Fatal("could not create new file")
	}
	defer file.Close()

	myData := []Data{
		{Name: "Electero", Ip1: "78.157.42.101", Ip2: "78.157.42.100"},
		{Name: "Shecan", Ip1: "185.51.200.2", Ip2: "178.22.122.100"},
		{Name: "Radar", Ip1: "10.202.10.10", Ip2: "10.202.10.11"},
		{Name: "403", Ip1: "10.202.10.202", Ip2: "10.202.10.102"},
		{Name: "Begzar", Ip1: "185.55.226.26", Ip2: "185.55.225.25"},
		{Name: "Shelter", Ip1: "94.103.125.157", Ip2: "94.103.125.158"},
		{Name: "Beshkan", Ip1: "181.41.194.177", Ip2: "181.41.194.186"},
		{Name: "Pishgaman", Ip1: "5.202.100.100", Ip2: "5.202.100.101"},
		{Name: "Level3", Ip1: "209.244.0.3", Ip2: "209.244.0.4"},
		{Name: "Google", Ip1: "8.8.8.8", Ip2: "8.8.4.4"},
	}

	myIPS := make(map[string][]string)

	for _, data := range myData {
		myIPS[data.Name] = []string{data.Ip1, data.Ip2}
		jsonbytes, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = file.Write(append(jsonbytes, '\n'))
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return myIPS
}

func (a *App) clearCache() {
	cmd := exec.Command("ipconfig", "/flushdns")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ Error flushing DNS cache: %v\n%s", err, string(out))
	}
}

func getInterfaceName() string {
	network, err := net.Interfaces()
	if err != nil {
		log.Fatal("error in interfaces")
	}

	for _, iface := range network {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			lowerName := strings.ToLower(iface.Name)
			if strings.Contains(lowerName, "virtual") || strings.Contains(lowerName, "vethernet") || strings.Contains(lowerName, "docker") {
				continue
			}
			return iface.Name
		}
	}
	return ""
}
