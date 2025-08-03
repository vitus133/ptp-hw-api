package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// NewPluginManager creates a new plugin manager and loads all plugins from the plugins directory
func NewPluginManager(pluginsDir string) (*PluginManager, error) {
	pm := &PluginManager{
		plugins: make(map[string]*HardwarePluginConfig),
	}

	// Load all plugin files from the plugins directory
	if err := pm.LoadPlugins(pluginsDir); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %w", err)
	}

	return pm, nil
}

// LoadPlugins loads all YAML plugin files from the specified directory
func (pm *PluginManager) LoadPlugins(pluginsDir string) error {
	// Check if plugins directory exists
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		// No plugins directory - not an error, just no hardware defaults available
		return nil
	}

	// Read all files in the plugins directory
	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory %s: %w", pluginsDir, err)
	}

	// Load each YAML file as a plugin
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml") {
			pluginPath := filepath.Join(pluginsDir, file.Name())
			if err := pm.LoadPlugin(pluginPath); err != nil {
				return fmt.Errorf("failed to load plugin %s: %w", pluginPath, err)
			}
		}
	}

	return nil
}

// LoadPlugin loads a single plugin file
func (pm *PluginManager) LoadPlugin(pluginPath string) error {
	data, err := ioutil.ReadFile(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin file: %w", err)
	}

	var plugin HardwarePluginConfig
	if err := yaml.Unmarshal(data, &plugin); err != nil {
		return fmt.Errorf("failed to parse plugin YAML: %w", err)
	}

	// Validate plugin has required fields
	if plugin.PluginInfo.Name == "" {
		return fmt.Errorf("plugin must have a name")
	}

	// Store plugin by name
	pm.plugins[plugin.PluginInfo.Name] = &plugin
	return nil
}

// GetPlugin returns a plugin by name, or nil if not found
func (pm *PluginManager) GetPlugin(name string) *HardwarePluginConfig {
	return pm.plugins[name]
}

// ListPlugins returns a list of all loaded plugin names
func (pm *PluginManager) ListPlugins() []string {
	var names []string
	for name := range pm.plugins {
		names = append(names, name)
	}
	return names
}

// ApplyPluginDefaults applies hardware plugin defaults to a condition's desired states
// This function merges plugin defaults with user-specified desired states
func (pm *PluginManager) ApplyPluginDefaults(clockChain *ClockChain, condition *Condition) error {
	// Only apply defaults for "default" condition type
	if len(condition.Sources) == 0 || condition.Sources[0].ConditionType != "default" {
		return nil
	}

	// If user has already specified desired states, we'll merge with plugin defaults
	// Build a map of existing desired states by clockID+boardLabel for quick lookup
	existingStates := make(map[string]*DesiredState)
	for i := range condition.DesiredStates {
		key := condition.DesiredStates[i].ClockID + ":" + condition.DesiredStates[i].BoardLabel
		existingStates[key] = &condition.DesiredStates[i]
	}

	// Process each subsystem and apply plugin defaults
	for _, subsystem := range clockChain.Structure {
		if subsystem.HardwarePlugin == "" {
			continue // No plugin specified, skip
		}

		plugin := pm.GetPlugin(subsystem.HardwarePlugin)
		if plugin == nil {
			// Plugin not found - this might be a warning, but not an error
			continue
		}

		// Apply defaults for all pins in this subsystem
		if err := pm.applySubsystemDefaults(subsystem, plugin, existingStates, &condition.DesiredStates); err != nil {
			return fmt.Errorf("failed to apply defaults for subsystem %s: %w", subsystem.Name, err)
		}
	}

	return nil
}

// applySubsystemDefaults applies plugin defaults for a specific subsystem
func (pm *PluginManager) applySubsystemDefaults(
	subsystem Subsystem,
	plugin *HardwarePluginConfig,
	existingStates map[string]*DesiredState,
	desiredStates *[]DesiredState,
) error {
	// Helper function to apply defaults for a pin if not already specified by user
	applyPinDefaults := func(boardLabel string) {
		key := subsystem.DPLL.ClockID + ":" + boardLabel

		// Skip if user has already specified this pin
		if existingStates[key] != nil {
			return
		}

		// Create new desired state with plugin defaults
		newState := DesiredState{
			ClockID:    subsystem.DPLL.ClockID,
			BoardLabel: boardLabel,
		}

		// Check for specific defaults for this pin
		if specificDefaults, exists := plugin.SpecificDefaults[boardLabel]; exists {
			if specificDefaults.EEC != nil {
				newState.EEC = &PinState{
					Priority: specificDefaults.EEC.Priority,
					State:    specificDefaults.EEC.State,
				}
			}
			if specificDefaults.PPS != nil {
				newState.PPS = &PinState{
					Priority: specificDefaults.PPS.Priority,
					State:    specificDefaults.PPS.State,
				}
			}
		} else {
			// No specific defaults for this pin - skip it
			return
		}

		// Add the new state if we created any pin configurations
		if newState.EEC != nil || newState.PPS != nil {
			*desiredStates = append(*desiredStates, newState)
			existingStates[key] = &(*desiredStates)[len(*desiredStates)-1]
		}
	}

	// Apply defaults for all pins in the subsystem
	for boardLabel := range subsystem.DPLL.PhaseInputs {
		applyPinDefaults(boardLabel)
	}
	for boardLabel := range subsystem.DPLL.PhaseOutputs {
		applyPinDefaults(boardLabel)
	}
	for boardLabel := range subsystem.DPLL.FrequencyInputs {
		applyPinDefaults(boardLabel)
	}
	for boardLabel := range subsystem.DPLL.FrequencyOutputs {
		applyPinDefaults(boardLabel)
	}

	return nil
}

// MergeUserConfigWithDefaults merges user-provided configuration with plugin defaults
// This is the main entry point for applying plugin defaults to a clock chain configuration
func (pm *PluginManager) MergeUserConfigWithDefaults(clockChain *ClockChain) error {
	if clockChain.Behavior == nil {
		return nil // No behavior section, nothing to merge
	}

	// Apply plugin defaults to each condition
	for i := range clockChain.Behavior.Conditions {
		if err := pm.ApplyPluginDefaults(clockChain, &clockChain.Behavior.Conditions[i]); err != nil {
			return fmt.Errorf("failed to apply plugin defaults to condition %s: %w",
				clockChain.Behavior.Conditions[i].Name, err)
		}
	}

	return nil
}
