#!/usr/bin/env python3

import subprocess
import os
import sys

MSB_NAME = "msb"
MS_SUFFIX = ""
MSB_ADDR = ""

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
    return saferun(("sudo", "docker", "inspect", "--format", "{{ .NetworkSettings.IPAddress }}", MSB_NAME))

def volumeGet(imageName):
    return saferun(("sudo", "docker", "inspect", "--format", "{{ .Config.Labels.VOLUME }}", imageName))

def dockerGateway():
    cmd = ("sudo", "docker", "network", "inspect", "bridge", "--format", '{{(index .IPAM.Config 0).Gateway}}')
    return saferun(cmd)

def dockerRun(imageName, msbIP, backrun, append):
    container = imageName + MS_SUFFIX
    if append:
        index = 0
        tmpName = container
        while True:
            tmpName = container if index == 0 else container + "_%d" % index
            index += 1
            cmd = ["sudo", "docker", "ps", "-aq", "--filter", r"""name=^/%s$""" % tmpName]
            print(">>> ", " ".join(cmd))
            if not subprocess.check_output(cmd):
                break
        container = tmpName
        print(container)

    cmd = [
            "sudo", 
            "docker", 
            "run", 
            "-it", 
            "--restart=always", 
            "--log-opt", "max-size=2m",
            "--log-opt", "max-file=5",
            "--name", container,
            ]

    if backrun:
        cmd.extend(["-d"])
    cmd.extend(["-v", "/tmp/.conf.%s:/tmp/conf" % container])
    cmd.extend(["-e", "MSBHOST=%s" % msbIP])
    cmd.extend(["-e", "DOCKER_GATEWAY=%s" % dockerGateway()])

    volumeMap = volumeGet(imageName)
    if volumeMap and volumeMap[0] != "<":
        cmd.extend(["-v", volumeMap])

    cmd.extend([imageName])
    return saferun(cmd)

def killContainer(imageName, killFirst, killLast):
    container = imageName + MS_SUFFIX
    killFirst = killFirst or "0"
    killLast = killLast or "99999999999"

    cmd = ["sudo", "docker", "ps", "-aq", "--filter", r'''name=\b%s\b|\b%s_.*''' % (container, container)]
    idList = subprocess.check_output(cmd).strip().decode("utf-8").split()
    print(idList)
    if idList:
        a = int(killFirst)
        b = int(killLast)
        cmd = ["sudo", "docker", "rm", "-f"]
        cmd.extend(idList[a:b])
        saferun(cmd)

def main():
    global MSB_NAME
    global MS_SUFFIX
    global MSB_ADDR

    if len(sys.argv) == 1 or "--help" in sys.argv:
        print("Directly run msa services from the image.")
        print("It fetch the MSB's IPAddress and set to the container")
        print("Usage: mscontainer.py [-k:s:e=kill] [-b=backrun] [-a=append] [--msbName=MSBName] [--msbAddr=MSBAddr] imageNames ...")
        return

    x = os.environ.get("MSB_NAME")
    if x:
        MSB_NAME = x
    x = os.environ.get("MS_SUFFIX")
    if x:
        MS_SUFFIX = x
    x = os.environ.get("MSB_ADDR")
    if x:
        MSB_ADDR = x

    #
    # Another MSB?
    #
    for name in sys.argv[1:]:
        if name.startswith("--msbName="):
            MSB_NAME = name[6:]
            continue

        if name.startswith("--suffix="):
            MS_SUFFIX = name[9:]
            continue

        if name.startswith("--msbAddr="):
            MSB_ADDR = name[10:]

    if MSB_ADDR:
        msbIP = MSB_ADDR
    else:
        msbIP = msbIPAddress()

    backrun = "-b" in sys.argv
    append  = "-a" in sys.argv

    msNames = []
    for name in sys.argv[1:]:
        if name[0] == "-":
            continue
        msNames.append(name)

    if not msNames:
        x = os.environ.get("MS_NAME")
        if x:
            msNames.append(x)

    #
    # Kill OLD?
    #
    killSome, killFirst, killLast = False, "0", "999999999999999"
    for name in sys.argv[1:]:
        if name.startswith("-k"):
            # -k: => container[:]
            # -k: => container[:]
            # -k3:-1 => container[3:-1]
            segs = name[2:].split(":")
            if len(segs) > 1:
                killLast = segs[1]

            if len(segs) > 0:
                killFirst = segs[0]

            for xname in msNames:
                if xname[0] == "-":
                    continue
                killContainer(xname, killFirst, killLast)
            break

    for name in msNames:
        dockerRun(name, msbIP, backrun, append)

if __name__ == "__main__":
    main()

