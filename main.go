package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Page struct {
	Title      string
	Body       []byte
	Contains   string
	NotPresent string
	Warning    string
}

type HAProxyIPsFile struct {
	Contains   map[string]bool
	NotPresent map[string]bool
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Index"}
	renderTemplate(w, "index", p)
}

func updateAllowedIPsHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Update HAProxy Allowed IPs"}
	renderTemplate(w, "updateallowedips", p)
}

func updateAllowedIPsProcessHandler(w http.ResponseWriter, r *http.Request) {
	replacer := strings.NewReplacer(" ", "", "-", "", "	", "")
	var oldIPs HAProxyIPsFile
	oldIPs.Contains = make(map[string]bool, 0)
	oldIPs.NotPresent = make(map[string]bool, 0)
	var newIPs HAProxyIPsFile
	newIPs.Contains = make(map[string]bool, 0)
	newIPs.NotPresent = make(map[string]bool, 0)
	var resultIPs HAProxyIPsFile
	resultIPs.Contains = make(map[string]bool, 0)
	resultIPs.NotPresent = make(map[string]bool, 0)

	for _, line := range strings.Split(strings.TrimSuffix(r.FormValue("oldcontains"), "\n"), "\n") {
		oldIPs.Contains[replacer.Replace(line)] = false
	}
	for _, line := range strings.Split(strings.TrimSuffix(r.FormValue("oldnotpresent"), "\n"), "\n") {
		oldIPs.NotPresent[replacer.Replace(line)] = false
	}
	for _, line := range strings.Split(strings.TrimSuffix(r.FormValue("newcontains"), "\n"), "\n") {
		newIPs.Contains[replacer.Replace(line)] = false
		resultIPs.Contains[replacer.Replace(line)] = true
	}
	for _, line := range strings.Split(strings.TrimSuffix(r.FormValue("newnotpresent"), "\n"), "\n") {
		newIPs.NotPresent[replacer.Replace(line)] = false
		resultIPs.NotPresent[replacer.Replace(line)] = true
	}

	for ip := range oldIPs.Contains {
		if _, ok := resultIPs.Contains[ip]; !ok {
			resultIPs.NotPresent[ip] = ok
		}
	}

	for ip := range oldIPs.NotPresent {
		if _, ok := resultIPs.Contains[ip]; !ok {
			resultIPs.NotPresent[ip] = ok
		}
	}

	resultContains := ""
	for ip := range resultIPs.Contains {
		resultContains += ip + "\n"
	}

	resultNotPresent := ""
	for ip := range resultIPs.NotPresent {
		resultNotPresent += ip + "\n"
	}

	warning := ""
	for ip := range resultIPs.NotPresent {
		if _, ok := resultIPs.Contains[ip]; ok {
			warning = "WARNING: There are IPs both in contains and not present"
		}
	}

	p := &Page{Title: "Result for Allowed IPs", Contains: resultContains, NotPresent: resultNotPresent, Warning: warning}
	renderTemplate(w, "updateallowedipsresult", p)
}

func findDuplicatesHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "HAProxy Find Duplicate IPs"}
	renderTemplate(w, "findduplicates", p)
}

func findDuplicatesProcessHandler(w http.ResponseWriter, r *http.Request) {
	replacer := strings.NewReplacer(" ", "", "-", "", "	", "")
	var listips HAProxyIPsFile
	listips.Contains = make(map[string]bool, 0)

	for _, line := range strings.Split(strings.TrimSuffix(r.FormValue("iplist"), "\n"), "\n") {
		if _, ok := listips.Contains[replacer.Replace(line)]; ok {
			listips.Contains[replacer.Replace(line)] = ok
		} else {
			listips.Contains[replacer.Replace(line)] = false
		}
	}

	resultDuplicates := ""
	for ip := range listips.Contains {
		if listips.Contains[ip] {
			resultDuplicates += ip + "\n"
		}
	}

	resultUnique := ""
	for ip := range listips.Contains {
		resultUnique += ip + "\n"
	}

	warning := ""

	p := &Page{Title: "Result for Allowed IPs", Contains: resultDuplicates, NotPresent: resultUnique, Warning: warning}
	renderTemplate(w, "findduplicatesresult", p)
}

var templates = template.Must(template.ParseFiles("pages/index.html", "pages/updateallowedips.html", "pages/updateallowedipsresult.html", "pages/findduplicates.html", "pages/findduplicatesresult.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods(http.MethodGet)
	r.HandleFunc("/updateAllowedIPs/", updateAllowedIPsHandler).Methods(http.MethodGet)
	r.HandleFunc("/updateAllowedIPs/process", updateAllowedIPsProcessHandler).Methods(http.MethodPost)
	r.HandleFunc("/findDuplicates/", findDuplicatesHandler).Methods(http.MethodGet)
	r.HandleFunc("/findDuplicates/process", findDuplicatesProcessHandler).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
