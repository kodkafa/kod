package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"kodkafa/internal/app/usecases"
	"kodkafa/internal/infra/exec"
	"kodkafa/internal/infra/repo"
	"kodkafa/internal/infra/runtime"
	"kodkafa/internal/infra/store"
	"kodkafa/internal/ui"

	tea_pkg "github.com/charmbracelet/bubbletea"
)

// Run initializes dependencies and starts the TUI or executes CLI commands.
func Run() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not determine home directory:", err)
	}
	baseDir := filepath.Join(home, ".kodkafa")

	// 1. Initialize Essential Infrastructure
	configStore := store.NewConfigStore(baseDir)
	initUC := usecases.NewInitLayoutUseCase(baseDir, "default_config.json", configStore)

	// Always ensure layout on start
	if err := initUC.Execute(); err != nil {
		log.Fatal("Failed to initialize layout:", err)
	}

	pluginRepo := repo.NewPluginRepository(baseDir)
	usageStore := store.NewUsageStore(baseDir)
	stateStore := store.NewStateStore(baseDir)
	runner := exec.NewProcessRunner(baseDir)
	installer := runtime.NewFSInstaller(baseDir, configStore)

	// Initialize Use Cases
	listUC := usecases.NewListPluginsUseCase(pluginRepo, usageStore, configStore, stateStore)
	addUC := usecases.NewAddPluginUseCase(pluginRepo, stateStore, configStore)
	deleteUC := usecases.NewDeletePluginUseCase(pluginRepo, stateStore, usageStore, installer)
	loadUC := usecases.NewLoadPluginDepsUseCase(pluginRepo, installer)
	infoUC := usecases.NewGetPluginInfoUseCase(pluginRepo, stateStore)
	runUC := usecases.NewRunPluginUseCase(pluginRepo, stateStore, usageStore, configStore, runner)

	// 2. Dispatch CLI or TUI
	if len(os.Args) > 1 {
		handleCLI(os.Args[1:], initUC, addUC, deleteUC, loadUC, infoUC, runUC, listUC)
		return
	}

	// 3. Start TUI
	cfg, err := configStore.Read()
	showSplash := true
	if err == nil && cfg != nil {
		showSplash = cfg.Splash
	}

	rootModel := ui.NewModel(listUC, addUC, deleteUC, loadUC, infoUC, runUC, initUC, showSplash)
	p := tea_pkg.NewProgram(rootModel, tea_pkg.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func handleCLI(args []string, initUC *usecases.InitLayoutUseCase, addUC *usecases.AddPluginUseCase, deleteUC *usecases.DeletePluginUseCase, loadUC *usecases.LoadPluginDepsUseCase, infoUC *usecases.GetPluginInfoUseCase, runUC *usecases.RunPluginUseCase, listUC *usecases.ListPluginsUseCase) {
	cmd := args[0]
	switch cmd {
	case "init":
		if err := initUC.Execute(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("System initialized successfully.")
	case "add", "a":
		if len(args) < 2 {
			fmt.Println("Usage: kodkafa add <path|url>")
			os.Exit(1)
		}
		res, err := addUC.Execute(usecases.AddPluginInput{Source: args[1]})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Plugin added: %s\n", res.Plugin.Name)

		// Auto-load dependencies
		fmt.Printf("loading dependencies for %s...\n", res.Plugin.Name)
		loadRes, err := loadUC.Execute(usecases.LoadPluginDepsInput{PluginName: res.Plugin.Name})
		if err != nil {
			fmt.Printf("Warning: Failed to load dependencies: %v\n", err)
		} else {
			fmt.Printf("Dependencies loaded: %s\n", loadRes.Status)
		}
	case "del", "d":
		if len(args) < 2 {
			fmt.Println("Usage: kodkafa del <name>")
			os.Exit(1)
		}
		name := args[1]

		// 0. Check existence first
		_, err := infoUC.Execute(usecases.GetPluginInfoInput{PluginName: name})
		if err != nil {
			fmt.Printf("Error: Plugin '%s' not found.\n", name)
			return
		}

		fmt.Printf("REMOVE PLUGIN: %s\n", name)
		fmt.Printf("Note: Dependencies will be removed. Run 'kod load' to reinstall.\n")
		fmt.Print("Are you sure? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			fmt.Println("Operation cancelled.")
			return
		}

		fmt.Print("Remove dependencies as well? (y/N): ")
		var remDeps string
		fmt.Scanln(&remDeps)
		removeDeps := strings.ToLower(remDeps) == "y"

		res, err := deleteUC.Execute(usecases.DeletePluginInput{PluginName: name, RemoveDeps: removeDeps})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Success: %s\n", res.Message)
	case "run", "r":
		if len(args) < 2 {
			fmt.Println("Usage: kodkafa run <name>")
			os.Exit(1)
		}
		name := args[1]

		// Launch TUI for run
		rootModel := ui.NewModel(listUC, addUC, deleteUC, loadUC, infoUC, runUC, initUC, false)
		rootModel.StartRun(name)
		p := tea_pkg.NewProgram(rootModel, tea_pkg.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	case "info", "i":
		if len(args) < 2 {
			fmt.Println("Usage: kodkafa info <name>")
			os.Exit(1)
		}
		res, err := infoUC.Execute(usecases.GetPluginInfoInput{PluginName: args[1]})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Plugin: %s\nInterpreter: %s\nDescription: %s\n", res.Plugin.Name, res.Plugin.Interpreter, res.Plugin.Description)
	case "load", "l":
		if len(args) < 2 {
			fmt.Println("Usage: kodkafa load <name>")
			os.Exit(1)
		}
		fmt.Printf("Installing dependencies for %s...\n", args[1])
		res, err := loadUC.Execute(usecases.LoadPluginDepsInput{PluginName: args[1]})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Dependencies loaded for %s. Status: %s\n", args[1], res.Status)
	case "log":
		fmt.Println("Log view not implemented in CLI yet.")
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
