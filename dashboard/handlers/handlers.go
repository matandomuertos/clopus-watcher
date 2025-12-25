package handlers

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/kubeden/clopus-watcher/dashboard/db"
)

type Handler struct {
	db       *db.DB
	tmpl     *template.Template
	partials *template.Template
	logPath  string
}

func New(database *db.DB, tmpl, partials *template.Template, logPath string) *Handler {
	return &Handler{
		db:       database,
		tmpl:     tmpl,
		partials: partials,
		logPath:  logPath,
	}
}

type PageData struct {
	Fixes   []db.Fix
	Total   int
	Success int
	Failed  int
	Pending int
	Log     string
}

func (h *Handler) readLog() string {
	data, err := os.ReadFile(h.logPath)
	if err != nil {
		return "No watcher log available yet. Waiting for first run..."
	}
	// Get last 100 lines
	lines := strings.Split(string(data), "\n")
	if len(lines) > 100 {
		lines = lines[len(lines)-100:]
	}
	return strings.Join(lines, "\n")
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	fixes, err := h.db.GetFixes(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, success, failed, pending, err := h.db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Fixes:   fixes,
		Total:   total,
		Success: success,
		Failed:  failed,
		Pending: pending,
		Log:     h.readLog(),
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Fixes(w http.ResponseWriter, r *http.Request) {
	fixes, err := h.db.GetFixes(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, success, failed, pending, err := h.db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Fixes:   fixes,
		Total:   total,
		Success: success,
		Failed:  failed,
		Pending: pending,
	}

	err = h.partials.ExecuteTemplate(w, "fixes-table.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) Logs(w http.ResponseWriter, r *http.Request) {
	log := h.readLog()
	w.Header().Set("Content-Type", "text/html")
	// Escape HTML and preserve newlines
	escaped := template.HTMLEscapeString(log)
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")
	w.Write([]byte(escaped))
}
