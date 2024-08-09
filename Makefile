test:
	go build cmd/visualuze-gode-deps/main.go
	mv ./main testdata
	cd ./testdata && ./main visualize > plantuml/preview.puml
	rm -rf ./testdata/main


debug:
	go build cmd/visualuze-gode-deps/main.go
	mv ./main testdata
	cd ./testdata && ./main visualize
	rm -rf ./testdata/main
