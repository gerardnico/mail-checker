// Package cmd

package cmd

import (
	"errors"
	"fmt"
	"github.com/gerardnico/mail-checker/pkg/report"
	"github.com/go-playground/validator/v10"
	"github.com/mileusna/spf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"strings"
)

type Configuration struct {
	Resolver    string             `yaml:"resolver" validate:"required,ipv4"`
	Mailers     []string           `yaml:"mailers" validate:"required,dive,ipv4"`
	Domains     []string           `yaml:"domains" validate:"required,dive,hostname"`
	PushGateway report.PushGateway `yaml:"pushgateway"`
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

		metaCheck := report.MetaCheck{
			Job: "mail-checker",
		}

		// Set DNS server which will be used by resolver.
		// Default is Google's 8.8.8.8:53
		spf.DNSServer = config.Resolver + ":53"

		// Spf Check
		_, _ = fmt.Fprintln(os.Stderr, "Spf Check:")
		for _, mailer := range config.Mailers {
			ip := net.ParseIP(mailer)
			for _, domain := range config.Domains {

				// Metrics
				spfCheckMetrics := report.MetricDefinition{
					Name: "check",
					Type: report.Gauge,
					Labels: map[string]string{
						"domain": domain,
						"mailer": mailer,
						"type":   "spf",
					},
				}
				metaCheck.Metrics = append(metaCheck.Metrics, spfCheckMetrics)

				r := spf.CheckHost(ip, domain, "foo@"+domain, "")

				format := "  * For domain %s, mailer %s, result is, %s!\n"
				result := fmt.Sprintf(format, domain, mailer, r)
				if r != "PASS" {
					metaCheck.Errors = append(metaCheck.Errors, result)
					spfCheckMetrics.Value = 1
					continue
				}
				spfCheckMetrics.Value = 0
				_, _ = fmt.Println(result)

			}
		}

		// Pointer to check if the error was reported to a sinc
		// if error and reported
		//   exit status zero
		//  else
		//   exit status 1
		isReported := false

		// KuberHealthy
		if os.Getenv("KH_REPORTING_URL") != "" {
			report.ToKuberHealthy(metaCheck)
			isReported = true
		}

		// PushGateway
		pushGatewayUrl := getConfigValue("pushgateway.url", "").(string)
		if pushGatewayUrl != "" {
			config.PushGateway.Url = pushGatewayUrl
			err := report.ToPushgateway(config.PushGateway, metaCheck)
			if err != nil {
				log.Fatal("Unable to report to pushgateway")
			}
			isReported = true
		}

		// Cli run
		if isReported {
			return
		}
		if len(metaCheck.Errors) > 0 {
			log.Fatal(metaCheck.Errors)
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

	// Env
	// Read in environment variables that match
	// Environment variable prefix
	viper.SetEnvPrefix("MAIL_CHECKER")
	// Automatically use environment variables that match
	viper.AutomaticEnv()
	// Replace dots with underscores in env var names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Config file
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

// Example function to demonstrate config retrieval
func getConfigValue(key string, defaultValue interface{}) interface{} {
	// Priority order:
	// 1. Command line flag (handled by Cobra/Viper automatically)
	// 2. Environment variable
	// 3. Config file
	// 4. Default value

	// Check if the value is set via flag or env
	if viper.IsSet(key) {
		return viper.Get(key)
	}

	// Return default if no value found
	return defaultValue
}
