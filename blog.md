# Akamai CLI plugin – investigate CDN configurations and traffic

## The challenge

As a Technical Advisor for Akamai, I get quite often questions from customers about their configurations on the Akamai platform. Some are simple to answer, but quite often you need to dig a little bit deeper in the configurations to get the answers you need.

Not a problem when customer has 5 configurations, but most of them don’t and they have 500 or 5000 configurations. Time for some tooling!

![aktool|110x100](media/aktool.jpg)
## Common questions

Questions answered with this tooling:

* Can you provide me with a list of all origins used in my account and check if they do resolve in an IP-address?

* Can you provide me an overview of the billed traffic for last month per CPCode and calculate the diversions from the month before?

* Can you provide me with a list of certificates & SANs and verify if the names are served from akamai CDN (or can they be removed from the certificate)?

* Can you provide me with a list of all active hostnames on Akamai and identify if they are served via the Akamai CDN and identify if the certificate is about to expire?

* Can you provide me a list with all hosts in my account and identify if there is a redirect implemented for HTTP=>HTTPS

* Can you provide me with a list of all properties that use the Edge Redirector Cloudlet?

## Akamai Command Line

Besides the Akamai Control Center (the User Interface), which can be used to view and change your Akamai setup, you can also use the [Akamai Command Line Interface (Akamai-CLI)](https://techdocs.akamai.com/developer/docs/about-clis) to perform tasks. Under the hood the CLI uses the API’s to get the tasks done.

Sometimes you can use the standard subcommands to get the data and massage the output (e.g. with jq) to get the answers you need. But if there is no good existing way to do that, you can also use the API’s yourself and create your own plugin.

Creating a plugin for the Akamai-CLI is not very complex (if you create the plugin with Node, Python or Golang). An akamai-CLI plugin is just a piece of software which can also run on its own. To use it as a plugin for the Akamai CLI, you need to make it available on GitHub and add some additional info.

## The akamai-review plugin

For these common questions mentioned above I have createed such a plugin. Multiple reports in the form of a CSV can be created. It is still up to you to use the data correctly, review the correctness if in doubt and represent it in a nice format.

You can find the repository / binaries on GitHub, download the correct binaries, or you can install the plugin from the Akamai CLI:

```

% akamai install ericdebeij/akamai-review

% akamai review

```

PS: As you can see this is a private repository, so there is no Akamai support for this.

# Repository / more info

More info and the repository:

https://github.com/ericdebeij/akamai-review