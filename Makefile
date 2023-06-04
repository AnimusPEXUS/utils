GONOPROXY="github.com/AnimusPEXUS/*"

all: get

get:
		go get -u -v "./..."
		go mod tidy
