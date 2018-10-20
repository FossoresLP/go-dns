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
- There are no tests (I will add some when I find time but feel free to open a PR if you want something specific tested)

