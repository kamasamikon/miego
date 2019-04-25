#!/usr/bin/env python3

import subprocess
import os
import sys


def saferun(cmd, debug=True):
    try:
        if debug:
            color = 3
            cmdline = " ".join(cmd)
            print('\033[1;3{}m{}\033[0m'.format(color, cmdline))

        return subprocess.check_output(cmd).strip().decode("utf-8")
    except:
        return None


def serviceNameGet():
    with open("msa.cfg", "rt") as f:
        for line in f.readlines():
            if line.startswith("s:/ms/name="):
                serviceName = line[11:].strip()
                return serviceName
    return None

def currentDir():
    return os.path.realpath(os.getcwd())

def msbIPAddress():
    return saferun(("sudo", "docker", "inspect", "--format", "{{ .NetworkSettings.IPAddress }}", "msb"))

def dockerRun(name, pwd, msbIP, backrun):
    cmd = ["sudo", "docker", "run", "-it", "--name", name]
    if backrun:
        cmd.extend(["-d"])
    cmd.extend(["-v", "/tmp/.conf.%s:/tmp/conf" % name])
    cmd.extend(["-v", "%s:/root/ms" % pwd])
    cmd.extend(["-e", "MSBHOST=%s" % msbIP])
    cmd.extend(["msa"])
    return saferun(cmd)

def dockerKill(name):
    saferun(("sudo", "docker", "rm", "-f", name))

def main():
    if "--help" in sys.argv:
        print("1. Fetch the MSB's IPAddress and set to the container")
        print("2. Get the container name from ./msa.cfg")
        print("3. Share current directory to container's /root/ms")
        print("4. Run the docker container")
        print("Usage: msahere.py [-k:kill] [-b:backrun]")
        return

    name = serviceNameGet()
    pwd = currentDir()
    msbIP = msbIPAddress()

    backrun = "-b" in sys.argv
    killold = "-k" in sys.argv

    if killold:
        dockerKill(name)

    dockerRun(name, pwd, msbIP, backrun)

if __name__ == "__main__":
    main()

