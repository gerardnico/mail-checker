// Package cmd

package cmd

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mileusna/spf"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Configuration struct {
	Resolver string   `yaml:"resolver" validate:"required,ipv4"`
	Mailers  []string `yaml:"mailers" validate:"required,dive,ipv4"`
	Domains  []string `yaml:"domains" validate:"required,dive,hostname"`
}

var cfgFile string
var config Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mail-checker",
	Short: "Check a mail installation",
	Long:  `Check Spf, Dkim`,
	Run: func(cmd *cobra.Command, args []string) {

		// Set DNS server which will be used by resolver.
		// Default is Google's 8.8.8.8:53
		spf.DNSServer = config.Resolver + ":53"

		_, _ = fmt.Fprintln(os.Stderr, "Spf Check:")
		for _, mailer := range config.Mailers {
			ip := net.ParseIP(mailer)
			for _, domain := range config.Domains {
				r := spf.CheckHost(ip, domain, "foo@"+domain, "")

				format := "  * For domain %s, mailer %s, result is, %s!\n"
				if r != "PASS" {
					log.Fatalf(format, domain, mailer, r)
				}

				_, _ = fmt.Fprintf(os.Stderr, format, domain, mailer, r)

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mail-checker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
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
