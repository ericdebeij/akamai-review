## akamai-review usage-cpcode

An overview of the usage for a month per cpcode and a comparison with the previous month

### Synopsis

Uses the billing API to get an overview of the usage for a specific month and compares this with the previous month, both bytes and hits

```
akamai-review usage-cpcode [flags]
```

### Options

```
      --contract string   contract to be used
      --export string     name of the exportfile (default "usage-cpcode_PERIOD.csv")
  -h, --help              help for usage-cpcode
      --period string     period to be investigated
      --product string    product code to be used
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
