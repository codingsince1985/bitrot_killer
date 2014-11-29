Bitrot Killer
=
A backup utility that checks checksum to prevent bitrot.

Generate checksum file
-
go run backup.go --**create** /home/jerry/shared /home/jerry/shared.json

Check update
-
go run backup.go --**check** /home/jerry/shared /home/jerry/shared.json *smb://192.168.8.140/public/sda1/Jerry/shared*

Check duplicated
-
go run backup.go --**dedup** /home/jerry/shared.json