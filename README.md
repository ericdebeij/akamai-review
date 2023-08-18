# akamai-review
Several utilities that can be used to review your Akamai account.

This is sample software, will likely remain in draft state forever and there is no plan to maintain this, use at your own risk.

## Akamai CLI Install
```bash
%  akamai install https://github.com/ericdebeij/akamai-review
```

or just download the binary for your operating system from the release folder and use directly.

## Requirements
* access credentials

## Warranty
This is sample software. As such this software comes with absolutely no warranty.

## Command line
```bash
% akamai review
```
or use the downloaded binary

## Usage
```bash
akamai-review is a utility collection to extract information from
your akamai account and perform checks on it that need to be performed
on a regular base.

Usage:
  akamai-review [command]

Available Commands:
  cps-certificates  List certificates as defined in cps
  help              Help about any command
  hosts-certificate List of all hostnames in your account per property with dns and certificate information
  pm-behaviors      An overview of the behaviors in a propery
  pm-hosts          List of all hostnames in your account per property with dns and certificate information
  pm-origins        An overview of the origins
  usage-cpcode      An overview of the usage for a month per cpcode and a comparison with the previous month

Flags:
      --accountkey string   akamai account switch key
      --cache string        cache folder (default "~/.akamai-cache")
      --config string       config file with all default parameters (default ".akamai-review.yaml")
      --edgerc string       akamai location of the credentials file (default "~/.edgerc")
  -h, --help                help for akamai-review
      --logfile string      logging output
      --loglevel string     logging level (default "FATAL")
      --resolver string     resolver to be used (default "8.8.8.8:53")
      --section string      akamai section of the credentials file (default "default")
      --warningdays int     warning days for certificate issues (default 14)

Use "akamai-review [command] --help" for more information about a command.
```
## Config file
A config file can be used for global parameters which are often used as well as for some default values of command specific parameters (like config and product).
The default config name is .akamai-review.yaml, the file will be searched for in the current directory as well as in the users home directory.

For an example of the config file, see the example directory.

# Contribution

By submitting a contribution (the “Contribution”) to this project, and for good and valuable consideration, the receipt and sufficiency of which are hereby acknowledged, you (the “Assignor”) irrevocably convey, transfer, and assign the Contribution to the owner of the repository (the “Assignee”), and the Assignee hereby accepts, all of your right, title, and interest in and to the Contribution along with all associated copyrights, copyright registrations, and/or applications for registration and all issuances, extensions and renewals thereof (collectively, the “Assigned Copyrights”). You also assign all of your rights of any kind whatsoever accruing under the Assigned Copyrights provided by applicable law of any jurisdiction, by international treaties and conventions and otherwise throughout the world. 

# Notice

Copyright 2021-2023 – Akamai Technologies, Inc.
 
All works contained in this repository, excepting those explicitly otherwise labeled, are the property of Akamai Technologies, Inc.