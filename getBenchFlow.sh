#!/bin/sh
set -e
#
# This script is meant for quick & easy install via:
#   'curl -sSL https://github.com/benchflow/client/releases/download/v-dev/benchflow | sh'
# or:
#   'wget -qO- https://github.com/benchflow/client/releases/download/v-dev/benchflow | sh'
#
# How to use this script to install BenchFlow:
#   1. Log into your Ubuntu installation as a user with `sudo` privileges.
#	2. Verify that you have `wget` installed.

# 		$ which wget

# 	3. If `wget` isn’t installed, install it after updating your manager:

# 		$ sudo apt-get update

# 		$ sudo apt-get install wget
#
# 	4. Get the latest BenchFlow package.

#		$ wget -qO- https://github.com/benchflow/client/releases/download/v-dev/getBenchFlow.sh | sh

# 	   The system prompts you for your sudo password. Then, it downloads and installs BenchFlow and its dependencies.

# 	5. Verify `benchflow` is installed correctly.

#		$ benchflow
# 		Usage: benchflow.py [OPTIONS] COMMAND [ARGS]...
# 	    Options:
#   		--help  Show this message and exit.
#
#
# Inspired from: https://get.docker.com and https://docs.docker.com/linux/step_one/

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

do_install() {

	if command_exists benchflow; then
		cat >&2 <<-'EOF'
			Warning: the "benchflow" command appears to already exist on this system.

			You may press Ctrl+C now to abort this script.
		EOF
		( set -x; sleep 20 )
	fi

	user="$(id -un 2>/dev/null || true)"

	sh_c='sh -c'
	if [ "$user" != 'root' ]; then
		if command_exists sudo; then
			sh_c='sudo -E sh -c'
		elif command_exists su; then
			sh_c='su -c'
		else
			cat >&2 <<-'EOF'
			Error: this installer needs the ability to run commands as root.
			We are unable to find either "sudo" or "su" available to make this happen.
			EOF
			exit 1
		fi
	fi

	curl=''
	if command_exists curl; then
		curl='curl -sSL'
	elif command_exists wget; then
		curl='wget -qO-'
	fi

	$sh_c 'curl -L https://github.com/benchflow/client/releases/download/v-dev/benchflow > /usr/local/bin/benchflow'
	$sh_c 'chmod +x /usr/local/bin/benchflow'
}

# wrapped up in a function so that we have some protection against only getting
# half the file during "curl | sh"
do_install