package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Process struct {
	PID     string
	Command string
	State   string
}

func getProcesses() []Process {
	var processes []Process

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return processes
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if directory name is a PID (numeric)
		if _, err := strconv.Atoi(entry.Name()); err != nil {
			continue
		}

		pid := entry.Name()
		cmdPath := filepath.Join("/proc", pid, "comm")
		statusPath := filepath.Join("/proc", pid, "status")

		// Read command name
		cmdBytes, err := os.ReadFile(cmdPath)
		cmd := "unknown"
		if err == nil {
			cmd = strings.TrimSpace(string(cmdBytes))
		}

		// Read state from status
		state := "unknown"
		statusBytes, err := os.ReadFile(statusPath)
		if err == nil {
			lines := strings.Split(string(statusBytes), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "State:") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						state = parts[1]
					}
					break
				}
			}
		}

		processes = append(processes, Process{
			PID:     pid,
			Command: cmd,
			State:   state,
		})
	}

	return processes
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>topz - Process Monitor</title>
    <style>
        body { font-family: monospace; margin: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; }
        table { border-collapse: collapse; width: 100%; }
        th, td { padding: 8px 12px; text-align: left; border: 1px solid #333; }
        th { background: #333; }
        tr:hover { background: #2a2a2a; }
    </style>
    <meta http-equiv="refresh" content="2">
</head>
<body>
    <h1>topz - Process Monitor (Sidecar)</h1>
    <p>Processes visible in shared PID namespace:</p>
    <table>
        <tr><th>PID</th><th>Command</th><th>State</th></tr>
        {{range .}}
        <tr><td>{{.PID}}</td><td>{{.Command}}</td><td>{{.State}}</td></tr>
        {{end}}
    </table>
</body>
</html>
`

func topzHandler(w http.ResponseWriter, r *http.Request) {
	processes := getProcesses()
	tmpl := template.Must(template.New("topz").Parse(htmlTemplate))
	tmpl.Execute(w, processes)
}

func main() {
	addr := ":8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	http.HandleFunc("/topz", topzHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/topz", http.StatusFound)
	})

	fmt.Printf("topz sidecar serving on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
