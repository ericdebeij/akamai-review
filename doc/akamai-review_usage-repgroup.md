## akamai-review usage-repgroup

An overview of the usage per reportinggroup

### Synopsis

Uses the billing API to get an overview of the usage per reportinggroup. Calculation is based on CPCode details and the actual reporting-groups are used (not the historical reportinggroups)

```
akamai-review usage-repgroup [flags]
```

### Options

```
      --contract string      contract to be used
      --export string        name of the exportfile (default "usage-repgroup.csv")
      --from string          from month (format YYYY-MM) (default "2023-12")
  -h, --help                 help for usage-repgroup
      --product string       product code to be used
      --rgroup stringArray   reporting groups (default all)
      --to string            to month (format YYYY-MM) (default "2023-12")
      --type string          Statistic type [Bytes|Hits] (default "Bytes")
      --unit string          Unit (default Bytes:GB, Hits:Hit)
```

### Options inherited from parent commands

```
      --accountkey string   akamai account switch key
      --cache string        cache folder (default "~/.akamai-cache")
      --config string       config file with all default parameters (default ".akamai-review.yaml")
      --edgerc string       akamai location of the credentials file (default "~/.edgerc")
      --logfile string      logging output
      --loglevel string     logging level
      --resolver string     resolver to be used
      --section string      akamai section of the credentials file (default "default")
      --warningdays int     warning days for certificate issues (default 14)
```

### SEE ALSO

* [akamai-review](akamai-review.md)	 - Review your account assets

###### Auto generated by spf13/cobra on 15-Jan-2024
