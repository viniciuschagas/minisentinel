green = \033[0;92m
blue = \033[0;94m
clear = \033[00m

done = echo "${green}Done${clear}"; echo

.SILENT:

install-deps:
	echo "${blue}--- Installing necessary Go packages ---${clear}"
	go install gotest.tools/gotestsum@latest

tidy:
	echo "${blue}Running go mod tidy...${clear}"
	go mod tidy
	${done}

setup:
	echo "${blue}Setting up the project...${clear}"
	${MAKE} install-deps tidy
	${done}

test:
	echo "${blue}Running tests...${clear}"
	gotestsum --format testname -- -count=1 -v ./...
	${done}