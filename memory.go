package main

import (
    "io/ioutil"
    "log"
    "encoding/xml"
    "os"
    "time"
)

const memoryFileName = "amigo-memory.xml"

// Memory stores all Amigo's knowledge.
// Persists itself to a fixed XML file which will be created automatically
// on first call into the directory from where the Amigo bot its executed
type Memory struct {
    // Connection params
	Host, Channel, Nick, Password string

    Masters []string
}

// LoadMemory creates and returns a new Memory instance and initializes it.
// If there is a previous saved Memory file, will load it as well.
func LoadMemory() *Memory {
    mem := new(Memory)

    raw, err := ioutil.ReadFile(memoryFileName)
    if err != nil {
        log.Println("AMIGO WARNING: Memory file not found. Loading factory memory.")
        mem.factoryMemory()
    } else {
        err = xml.Unmarshal(raw, mem)
        if err != nil {
            log.Println("AMIGO WARNING: Invalid memory file. Loading factory memory.")
            mem.factoryMemory()
        }
    }

    mem.persist()

    return mem
}

// factoryMemory sets default memory information needed to be operative.
func (m *Memory) factoryMemory() {
    m.Host = "irc.freenode.org:6667"
	m.Channel = "#amigo-bot"
	m.Nick = "amigobot"

    // Generate password based on host data
    var pass string
    pass = "eh,amigo!"

    hostname, err := os.Hostname()
    if err == nil {
        pass += hostname
    }

    pid := os.Getpid()
    pass += string(pid)

    wd, err := os.Getwd()
    if err == nil {
        pass += wd
    }

	m.Password = pass

    m.Write()
}

// Write saves memory to an XML file
func (m *Memory) Write() {
    raw, err := xml.MarshalIndent(m, "", "    ")
    if err != nil {
        log.Println("AMIGO ERROR: Memory data cannot be encoded to be saved.")
        return
    }

    err = ioutil.WriteFile(memoryFileName, append([]byte(xml.Header), raw...), os.ModePerm)
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
