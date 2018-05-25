#testergo

name=testergo
wd=vendor/testergo
assets=${wd}/assets

build: bindata
	go build -o ${name} ./${wd}


install: bindata
	go install ./${wd}

run: build
	./${name} ./testdata




bindata:
	go-bindata -debug -ignore='\.scss' -ignore='\.go' -pkg=assets -o=${assets}/assets.go -prefix=${wd} -nocompress -nomemcopy ${assets}/fonts/ ${assets}