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


def currentDir():
    return os.path.realpath(os.getcwd())

def msbIPAddress():
    return saferun(("sudo", "docker", "inspect", "--format", "{{ .NetworkSettings.IPAddress }}", "msb"))

def dockerRun(name, msbIP, backrun):
    cmd = ["sudo", "docker", "run", "-it", "--name", name]
    if backrun:
        cmd.extend(["-d"])
    cmd.extend(["-v", "/tmp/.conf.%s:/tmp/conf" % name])
    cmd.extend(["-e", "MSBHOST=%s" % msbIP])
    cmd.extend([name])
    return saferun(cmd)

def dockerKill(name):
    saferun(("sudo", "docker", "rm", "-f", name))

def main():
    if "--help" in sys.argv:
        print("Usage: mscontainer.py imageName [-k:kill] [-b:backrun]")
        return

    name = sys.argv[1]
    msbIP = msbIPAddress()

    if "-k" in sys.argv:
        dockerKill(name)

    backrun = "-b" in sys.argv

    dockerRun(name, msbIP, backrun)

if __name__ == "__main__":
    main()

