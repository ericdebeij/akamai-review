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

## Caching and rate limits
The first time you run some of the reports, it might take quite a while as properties and rules need to be downloaded. For this reason that kind of (immutable) information is cached. You might also run into rate limits. On subsequent call only changed configurations need to be downloaded. 

## Config file
A config file can be used for global parameters which are often used as well as for some default values of command specific parameters (like config and product).
The default config name is .akamai-review.yaml, the file will be searched for in the current directory as well as in the users home directory.

For an example of the config file, see the example directory.

# Usage
See [akamai-review](akamai-review.md)

# Reports
The following reports are currently supported

## cps-certificates
See [akamai-review cps-certificates](akamai-review_cps-certificates.md)

List certificates as defined in cps.

Columns: cn,san,cdn

## hosts-certificate
See [akamai-review hosts-certificate](akamai-review_hosts-certificate.md)

List of all hostnames in your account per property with dns and certificate information

Columns: host,cdn,security,subject-cn,issuer-cn,expires,expire-days

## usage-cpcode
See [akamai-review usage-cpcode](akamai-review_usage-cpcode.md)

An overview of the usage for a month per cpcode and a comparison with the previous month

Columns: cpcode,cpname,repgrp,2023-07(GB),2023-06(GB),diff(GB),2023-07(Hits),2023-06(Hits),diff(Hits)

## pm-hosts
See [akamai-review pm-hosts](akamai-review_pm-hosts.md)

List of all hostnames in your account per property with dns and certificate information

Columns: group,property,host,edgehost,cdn,ips,cert-subject,cert-issuer,cert-expire

## pm-origins
See [akamai-review pm-origins](akamai-review_pm-origins.md)

An overview of the origins

Columns: group,property,origin,origintype,forward,hostmatch,pathmatch,siteshield,ips

## pm-behaviors
See [akamai-review pm-behaviors](akamai-review_pm-behaviors.md)

An overview of the behaviors in a propery

Columns: group,property,behaviors

# Contribution

By submitting a contribution (the “Contribution”) to this project, and for good and valuable consideration, the receipt and sufficiency of which are hereby acknowledged, you (the “Assignor”) irrevocably convey, transfer, and assign the Contribution to the owner of the repository (the “Assignee”), and the Assignee hereby accepts, all of your right, title, and interest in and to the Contribution along with all associated copyrights, copyright registrations, and/or applications for registration and all issuances, extensions and renewals thereof (collectively, the “Assigned Copyrights”). You also assign all of your rights of any kind whatsoever accruing under the Assigned Copyrights provided by applicable law of any jurisdiction, by international treaties and conventions and otherwise throughout the world. 

# Notice

Copyright 2021-2023 – Akamai Technologies, Inc.
 
All works contained in this repository, excepting those explicitly otherwise labeled, are the property of Akamai Technologies, Inc.