package core

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/huuhoait/gin-vue-admin/server/core/internal"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Viper configuration
func Viper() *viper.Viper {
	config := getConfigPath()

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.GVA_CONFIG); err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}

	// Environment-variable overrides for sensitive values. Config files are
	// often committed to VCS; secrets must be overridable at deploy time.
	applySecretOverrides()

	// root AdaptNature According torootPositionSetGoFindToCorrespondingmigrateMovePositionSet,EnsurerootpathHaveEffect
	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}

// getConfigPath getconfigurationFilepath, Priority: CommandRow > environment variables > defaultValue
func getConfigPath() (config string) {
	// `-c` flag parse
	flag.StringVar(&config, "c", "", "choose config file.")
	flag.Parse()
	if config != "" { // CommandRowParameterNotEmpty WillValueAssignValueAtconfig
		fmt.Printf("you are usingCommandRowof '-c' ParameterPassofValue, config 's path is %s\n", config)
		return
	}
	if env := os.Getenv(internal.ConfigEnv); env != "" { // Judgeenvironment variables GVA_CONFIG
		config = env
		fmt.Printf("you are using %s environment variables, config 's path is %s\n", internal.ConfigEnv, config)
		return
	}

	switch gin.Mode() { // According to gin Modefile name
	case gin.DebugMode:
		config = internal.ConfigDebugFile
	case gin.ReleaseMode:
		config = internal.ConfigReleaseFile
	case gin.TestMode:
		config = internal.ConfigTestFile
	}
	fmt.Printf("you are using gin of %s ModeRun, config 's path is %s\n", gin.Mode(), config)

	_, err := os.Stat(config)
	if err != nil || os.IsNotExist(err) {
		config = internal.ConfigDefaultFile
		fmt.Printf("configurationFilepathDoes not exist, UsedefaultconfigurationFilepath: %s\n", config)
	}

	return
}
