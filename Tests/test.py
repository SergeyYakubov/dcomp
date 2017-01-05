#!/usr/bin/python

import os
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

failed=[]
include=[]
exclude=[]

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



def check_test(path,fname,exit_code,output):
    fullname = os.path.join(dirpath,fname)
    with open(fullname, 'r') as stream:
            struct=yaml.load(stream)

    output_verbose = output


    if struct['exit_code'] !=  exit_code:
        print "\nexit_code: " ,  exit_code,
        print output_verbose
        return False

    exactly=True
    if 'exactly' in struct['output']:
        exactly=struct['output']['exactly']

    data=""
    if 'string' in struct['output']:
        data=struct['output']['string']
    elif 'file' in struct['output']:
        fullname = os.path.join(dirpath,struct['output']['file'])
        with open(fullname, 'r') as myfile:
            data=myfile.read()

    data=data.lower().replace(" ", "").replace('\n','')

    output=re.sub(r'^\+.*\n?', '', output, flags=re.MULTILINE)

    output=output.lower().replace(" ", "").replace('\n','')

    if exactly:
        ok = data == output
    else:
        ok = output in data or data in output

    if not ok:
        print "\nexpected: " +  data
        print "got: " +  output
        print output_verbose

    return ok

for dirpath, dirs, files in os.walk("."):
    if dirpath in exclude:
        continue
    if len(include)>0 and dirpath not in include:
        continue

    path=os.path.abspath(dirpath)


    if "run.sh" in files:
        print "Testing " + dirpath + "...",
        if "init.sh" in files:
            fname = os.path.join(path,"init.sh")
            subprocess.call([fname],stdout=DEVNULL,cwd=path)
        fname = os.path.join(path,"run.sh")
        p = Popen([fname], stdin=PIPE, stdout=PIPE, stderr=subprocess.STDOUT,cwd=path)
        output, err = p.communicate()
        exit_code = p.returncode

        result=check_test(path,"check.yaml",exit_code,output)
        if "cleanup.sh" in files:
            fname = os.path.join(path,"cleanup.sh")
            subprocess.call([fname],cwd=path)
        if result:
            print bcolors.OKGREEN + "OK" + bcolors.ENDC
        else:
            print bcolors.FAIL + "FAIL" + bcolors.ENDC
            failed.append(dirpath)

if len(failed)>0:
    print bcolors.FAIL + "FAILED:" + bcolors.ENDC
    for failure in failed:
        print bcolors.FAIL + failure + bcolors.ENDC
    exit(1)
else:
    print bcolors.OKGREEN + "ALL OK" + bcolors.ENDC
    exit(0)