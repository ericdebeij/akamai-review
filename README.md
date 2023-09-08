# akamai-review
Several utilities that can be used to review your Akamai account.

This is sample software, will likely remain in draft state forever and there is no plan to maintain this, use at your own risk.

## Akamai CLI Install
```bash
%  akamai install ericdebeij/akamai-review
```

or just download the binary for your operating system from the release folder and use directly.

## Requirements
* access credentials

## Warranty
This is sample software. As such this software comes with absolutely no warranty.

## Command line
```bash
% akamai review subcommand
```
or use the downloaded binary

Example:
```bash
% akamai review cps-certificates
```

## Caching and rate limits
The first time you run some of the reports, it might take quite a while as properties and rules need to be downloaded. For this reason that kind of (immutable) information is cached. You might also run into rate limits. On subsequent call only changed configurations need to be downloaded. 

## Config file
A config file can be used for global parameters which are often used as well as for some default values of command specific parameters (like config and product).
The default config name is .akamai-review.yaml, the file will be searched for in the current directory as well as in the users home directory.

For an example of the config file, see the example directory.

# Usage
```
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
  usage-repgroup    An overview of the usage per reportinggroup
```

## Options
See [akamai-review](doc/akamai-review.md)

# Reports
The following reports are currently supported

## cps-certificates
List certificates as defined in cps.

Columns: cn,san,cdn

See [akamai-review cps-certificates](doc/akamai-review_cps-certificates.md)

## hosts-certificate
List of all hostnames in your account per property with dns and certificate information

Columns: host,cdn,security,subject-cn,issuer-cn,expires,expire-days

See [akamai-review hosts-certificate](doc/akamai-review_hosts-certificate.md)

## usage-cpcode
An overview of the usage for a month per cpcode and a comparison with the previous month

Columns: cpcode,cpname,repgrp,2023-07(GB),2023-06(GB),diff(GB),2023-07(Hits),2023-06(Hits),diff(Hits)

See [akamai-review usage-cpcode](doc/akamai-review_usage-cpcode.md)

## usage-repgroup
An overview of the usage per reporting groups, multiple months

Columns: month,reportinggroup...

See [akamai-review usage-repgroup](doc/akamai-review_usage-repgroup.md)

## pm-hosts
List of all hostnames in your account per property with dns and certificate information

Columns: group,property,host,edgehost,cdn,ips,cert-subject,cert-issuer,cert-expire

See [akamai-review pm-hosts](doc/akamai-review_pm-hosts.md)

## pm-origins
An overview of the origins

Columns: group,property,origin,origintype,forward,hostmatch,pathmatch,siteshield,ips

See [akamai-review pm-origins](doc/akamai-review_pm-origins.md)

## pm-behaviors
An overview of the behaviors in a propery

Columns: group,property,behaviors

See [akamai-review pm-behaviors](doc/akamai-review_pm-behaviors.md)

# Release notes
See [Release notes](release.md)

# Known issues

Rate controls might hinder the process. 

* Some services handle rate controls correctly (if we run into a 429 whole checking if an IP-address is an Akamai edge-server we just wait the required number of seconds and try again), but others result in an error. 
* There is caching for some immutable elements (e.g. property rules of activated versions), in that case just restarting the process after a while will solve the problem.
* You might want to run the command with "--loglevel debug --logfile debug.log" to capture all requests and check for the last errors. In that case all output including INFO and FATAL will go only to the logfile.
* The hosts-certificate command can only be used once a minute as the rate control for the API to get all hosts of an account can only be run once a minute.

Not all errors are not handled in a consistent way.

See also the [ToDo list](todo.md)

# Contribution

By submitting a contribution (the “Contribution”) to this project, and for good and valuable consideration, the receipt and sufficiency of which are hereby acknowledged, you (the “Assignor”) irrevocably convey, transfer, and assign the Contribution to the owner of the repository (the “Assignee”), and the Assignee hereby accepts, all of your right, title, and interest in and to the Contribution along with all associated copyrights, copyright registrations, and/or applications for registration and all issuances, extensions and renewals thereof (collectively, the “Assigned Copyrights”). You also assign all of your rights of any kind whatsoever accruing under the Assigned Copyrights provided by applicable law of any jurisdiction, by international treaties and conventions and otherwise throughout the world. 

# Notice

Copyright 2021-2023 – Akamai Technologies, Inc.
 
All works contained in this repository, excepting those explicitly otherwise labeled, are the property of Akamai Technologies, Inc.