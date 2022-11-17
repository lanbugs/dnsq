# dnsq - Tool to check DNS zones syncronity

This is my first GO project to learn the programming language :-)

[![asciicast](https://asciinema.org/a/538472.svg)](https://asciinema.org/a/538472)

## What can dnsq do for you?

dnsq has several options to check with the SOA serial if the replication of zones are completed.

You can predefine a config.yml or several config.yml files in the program folder to define the nameservers to use and which domains should be checked.

## The config format

```
InitialNameserver: 8.8.8.8 
Nameservers:
  - 1.1.1.1
  - 1.0.0.1
Domains:
  - example.org
  - example.net
```

### InitialNameserver
This nameserver is used to determine the NS records from the zone if you have not defined "Nameservers" in the config.

### Nameservers
You can static define the nameservers which should be checked.

### Domains
You can define a list of domains which should be all checked or be present  for the interactive mode

## Usage

```
Usage of dnsq:
  -c, --config string   Static nameservers to query for domains, if not defined NS records will be used as nameservers. (default "config.yml")
  -d, --domain string   Domain to query
  -h, --help            Show help
  -i, --interactive     Use interactive shell
  -n, --ns string       Nameserver for inital query (default "8.8.8.8")
```

Per default dnsq searches for config.yml. Here you can define all options or only the "InitialNameserver".

If you run the interactive mode you got an menu to select a config and in the second step an domain.

```
dnsq -i
? Choose a config  [Use arrows to move, type to filter]
> config.yml
  config_example.yml

dnsq -i
? Choose a config config.yml
? Choose a domain  [Use arrows to move, type to filter]
> example.org
  example.net

dnsq -i
? Choose a config config2.yml
? Choose a domain example.org
Run executed at: 2022-11-17 18:26:09.23681 +0100 CET m=+22.341775301

Domain       Nameserver  Serial      Difference  
example.org  1.1.1.1     2022091149  0
example.org  1.0.0.1     2022091149  0
```

