Bitrot Killer
==
[![PkgGoDev](https://pkg.go.dev/badge/github.com/codingsince1985/bitrot_killer)](https://pkg.go.dev/github.com/codingsince1985/bitrot_killer)

A backup utility that checks checksum to prevent bitrot.

Generate checksum file
--
`$ bitrot_killer --create /home/jerry/shared /home/jerry/shared.json`

Check & update checksum file
--
`$ bitrot_killer --check /home/jerry/shared /home/jerry/shared.json`

Find duplicated files and empty folder
--
`$ bitrot_killer --dedup /home/jerry/shared.json`

License
==
Bitrot Killer is distributed under the terms of the MIT license. See LICENSE for details.
