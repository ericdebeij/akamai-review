# Release notes

## 0.1.10+11+12 - remove custom resolver
- Miekeg/dns removed, built-in resolver used to lookup CNAME and IP-addresses, skip netstorage, cache results, log errors at debug level

## 0.1.8+9 - fix empty export
- Irregular empty export fix

## 0.1.7 - add cps-overview
- Add contract column to pm-hosts
- Rename cps-certificates to cps-sans
- Add cps-overview to show all cps certificates, their cipher usage

## 0.1.6 - log file improvement
- Multiple log files supported, with logs cumulated to a provided level
- Syntax change in .akamai-review.yaml
```
log:
- level: info
- file: filename.log
```
- Small bug fixes

## 0.1.5 - usage-repgroup
- New subcommand added to collect usage data based on reporting groups
- Multiple small fixes

## 0.1.4 - bug fixes
- Bug fixes only

## 0.1.3 - aliases and debug info
- Adding aliases for contracts and products

## 0.1.2 - pm-behavior
- New subcommand pm-behavior used to grep the behaviors used in a property or the usage of a specific behavior

## 0.1.1 - change command names, hosts-certificate
- Standardize the naming convention for command names
- New subcommand added to collect hosts (from security api) and check the related certificates

## 0.1.0 - restructure command and package structure
- Major overhaul of the command and package structure