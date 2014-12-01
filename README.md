Bitrot Killer
=
A backup utility that checks checksum to prevent bitrot.

Generate checksum file
-
`bitrot_killer --create /home/jerry/shared /home/jerry/shared.json`

Check update
-
`bitrot_killer --check /home/jerry/shared /home/jerry/shared.json [smb://192.168.8.140/public/sda1/Jerry/shared]`

Check duplicated files and empty folder
-
`bitrot_killer --dedup /home/jerry/shared.json`
