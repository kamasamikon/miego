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

def dockerGateway():
    cmd = ("sudo", "docker", "network", "inspect", "bridge", "--format", '{{(index .IPAM.Config 0).Gateway}}')
    return saferun(cmd)

def dockerRun(name, msbIP, backrun):
    cmd = [
            "sudo", 
            "docker", 
            "run", 
            "-it", 
            "--restart=always", 
            "--log-opt", "max-size=2m",
            "--log-opt", "max-file=5",
            "--name", name,
            ]
    if backrun:
        cmd.extend(["-d"])
    cmd.extend(["-v", "/tmp/.conf.%s:/tmp/conf" % name])
    cmd.extend(["-e", "MSBHOST=%s" % msbIP])
    cmd.extend(["-e", "DOCKER_GATEWAY=%s" % dockerGateway()])
    cmd.extend([name])
    return saferun(cmd)

def dockerKill(name):
    saferun(("sudo", "docker", "rm", "-f", name))

def main():
    if len(sys.argv) == 1 or "--help" in sys.argv:
        print("Directly run msa services from the image.")
        print("It fetch the MSB's IPAddress and set to the container")
        print("Usage: mscontainer.py [-k:kill] [-b:backrun] imageNames ...")
        return

    msbIP = msbIPAddress()

    backrun = "-b" in sys.argv
    killold = "-k" in sys.argv

    for name in sys.argv[1:]:
        if name[0] == "-":
            continue
        if killold:
            dockerKill(name)
        dockerRun(name, msbIP, backrun)

if __name__ == "__main__":
    main()

