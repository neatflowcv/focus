#!/bin/bash
set -euo pipefail
# ARG_POSITIONAL_SINGLE([id],[id])
# ARG_POSITIONAL_SINGLE([parent],[parent])
# ARG_POSITIONAL_SINGLE([next],[next])

die()
{
	local _ret="${2:-1}"
	test "${_PRINT_HELP:-no}" = yes && print_help >&2
	echo "$1" >&2
	exit "${_ret}"
}

_positionals=()

print_help()
{
	printf 'Usage: %s <id> <parent> <next>\n' "$0"
	printf '\t%s\n' "<id>: id"
	printf '\t%s\n' "<parent>: parent"
	printf '\t%s\n' "<next>: next"
}


parse_commandline()
{
	_positionals_count=0
	local _key
	while test $# -gt 0
	do
		_last_positional="$1"
		_positionals+=("$_last_positional")
		_positionals_count=$((_positionals_count + 1))
		shift
	done
}


handle_passed_args_count()
{
	local _required_args_string="'id', 'parent' and 'next'"
	test "${_positionals_count}" -ge 3 || _PRINT_HELP=yes die "FATAL ERROR: Not enough positional arguments - we require exactly 3 (namely: $_required_args_string), but got only ${_positionals_count}." 1
	test "${_positionals_count}" -le 3 || _PRINT_HELP=yes die "FATAL ERROR: There were spurious positional arguments --- we expect exactly 3 (namely: $_required_args_string), but got ${_positionals_count} (the last one was: '${_last_positional}')." 1
}


assign_positional_args()
{
	local _positional_name _shift_for=$1
	_positional_names="_arg_id _arg_parent _arg_next "

	shift "$_shift_for"
	for _positional_name in ${_positional_names}
	do
		test $# -gt 0 || break
		eval "$_positional_name=\${1}" || die "Error during argument parsing, possibly an Argbash bug." 1
		shift
	done
}

parse_commandline "$@"
handle_passed_args_count
assign_positional_args 1 "${_positionals[@]}"

curl -H 'Content-Type: application/json' \
-X PATCH \
http://localhost:8080/focus/tasks/${_arg_id} \
-H "Authorization: Bearer $TOKEN" \
-d "{
    \"parent_id\": \"${_arg_parent}\",
    \"next_id\": \"${_arg_next}\"
}"