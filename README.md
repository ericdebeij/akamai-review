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

## Usage
```bash
% akamai-review is a utility collection to extract information from
your akamai account and perform checks on it that need to be performed
on a regular base.

Usage:
  akamai-review [command]

Available Commands:
  alb                  Priving an overview of ALB configuration
  certificates         Check your certificates
  completion           Generate the autocompletion script for the specified shell
  help                 Help about any command
  properties           report on properties in the account
  properties-host      report on hosts used in the properties in the account
  properties-origin    report on origins used in the properties in the account
  report               Run reports
  usage-cpcode         Reports based on usage as part of the billing data

Flags:
      --account string   account switch key
      --config string    config file with all default parameters (default ".akamai-review.yaml")
      --edgerc string    location of the credentials file
  -h, --help             help for akamai-review
      --section string   section of the credentials file

Use "akamai-review [command] --help" for more information about a command.
```

# Contribution

By submitting a contribution (the “Contribution”) to this project, and for good and valuable consideration, the receipt and sufficiency of which are hereby acknowledged, you (the “Assignor”) irrevocably convey, transfer, and assign the Contribution to the owner of the repository (the “Assignee”), and the Assignee hereby accepts, all of your right, title, and interest in and to the Contribution along with all associated copyrights, copyright registrations, and/or applications for registration and all issuances, extensions and renewals thereof (collectively, the “Assigned Copyrights”). You also assign all of your rights of any kind whatsoever accruing under the Assigned Copyrights provided by applicable law of any jurisdiction, by international treaties and conventions and otherwise throughout the world. 

# Notice

Copyright 2021 – Akamai Technologies, Inc.
 
All works contained in this repository, excepting those explicitly otherwise labeled, are the property of Akamai Technologies, Inc.