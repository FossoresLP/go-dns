DNS
===
This repository contains a (WIP) DNS server written in Go(lang) and all the packages that it is composed of.

Feel free to use these packages for your own projects but keep the license in mind.

Known issues:
-------------
- Does not perform recursive resolution itself but instead relies on Cloudflares 1.1.1.1
- Parsing issues for some messages (i.e. Microsoft.com and other Microsoft websites - these are being worked on)
- Does not support normal zone files
- Zones file format is somewhat awkward as of right now (This will be fixed by switching away from TOML)
- Not fully tested (See [#1](https://github.com/fossoreslp/go-dns/issues/1))

