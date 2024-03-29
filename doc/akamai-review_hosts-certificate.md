## akamai-review hosts-certificate

List of all hostnames in your account per property with dns and certificate information

### Synopsis

An overview of the properties and the hostnames associated within the property. In order to find this information the property manager hostnames are downloaded (and stored in a cache).
The related edgehost is shown and the host is checked to see if it is actually served by Akamai, resolves in a proper IP-address and information regarding the certificate being used

```
akamai-review hosts-certificate [flags]
```

### Options

```
      --export string   name of the exportfile (default "hosts-certificate.csv")
  -h, --help            help for hosts-certificate
      --httptest        run an http test to check if http->https redirect is implemented
      --match string    regular expression for hostmatch
      --skip string     regular expression for hostskip (default "^failover\\..*$")
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
