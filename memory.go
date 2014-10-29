package main

import (
    "io/ioutil"
    "log"
    "encoding/xml"
    "os"
)

const memoryFileName = "amigo-memory.xml"

type Memory struct {
    // Connection params
	Host, Channel, Nick, Master string

    Identities []string
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

    return mem
}

// factoryMemory sets default memory information needed to be operative.
func (m *Memory) factoryMemory() {
    m.Host = "irc.freenode.org:6667"
	m.Channel = "#amigo-bot"
	m.Nick = "amigobot"
	m.Master = "amigo-master"

    m.Identities = []string{"eh amigo!"}

    m.Write()
}

// Write saves memory to an XML file
func (m *Memory) Write() {
    raw, err := xml.MarshalIndent(m, "", "    ")
    if err != nil {
        log.Println("AMIGO ERROR: Memory data cannot be encoded to be saved.")
        return
    }

    log.Println("AMIGO: Writing memory to " + memoryFileName)

    err = ioutil.WriteFile(memoryFileName, append([]byte(xml.Header), raw...), os.ModePerm)
    if err != nil {
        log.Println("AMIGO ERROR: Cannot write memory file, please check for permissions.")
        return
    }
}
