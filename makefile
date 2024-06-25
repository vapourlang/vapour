default: rinstall
	go build .

rinstall: rcheck
	cd package && R -s -e "devtools::install()"

rcheck: rdocument
	cd package && R -s -e "devtools::check()"

rdocument:
	cd package && R -s -e "devtools::document()"

dev: rinstall
	go run .

lsp:
	go build .
