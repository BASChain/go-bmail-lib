BINDIR=bin

#.PHONY: pbs

all: m a i test
#
#pbs:
#	cd pbs/ && $(MAKE)
#
test:
	go build  -ldflags '-w -s' -o $(BINDIR)/ctest mac/*.go
m:
	CGO_CFLAGS=-mmacosx-version-min=10.11 \
	CGO_LDFLAGS=-mmacosx-version-min=10.11 \
	GOARCH=amd64 GOOS=darwin go build  --buildmode=c-archive -o $(BINDIR)/bmail.a mac/*.go
	cp mac/callback.h $(BINDIR)/
a:
	gomobile bind -v -o $(BINDIR)/BmailLib.aar -target=android github.com/BASChain/go-bmail-lib/android
i:
	gomobile bind -v -o $(BINDIR)/BmailLib.framework -target=ios github.com/BASChain/go-bmail-lib/ios

sol:
	cd resolver/ && $(MAKE)
clean:
	gomobile clean
	rm $(BINDIR)/*