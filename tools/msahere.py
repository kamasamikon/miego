#!/usr/bin/env python3

import subprocess
import os
import sys

def serviceNameGet():
    with open("msa.cfg", "rt") as f:
        for line in f.readlines():
            if line.startswith("s:/ms/name="):
                serviceName = line[11:].strip()
                return serviceName
    return None

def currentDir():
    return os.path.split(os.path.realpath(__file__))[0]

def msbIPAddress():
    cmd = ("sudo", "docker", "inspect", "--format", "{{ .NetworkSettings.IPAddress }}", "msb")
    print(" ".join(cmd))
    return subprocess.check_output(cmd).strip().decode("utf-8")


def dockerRun(name, pwd, msbIP, backrun):
    cmd = ["sudo", "docker", "run", "-it", "--name", name]
    if backrun:
        cmd.extend(["-d"])
    cmd.extend(["-v", "/tmp/.conf.%s:/tmp/conf" % name])
    cmd.extend(["-v", "%s:/root/ms" % pwd])
    cmd.extend(["-e", "MSBHOST=%s" % msbIP])
    cmd.extend(["msa"])
    print(" ".join(cmd))
    return subprocess.call(cmd)

def dockerKill(name):
    cmd = ("sudo", "docker", "rm", "-f", name)
    print(" ".join(cmd))
    return subprocess.call(cmd)

name = serviceNameGet()
pwd = currentDir()
msbIP = msbIPAddress()
print("MSBIP: <%s>" % msbIP)

if "k" in sys.argv:
    dockerKill(name)

backrun = "b" in sys.argv

dockerRun(name, pwd, msbIP, backrun)

