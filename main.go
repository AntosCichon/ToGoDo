package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)

type Entry struct {
	Id       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Modifier int       `json:"modifier"`
	Color    string    `json:"color"`
}

func getEntries() ([]Entry, error) {
	data, err := os.ReadFile("list.json")
	if err != nil {
		fmt.Printf("\033[31mCannot read file list.json\033[0m\n%s\n", err)
		return nil, err
	} else if len(data) < 2 {
		// when file is less than 2 characters long, just ignore it and use empty json list: []
		data = []byte("[]")
	}
	var entries []Entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		fmt.Printf("\033[31mCannot Unmarshal json from list.json\033[0m\n%s\n", err)
		return nil, err
	}
	return entries, nil
}
func saveEntry(entry Entry) error {
	entries, err := getEntries()
	if err != nil {
		fmt.Printf("\033[31mCannot get Entries from file list.json033[0m\n%s\n", err)
		return err
	}

	entries = append(entries, entry)

	entriesJson, err := json.Marshal(entries)
	if err != nil {
		fmt.Printf("\033[31mCannot Marshal entries slice\033[0m\n%s\n", err)
		return err
	}

	err = os.WriteFile("list.json", entriesJson, 0644)
	if err != nil {
		fmt.Printf("\033[31mCannot save file list.json\033[0m\n%s\n", err)
		return err
	}

	return nil
}
func removeEntry(entries []Entry, id string) ([]Entry, error) {
	deleteId, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("\033[31mCannot parse string uuid to uuid.UUID\033[0m\n%s\n", err)
		return nil, err
	}
	newEntries := entries
	for i, entry := range entries {
		if entry.Id == deleteId {
			newEntries = append(entries[:i], entries[i+1:]...)
			break
		}
	}

	entriesJson, err := json.Marshal(newEntries)
	if err != nil {
		fmt.Printf("\033[31mCannot Marshal entries slice\033[0m\n%s\n", err)
		return nil, err
	}

	err = os.WriteFile("list.json", entriesJson, 0644)
	if err != nil {
		fmt.Printf("\033[31mCannot save file list.json\033[0m\n%s\n", err)
		return nil, err
	}

	return newEntries, nil
}
func changeModifier(entries []Entry, id string, modifier int) ([]Entry, error) {
	modifyId, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("\033[31mCannot parse string uuid to uuid.UUID\033[0m\n%s\n", err)
		return nil, err
	}
	newEntries := entries
	for i, entry := range entries {
		if entry.Id == modifyId {
			entries[i].Modifier = modifier
			break
		}
	}

	entriesJson, err := json.Marshal(newEntries)
	if err != nil {
		fmt.Printf("\033[31mCannot Marshal entries slice\033[0m\n%s\n", err)
		return nil, err
	}

	err = os.WriteFile("list.json", entriesJson, 0644)
	if err != nil {
		fmt.Printf("\033[31mCannot save file list.json\033[0m\n%s\n", err)
		return nil, err
	}

	return newEntries, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := getEntries()
	if err != nil {
		http.Error(w, "Failed to read ToDos file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com")
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, entries)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid POST data", http.StatusBadRequest)
		return
	}

	entryTitle := r.FormValue("title")
	if entryTitle == "" {
		http.Error(w, "Field \"title\" cannot be empty", http.StatusBadRequest)
		return
	}
	entryColor := r.FormValue("color")

	entry := Entry{
		Id:       uuid.New(),
		Title:    html.EscapeString(entryTitle),
		Modifier: 0,
		Color:    entryColor,
	}
	err = saveEntry(entry)
	if err != nil {
		http.Error(w, "Failed to save new entry", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, "Failed to Marshal response body", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid POST data", http.StatusBadRequest)
		return
	}

	deleteId := r.FormValue("id")
	if deleteId == "" {
		http.Error(w, "Field \"id\" cannot be empty", http.StatusBadRequest)
		return
	}

	entries, err := getEntries()
	if err != nil {
		http.Error(w, "Failed to read ToDos file", http.StatusInternalServerError)
		return
	}

	newEntries, err := removeEntry(entries, deleteId)
	if err != nil {
		http.Error(w, "Failed to remove entry from file", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(newEntries)
	if err != nil {
		http.Error(w, "Failed to Marshal response body", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func modifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid POST data", http.StatusBadRequest)
		return
	}

	modifyId := r.FormValue("id")
	if modifyId == "" {
		http.Error(w, "Field \"id\" cannot be empty", http.StatusBadRequest)
		return
	}
	newModifier, err := strconv.Atoi(r.FormValue("modifier"))
	if err != nil {
		http.Error(w, "Invalid modifier received", http.StatusBadRequest)
		return
	} else if newModifier != 0 && newModifier != 1 {
		http.Error(w, "Modifier has to be an integer 0 or 1", http.StatusBadRequest)
		return
	}

	entries, err := getEntries()
	if err != nil {
		http.Error(w, "Failed to read ToDos file", http.StatusInternalServerError)
		return
	}

	newEntries, err := changeModifier(entries, modifyId, newModifier)
	if err != nil {
		http.Error(w, "Failed to change modifier of this entry", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(newEntries)
	if err != nil {
		http.Error(w, "Failed to Marshal response body", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/remove", removeHandler)
	http.HandleFunc("/modify", modifyHandler)

	fmt.Printf("\033[32mServer open, listening on port 10004\033[0m\n")
	log.Fatal(http.ListenAndServe(":10004", nil))
}
