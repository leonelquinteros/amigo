package amigo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const memoryFileName = "amigo-memory.json"

// Memory stores all Amigo's knowledge.
// Persists itself to a fixed JSON file which will be created automatically
// on first call into the directory from where the Amigo bot its executed
type Memory struct {
	// Master nicks
	Masters []string

	// Commands definitions
	Commands map[string]string

	// 'exec when' command configuration
	AutoCmd map[string]string
}

// LoadMemory creates and returns a new Memory instance and initializes it.
// If there is a previous saved Memory file, will load it as well.
func LoadMemory() *Memory {
	mem := new(Memory)

	raw, err := ioutil.ReadFile(memoryFileName)
	if err != nil {
		log.Println("AMIGO WARNING: Memory file not found. Loading empty memory.")
	} else {
		err = json.Unmarshal(raw, mem)
		if err != nil {
			log.Println("AMIGO WARNING: Invalid memory file. Loading empty memory.")
		}
	}

	// Init
	if mem.Commands == nil {
		mem.Commands = make(map[string]string)
	}
	if mem.AutoCmd == nil {
		mem.AutoCmd = make(map[string]string)
	}

	mem.persist()

	return mem
}

// Write saves memory to an JSON file
func (m *Memory) Write() {
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println("AMIGO ERROR: Memory data cannot be encoded to be saved: " + err.Error())
		return
	}

	err = ioutil.WriteFile(memoryFileName, raw, os.ModePerm)
	if err != nil {
		log.Println("AMIGO ERROR: Cannot write memory file, please check for permissions.")
		return
	}
}

// persist writes the memory to the memory file every 10 seconds.
func (m *Memory) persist() {
	log.Println("AMIGO: Writing memory to " + memoryFileName)

	go func() {
		for {
			time.Sleep(10000 * time.Millisecond)
			m.Write()
		}
	}()

}
