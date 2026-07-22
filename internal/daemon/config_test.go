package daemon

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/realconebob/trem/internal"
)

func Test_GetDefaultConfig(t *testing.T) {
	conf, err := GetDefaultConfig()
	if err != nil {
		t.Error("Could not get default config", err)
	}

	fmt.Printf("conf: %v\n", conf)
}

func Test_GenerateDefaultConfig(t *testing.T) {
	errs := GenerateDefaultConfig()
	errs = misc.NilErrSliceCheck(errs)
	for _, err := range errs {
		t.Errorf("Encountered error while generating default config: %v", err)
	}
}

// TODO: Consider breaking into several other functions &/or inverting if statements for readability
func Test_WriteDefaults(t *testing.T) {
	// Test that writing the default config, reminders, and commands work as expected
	dir := t.TempDir()
	cmdpath := path.Join(dir, "commands.gob")
	rempath := path.Join(dir, "reminders.gob")
	cfgpath := path.Join(dir, "tremrc")

	if cmds, err := os.Create(cmdpath); err != nil {
		t.Errorf("Could not create file \"%v\" for writing: %v", cmdpath, err)
		return
	} else {
		defer cmds.Close()
		err = writeDefaultCommands(cmds)
		if err != nil {
			t.Errorf("Ran into an error while writing default command file: %v", err)
			return
		}

		// TODO: Do some sanity checking

	}

	if rems, err := os.Create(rempath); err != nil {
		t.Errorf("Could not create file \"%v\" for writing: %v", rempath, err)
		return
	} else {
		defer rems.Close()
		err = writeDefaultReminders(rems)
		if err != nil {
			t.Errorf("Ran into an error while writing default reminder file: %v", err)
			return
		}

		// TODO: Do some sanity checking
	}

	dconf, err := GetDefaultConfig()
	if err != nil {
		t.Errorf("Could not get default config for writing: %v", err)
		return
	}

	if cfg, err := os.Create(cfgpath); err != nil {
		t.Errorf("Could not create file \"%v\" for writing: %v", cfgpath, err)
		return
	} else {
		defer cfg.Close()
		err = writeDefaultConfig(cfg, dconf)
		if err != nil {
			t.Errorf("Ran into error while writing default config file: %v", err)
			return
		}

		conf, err := GetConfigFromFile(cfgpath)
		if err != nil {
			t.Errorf("Could not get config from file \"%v\": %v", cfgpath, err)
			return
		}
		if !conf.Compare(dconf) {t.Errorf("conf from file doesn't match default config")}
	}
}

func Test_GetConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	p := path.Join(dir, "tremrc")
	file, err := os.Create(p)
	if err != nil {
		t.Errorf("Could not get file handle to write example config")
	}

	cmds := path.Join(dir, "commands.gob")
	rems := path.Join(dir, "reminders.gob")
	toWrite := fmt.Appendf(nil,
		"[trem]\n" +
		"config = %v\n" +
		"                    commands = %v\n\n\n" +
		"   reminders      =             %v        \n",
		p, cmds, rems,
	) // Extra spaces to check trimming behavior, not that you'd ever write a config like this

	if _, err := file.Write(toWrite); err != nil {
		t.Error("Could not write example config contents to testing file")
		file.Close()
		return
	}
	file.Close()

	conf, err := GetConfigFromFile(p)
	if err != nil {
		t.Errorf("Encountered error in getting config from file \"%v\"\n", p)
	}

	if conf.ConfigPath != p {t.Errorf("Actual config path does not match expected. Expected: %v, Got: %v", p, conf.ConfigPath)}
	if conf.CommandPath != cmds {t.Errorf("Actual command path does not match expected. Expected: %v, Got: %v", cmds, conf.CommandPath)}
	if conf.ReminderPath != rems {t.Errorf("Actual reminder path does not match expected. Expected: %v, Got: %v", rems, conf.ReminderPath)}

	fmt.Print(conf)
}