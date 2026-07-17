package daemon

import (
	"os"
	"fmt"
	"path"
	"errors"
)

var defaultConf Config
func init() {
	var err error
	defaultConf, err = GetDefaultConfig()
	if err != nil {fmt.Fprintf(os.Stderr, "Could not get path to default config locations: %v", err)}
}

type Config struct {
	ReminderPath string
	CommandPath string
	ConfigPath string
}
const CONFIG_FOLDER_NAME string = "trem"
const CONFIG_REMINDER_NAME string = "reminders.gob"
const CONFIG_COMMAND_NAME string = "commands.gob"
const CONFIG_CONFIG_NAME string = "tremrc"

func GetDefaultConfig() (Config, error) {
	confdir, err := os.UserConfigDir()
	if err != nil {return Config{}, err}

	rems := path.Join(confdir, CONFIG_FOLDER_NAME, CONFIG_REMINDER_NAME)
	cmds := path.Join(confdir, CONFIG_FOLDER_NAME, CONFIG_COMMAND_NAME)
	conf := path.Join(confdir, CONFIG_FOLDER_NAME, CONFIG_CONFIG_NAME)

	return Config{
		ReminderPath: rems,
		CommandPath: cmds,
		ConfigPath: conf,
	}, nil
}

func GenerateDefaultConfig() []error {
	// Create the config folder
	confdir, err := os.UserConfigDir()
	if err != nil {return []error{errors.New("Could not get user config dir"), err}}
	err = os.MkdirAll(path.Join(confdir, CONFIG_FOLDER_NAME), 0700)
	if err != nil {return []error{errors.New("Could not generate config dir"), err}}

	// Create the files
	dconf, err := GetDefaultConfig()
	if err != nil {return []error{errors.New("Could not generate default config struct"), err}}

	names := []string{dconf.ConfigPath, dconf.CommandPath, dconf.ReminderPath}
	files := make([]*os.File, 3)
	var errs []error = nil
	for i, name := range names {
		file, err := os.Create(name)
		if err != nil {errs = append(errs, err)}
		files[i] = file
	}
	if errs != nil {return errs}

	errs = append(errs, writeDefaultConfig(files[0]))
	errs = append(errs, writeDefaultCommands(files[1]))
	errs = append(errs, writeDefaultReminders(files[2]))

	return errs

	// TODO: There's probably a better way of doing error checking here
}
func writeDefaultReminders(*os.File) error {
	return errors.New("TODO: Unimplemented")
}
func writeDefaultCommands(*os.File) error {
	return errors.New("TODO: Unimplemented")
}
func writeDefaultConfig(*os.File) error {
	return errors.New("TODO: Unimplemented")
}