Bitrot Killer
=============

A backup utility that checks checksum to prevent bitrot.

##### Generate checksum file for a directory
go run backup.go --**create** /home/jerry/shared /home/jerry/shared.json

##### Check directory and optionally sync changes to remote directory
go run backup.go --**check** /home/jerry/shared /home/jerry/shared.json *smb://192.168.8.140/public/sda1/Jerry/shared*
