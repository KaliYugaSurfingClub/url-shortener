sqlite3 storage.db < ./storage/init.sql
sqlite3 cmd/storageTests/test.db < ./storage/init.sql
go test cmd/storageTests/initTestDB_test.go
go test cmd/storageTests/selectAll_test.go
