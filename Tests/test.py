#!/usr/bin/python

import os
from os import path
import subprocess
from subprocess import Popen, PIPE

import yaml
import re

DEVNULL = open(os.devnull, 'wb')


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


failed = []
include = []
exclude = []

try:
    with open("./include", 'r') as f:
        for line in f.readlines():
            if line.startswith('#'):
                continue
        include.append(line.strip())
except:
    pass

try:
    with open("./exclude", 'r') as f:
        for line in f.readlines():
            if line.startswith('#'):
                continue
    exclude.append(line.strip())
except:
    pass


def check_test(dirpath, fname, exit_code, output):
    fullname = os.path.join(dirpath, fname)
    with open(fullname, 'r') as stream:
        struct = yaml.load(stream)

    output_verbose = output

    if struct['exit_code'] != exit_code:
        print "\nexit_code: ", exit_code,
        print output_verbose
        return False

    if 'output' in  struct:
        exactly = True
        for output_field in struct['output']:
            if 'exactly' in output_field:
                exactly = output_field['exactly']
            data = ""
            if 'string' in output_field:
                data = output_field['string']
            elif 'file' in output_field:
                fullname = os.path.join(dirpath, output_field['file'])
                with open(fullname, 'r') as myfile:
                    data = myfile.read()

            data = data.lower().replace(" ", "").replace('\n', '')

            output = re.sub(r'^\+.*\n?', '', output, flags=re.MULTILINE)

            output = output.lower().replace(" ", "").replace('\n', '')

            if exactly:
                ok = data == output
            else:
                ok = output in data or data in output

            if not ok:
                print "\nexpected: " + data
                print "got: " + output
                print output_verbose
                return False

    return True

def cleanup(dirpath, files):
    if "cleanup.sh" in files:
        print "cleanup ",dirpath
        fname = os.path.join(os.path.abspath(dirpath), "cleanup.sh")
        subprocess.call([fname], stdout=DEVNULL, cwd=dirpath)


def init(dirpath, files):
    if "init.sh" in files:
        print "init ",dirpath
        fname = os.path.join(os.path.abspath(dirpath), "init.sh")
        subprocess.call([fname], stdout=DEVNULL, cwd=dirpath)


#modified os.walk
def walk(top,init,cleanup):
    islink, join, isdir = os.path.islink, os.path.join, os.path.isdir
    try:
        names = os.listdir(top)
    except os.error:
        return

    dirs, nondirs = [], []
    for name in names:
        if isdir(join(top, name)):
            dirs.append(name)
        else:
            nondirs.append(name)

    init(top,nondirs)

    yield top, dirs, nondirs

    for name in dirs:
        new_path = join(top, name)
        if not islink(new_path):
            for x in walk(new_path,init,cleanup):
                yield x

    cleanup(top,nondirs)

for dirpath, dirs, files in walk(".",init,cleanup):
    if dirpath in exclude:
        continue
    if len(include) > 0 and dirpath not in include:
        continue

    path = os.path.abspath(dirpath)

    if "run.sh" in files:
        print "Testing " + dirpath + "...",
        fname = os.path.join(path, "run.sh")
        p = Popen([fname], stdin=PIPE, stdout=PIPE, stderr=subprocess.STDOUT, cwd=path)
        output, err = p.communicate()
        exit_code = p.returncode
        result = check_test(path, "check.yaml", exit_code, output)
        if result:
            print bcolors.OKGREEN + "OK" + bcolors.ENDC
        else:
            print bcolors.FAIL + "FAIL" + bcolors.ENDC
            failed.append(dirpath)

if len(failed) > 0:
    print bcolors.FAIL + "FAILED:" + bcolors.ENDC
    for failure in failed:
        print bcolors.FAIL + failure + bcolors.ENDC
    exit(1)
else:
    print bcolors.OKGREEN + "ALL OK" + bcolors.ENDC
    exit(0)
