# Hello

This thing is a dumb utility that takes current rates from bitvavo and poops out
your total winnings (or losses), given what you put in the config. If you
purchase new coins you'll have to update that yaml file as well.

# How to use

    $ go install github.com/pubkraal/amirich

And copy the conf_sample.yaml to ~/.amirich.yaml. The contents are
self-explanatory. Then just run

    $ amirich

And magic output appears

You can also clone the repo and do the `go run main.go` thing, but you still
need the yaml file in place.