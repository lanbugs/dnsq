package main

/*
dnsq - Tool to check DNS zones syncronity
Written by Maximilian Thoma 2022
*/

import (
	"dnsq/dnsq"
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	NAME        = "dnsq"
	DESCRIPTION = "Tool to check DNS zones syncronity"
	VERSION     = "1.0.3"
)

func Help() {
	fmt.Printf("\n%v - %v - Version %v\n", NAME, DESCRIPTION, VERSION)
	fmt.Println("Written by Maximilian Thoma 2022")
	fmt.Printf("Check https://github.com/lanbugs/dnsq for updates!\n\n")
	fmt.Println("You can use static nameservers to check in config.yml file.")
	fmt.Println("YAML Format:")
	fmt.Println("Nameservers:")
	fmt.Println("  - 1.1.1.1")
	fmt.Printf("  - 1.0.0.1\n\n")
	fmt.Println("You can also define static initial nameserver:")
	fmt.Printf("InitialNameserver: 1.1.1.1\n\n")
	fmt.Println("For the interactive mode and auto mode you can also define a list of domains:")
	fmt.Println("Domains:")
	fmt.Println("  - example.org")
	fmt.Printf("  - example.net\n\n")
	fmt.Println("If you run dnsq with -i the program search after *.yml files in the program folder.")
	fmt.Printf("If you start dnsq only with -c and you have defined also domains in the config all domains will be checked.\n\n")

	flag.Usage()
	syscall.Exit(0)
}

func main() {
	var interactive *bool = flag.BoolP("interactive", "i", false, "Use interactive shell")
	var config *string = flag.StringP("config", "c", "config.yml", "Static nameservers to query for domains, if not defined NS records will be used as nameservers.")
	var domain *string = flag.StringP("domain", "d", "", "Domain to query")
	var ns *string = flag.StringP("ns", "n", "8.8.8.8", "Nameserver for inital query")
	var help *bool = flag.BoolP("help", "h", false, "Show help")

	flag.Parse()

	if *help {
		Help()
	}

	if *interactive {
		configs, _ := filepath.Glob("./*.yml")

		var conf []string

		conf = append(conf, configs...)

		var qs = []*survey.Question{
			{
				Name: "config",
				Prompt: &survey.Select{
					Message: "Choose a config",
					Options: conf,
				},
			},
		}

		answer := struct {
			Config string `survey:"config"`
		}{}

		err := survey.Ask(qs, &answer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		*config = answer.Config

	}

	// Autoload inital config file
	viper.AddConfigPath(".")
	viper.SetConfigFile(*config)
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	var cnf dnsq.Config

	err := viper.Unmarshal(&cnf)

	if err != nil {
		fmt.Printf("unable to decode config, %v", err)
	}

	if *interactive {
		if len(cnf.Domains) > 0 {
			var qsd = []*survey.Question{
				{
					Name: "domain",
					Prompt: &survey.Select{
						Message: "Choose a domain",
						Options: cnf.Domains,
					},
				},
			}

			answerd := struct {
				Domain string `survey:"domain"`
			}{}

			err := survey.Ask(qsd, &answerd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			*domain = answerd.Domain
		}
	}

	if *domain == "" && len(cnf.Domains) == 0 {
		Help()
	}

	var nameservers []string

	// Determine initial nameserver
	if *ns == "8.8.8.8" && len(cnf.InitialNameserver) > 0 {
		*ns = cnf.InitialNameserver
	}

	// Determine nameservers
	if len(cnf.Nameservers) == 0 {
		// Use dynamic determined nameservers from zone file
		nameservers = dnsq.GetNS(*domain, *ns)

	} else {
		// Use static nameservers from config file
		nameservers = cnf.Nameservers
	}

	// If Domains are defined in config file work on each domain
	if *domain == "" && len(cnf.Domains) > 0 {
		for _, d := range cnf.Domains {
			fmt.Printf("\n\n\nDOMAIN: %v\n", d)
			fmt.Println("==========================================================================")

			dnsq.Runner(d, *ns, nameservers)
		}

	} else {
		// Run only single domain
		dnsq.Runner(*domain, *ns, nameservers)
	}

}
