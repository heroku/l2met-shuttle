# `l2met-shuttle` is [Log Shuttle][log-shuttle] for [L2met][l2met]

`l2met-shuttle` will extract L2met compatible metrics from its input and rely on Log Shuttle to deliver them to the given URL.

[l2met]: https://github.com/ryandotsmith/l2met
[log-shuttle]: https://github.com/heroku/log-shuttle

## Usage

    $ <command> [<argument>...] | l2met-shuttle <url>

Or, to pass input through to stdout as well as transporting it to an l2met service

    $ <commmand> [<argument>...] | l2met-shuttle --tee <url>
