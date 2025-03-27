// Package cmd

package cmd

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kuberhealthy/kuberhealthy/v2/pkg/checks/external/checkclient"
	"github.com/mileusna/spf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
)

type Configuration struct {
	Resolver string   `yaml:"resolver" validate:"required,ipv4"`
	Mailers  []string `yaml:"mailers" validate:"required,dive,ipv4"`
	Domains  []string `yaml:"domains" validate:"required,dive,hostname"`
}

// Version The version
var Version = "set-with-x-flags" // set at compile time with x flags

var cfgFile string
var config Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "mail-checker",
	Short:   "Check a mail installation",
	Long:    `Check Spf, Dkim`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {

		// Set DNS server which will be used by resolver.
		// Default is Google's 8.8.8.8:53
		spf.DNSServer = config.Resolver + ":53"

		// A kuberhealthy will fail only if it can't report it
		// Docs: https://github.com/kuberhealthy/kuberhealthy/blob/master/docs/CHECK_CREATION.md
		isKhRun := false
		if os.Getenv("KH_REPORTING_URL") != "" {
			isKhRun = true
		}

		// Spf Check
		_, _ = fmt.Fprintln(os.Stderr, "Spf Check:")
		var errors []string
		for _, mailer := range config.Mailers {
			ip := net.ParseIP(mailer)
			for _, domain := range config.Domains {
				r := spf.CheckHost(ip, domain, "foo@"+domain, "")

				format := "  * For domain %s, mailer %s, result is, %s!\n"
				result := fmt.Sprintf(format, domain, mailer, r)
				if r != "PASS" {
					errors = append(errors, result)
					continue
				}
				_, _ = fmt.Println(result)

			}
		}

		// Error
		if len(errors) > 0 {
			if isKhRun {
				err := checkclient.ReportFailure(errors)
				if err != nil {
					log.Fatal("Unable to report the following errors to kuberhealthy", errors)
				}
				os.Exit(0)
			}
			log.Fatal(errors)
		}

		// Success
		if isKhRun {
			err := checkclient.ReportSuccess()
			if err != nil {
				log.Fatal("Could not report success")
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Version
	rootCmd.SetVersionTemplate("Version: {{.Version}}")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile == "" {
		_ = rootCmd.Help()
		log.Fatal("The config is mandatory")
	}

	// Use config file from the flag.
	viper.SetConfigFile(cfgFile)

	// Read in environment variables that match
	viper.AutomaticEnv()

	// Read it in.
	_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) || errors.Is(err, os.ErrNotExist) {
			log.Fatal("Config file not found")
		}
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to unmarshall the config %v", err)
	}
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		log.Fatalf("Error validating the config file: %v\n", err)
	}

}
