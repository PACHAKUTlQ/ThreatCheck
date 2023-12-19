package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"gopkg.in/yaml.v3"
)

type ScannerConfig struct {
	Name      string `yaml:"name"`
	ScanCmd   string `yaml:"scanCmd"`
	PosOutput string `yaml:"posOutput"`
	// NegOutput string `yaml:"negOutput"`
}

type Config struct {
	Scanners []ScannerConfig `yaml:"scanners"`
}

func parseConfig() (*Config, error) {
	path := "config.yaml"
	configInit(path)
	var config Config
	data, err := os.ReadFile(path) // Use os.ReadFile instead of ioutil.ReadFile
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func configInit(path string) {
	configTemplate := `
scanners:
  # - name: "AV name"
  #   scanCmd: "Command for scanning the target file. Use {{file}} as the file name to be scanned. The scanner executable is STRONGLY RECOMMENDED to be in PATH."
  #   posOutput: "A string in output of positive detection but not in negative"
  - name: "ESET"
    scanCmd: "ecls /clean-mode=none /no-quarantine {{file}}"
    posOutput: ">"
  - name: "Windows Defender"
    scanCmd: "MpCmdRun.exe -Scan -ScanType 3 -File {{file}} -DisableRemediation -Trace -Level 0x10"
    posOutput: "Threat information"
  # - name: "Any others"`

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If file doesn't exist, create it
		err := os.WriteFile(path, []byte(configTemplate), 0644)
		if err != nil {
			fmt.Println("create config.yaml error", err)
		}
	}
}

// enableVirtualTerminalProcessing enables virtual terminal processing for the given file descriptor.
func enableVirtualTerminalProcessing(fd uintptr) error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	var mode uint32
	handle := syscall.Handle(fd)

	// Get the current console mode
	r1, _, e1 := syscall.SyscallN(procGetConsoleMode.Addr(), uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if r1 == 0 {
		return os.NewSyscallError("GetConsoleMode", e1)
	}

	// Enable virtual terminal processing
	const enableVirtualTerminalProcessing uint32 = 0x0004
	mode |= enableVirtualTerminalProcessing

	r1, _, e1 = syscall.SyscallN(procSetConsoleMode.Addr(), uintptr(handle), uintptr(mode))
	if r1 == 0 {
		return os.NewSyscallError("SetConsoleMode", e1)
	}

	return nil
}
