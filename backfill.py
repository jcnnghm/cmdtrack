#!/bin/python
# This should be used by doing something like:
# fc -t %s -l 1 10 | python backfill.py
# The 1 and 10 in the example above specify which lines to capture, inclusive
import fileinput
import re
import subprocess
import sys

# https://regex101.com/r/1lguIy/1
regex = r"^\s*?(\d+\*?)\s*?(\d+)\s+(.*)$"

for line in fileinput.input():
    matches = re.finditer(regex, line, re.MULTILINE)
    for match in matches:
        command_num = match.group(1)
        timestamp = match.group(2)
        command = match.group(3)

        args = [
            "cmdtrack",
            "track",
            '--command', command,
            '--timestamp', timestamp,
            '--workdir', '~',
            # '--url', 'http://localhost:8080/',  # Uncomment this line for
            # testing
        ]
        exit_code = subprocess.call(args, stdout=sys.stdout, stderr=sys.stderr)
        if exit_code != 0:
            print "Executing command failed"
            print args
            sys.exit(-1)

        print command_num, timestamp, command
