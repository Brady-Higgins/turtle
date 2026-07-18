package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/adrg/xdg"
)

type State struct {
	TotalTime    string
	TotalSavings string
	LastUse      string
}

// UpdateStateStart : runs on local
func UpdateStateStart(state *State) {
	state.LastUse = time.Now().Format(time.RFC3339)
}

// UpdateStateStop : runs on cloud
func UpdateStateStop(state *State, hourly float64) (string, error) {
	lastTime, err := time.Parse(time.RFC3339, state.LastUse)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if state.TotalTime == "" {
		state.TotalTime = "0.0"
	}
	if state.TotalSavings == "" {
		state.TotalSavings = "0.0"
	}

	prevTotal, err := strconv.ParseFloat(state.TotalTime, 64)
	diff := time.Now().Sub(lastTime).Hours()
	total := prevTotal + diff
	state.TotalTime = strconv.FormatFloat(total, 'f', 6, 64)
	state.LastUse = ""

	// calculate cost
	prevCost, err := strconv.ParseFloat(state.TotalSavings, 64)
	state.TotalSavings = strconv.FormatFloat(prevCost+diff*hourly, 'f', 6, 64)
	return strconv.FormatFloat(diff*hourly, 'f', 6, 64), nil
}

// ReadStateFile : reads state file into a State struct
func ReadStateFile() (*State, error) {
	statePath, err := xdg.StateFile("turtle/state.json")
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(statePath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	fileBytes, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}
	s := &State{}
	err = json.Unmarshal(fileBytes, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// WriteStateFile : write state struct to state file
func WriteStateFile(state *State) error {
	statePath, err := xdg.StateFile("turtle/state.json")
	if err != nil {
		return err
	}
	//fullPath := filepath.Join(statePath, "state.json")
	//fileBytes, err := json.Marshal(state)
	//// create state file if it doesn't exist
	//if _, err := os.Stat(statePath); errors.Is(err, os.ErrNotExist) {
	//	// create folder
	//	// 0755 = rwxrwxrwx
	//	err = os.MkdirAll(statePath, 0755)
	//	if err != nil {
	//		return err
	//	}
	//}
	//var file *os.File
	//if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
	//	// 0755 = rwxrwxrwx
	//	file, err = os.Create(fullPath)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	file, err = os.Open(fullPath)
	//}
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return err
	}
	fileBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, fileBytes, 0644)

	//if err != nil || file == nil {
	//	fmt.Println("heyb")
	//	return err
	//}
	//defer file.Close()
	//_, err = file.Write(fileBytes)
	//if err != nil {
	//	return err
	//}
	//return nil
}
