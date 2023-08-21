## akamai-review cps-certificates

List certificates as defined in cps

### Synopsis

List of the certificates, the SAN in the certificates. Additional information is provided to check whether the CN or SAN entry is actually served via Akamai

```
akamai-review cps-certificates [flags]
```

### Options

```
      --contract string   contract to be used
      --export string     contract to be used (default "cps-certificates.csv")
  -h, --help              help for cps-certificates
```

### Options inherited from parent commands

```
      --accountkey string   akamai account switch key
      --cache string        cache folder (default "~/.akamai-cache")
      --config string       config file with all default parameters (default ".akamai-review.yaml")
      --edgerc string       akamai location of the credentials file (default "~/.edgerc")
      --logfile string      logging output
      --loglevel string     logging level (default "FATAL")
      --resolver string     resolver to be used (default "8.8.8.8:53")
      --section string      akamai section of the credentials file (default "default")
      --warningdays int     warning days for certificate issues (default 14)
```

### SEE ALSO

* [akamai-review](akamai-review.md)	 - Review your account assets

###### Auto generated by spf13/cobra on 21-Aug-2023