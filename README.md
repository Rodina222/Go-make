# gomake-Rodina

This is a useful app for parsing makefiles and running target commands while taking dependencies into account.

## __Features:__
- Parsing the command line.
- Reading a Makefile.
- Building a graph of vertices that represent targets, with each vertex containing dependencies and commands as string slices.
- Checking for cyclic dependence and exiting when one is detected.
- Command execution order, where commands of dependencies are executed first, followed by commands of the input target.

## __Formatting of makefiles for ease of use:__
- Each target is at the start of a line.
- Each command is on a line that starts with a single tab (/t).
- At least one command must be executed for every target.

## __Manual:__

1. Clone the repository:
```sh
$ git clone https://github.com/codescalersinternships/gomake-Rodina.git
```
2. Go to the repository directory:
```sh
$ cd gomake-Rodina
```
3. Install dependencies:
```sh
$ go get -d ./...
```
4. Build the package:
```sh
$ go build -o "bin/gomake" main.go
```
5. Go to bin:
```sh
$ cd bin
```
 ### __How to use?__
```sh
$ ./gomake -f Makefile -t target
```
### __Here is an example for a makefile:__

```sh
build:
	@echo 'executing build'
	echo 'cmd2'

test: build publish
	@echo 'executing test'

publish: test 
	@echo 'executing publish'

```
### __How to test?:__

Run all the tests as follows: 
```sh
go test ./....
```
If all tests pass on, the result should show that the tests were successful as follows:
```sh
PASS
ok      github.com/codescalersinternships/gomake-Rodina/internal        0.006s
```
If any test fails, the output will tell which test failed.
