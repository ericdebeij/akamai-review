# Release notes

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