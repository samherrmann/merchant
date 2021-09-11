
gobuild = mkdir -p dist/$$GOOS-$$GOARCH && go build -o dist/$$GOOS-$$GOARCH .
tar = (cd dist && tar -czvf $$GOOS-$$GOARCH.tar.gz $$GOOS-$$GOARCH/*)
zip = (cd dist && zip -r $$GOOS-$$GOARCH.zip $$GOOS-$$GOARCH/*)

build:
	mkdir -p dist && go build -o dist .

build.all: 
	export GOOS=linux; export GOARCH=386; $(gobuild) && $(tar)
	export GOOS=linux; export GOARCH=amd64; $(gobuild) && $(tar)
	export GOOS=windows; export GOARCH=amd64; $(gobuild) && $(zip)
	export GOOS=windows; export GOARCH=386; $(gobuild) && $(zip)

test:
	go test ./... -cover

clean:
	rm -rf dist

# Resources:
# List of available target OSs and architectures: 
# https://golang.org/doc/install/source#environment