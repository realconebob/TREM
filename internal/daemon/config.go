package daemon

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/realconebob/trem/internal"
	"github.com/realconebob/trem/internal/reminders"
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
		if err == nil {defer file.Close()
		} else {errs = append(errs, err)}

		files[i] = file
	}
	if es := misc.NilErrSliceCheck(errs); es != nil {return es}

	errs = append(errs, writeDefaultConfig(files[0]))
	errs = append(errs, writeDefaultCommands(files[1]))
	errs = append(errs, writeDefaultReminders(files[2]))

	return misc.NilErrSliceCheck(errs)
}
func writeDefaultReminders(file *os.File) error {
	if file == nil {return errors.New("file handle is nil")}

	rem, err := reminders.CreateEntryByDate(time.UnixDate, "Thu Jul 4 00:00:00 EDT 2030", "Happy Independence Day! Don't blow your hand off")
	if err != nil {return err}

	buf, err := misc.SerializeToBuffer([]reminders.Entry{rem}, func(_ reminders.Entry)bool {return true})
	if err != nil {return err}

	_, err = file.Write(buf)
	return err
}

func writeDefaultCommands(file *os.File) error {
	if file == nil {return errors.New("file handle is nil")}

	cmd := Command{Behavior: RELOAD_REMINDERS}
	buf, err := misc.SerializeToBuffer([]Command{cmd}, func(_ Command)bool{return true})
	if err != nil {return err}

	_, err = file.Write(buf)
	return err
}

func writeDefaultConfig(file *os.File) error {
	if file == nil {return errors.New("file handle is nil")}
	_, err := fmt.Fprintf(file, "[trem]\nTODO = Write actual config here")
	return err
}

func GetConfigFromFile(path string) (Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {return Config{}, err}
	return ParseConfigBuffer(contents)
}

func ParseConfigBuffer(buf []byte) (Config, error) {
	if buf == nil {return Config{}, errors.New("buf is nil")}
	var res Config

	// Convert to strings for ease of use
	contents := string(buf[:])
	lines := strings.Split(contents, "\n")

	const CUTSET = " \n"
	for idx, line := range lines {
		trimmed := strings.Trim(line, CUTSET)
		split := strings.Split(trimmed, "=")
		if l := len(split); l < 2 {
			if l > 0 {fmt.Fprintf(os.Stderr, "Warning - Line %v contains a keyword (%v), but no value", idx, split)}
			continue
		}

		switch {
		case idx == 0 && split[0] != "[trem]": return Config{}, errors.New("Config does not start with \"[trem]\" header")
		case split[0] == "config":		res.ConfigPath 		= strings.Trim(split[1], CUTSET)
		case split[0] == "commands":	res.CommandPath 	= strings.Trim(split[1], CUTSET)
		case split[0] == "reminders":	res.ReminderPath 	= strings.Trim(split[1], CUTSET)

		default: fmt.Fprintf(os.Stderr, "Warning - Encountered unknown keyword in config: %v", trimmed)
		}
	}

	// TODO: Add some sanity checks

	return res, nil
}