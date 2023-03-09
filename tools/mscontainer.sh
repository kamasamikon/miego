#!/usr/bin/env python3

import subprocess
import os
import sys
import shlex

MSB_NAME = "msb"
MS_SUFFIX = ""

f = open("/tmp/mscontainer.log", "w+")

def saferun(cmd, debug=True):
    try:
        print(cmd, file=f)
        if debug:
            color = 3
            cmdline = "'" + "' '".join(cmd) + "'"
            print('\033[1;3{}m{}\033[0m'.format(color, cmdline))

        return subprocess.check_output(cmd).strip().decode("utf-8")
    except:
        return None

def currentDir():
    return os.path.realpath(os.getcwd())

def getContainerIp(containerName):
    return saferun(("sudo", "docker", "inspect", "--format", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", containerName))

def getContainerPort(containerName, innerPort=None):
    innerPort = innerPort or "80/tcp"
    return saferun(("sudo", "docker", "inspect", "--format", "{{(index (index .NetworkSettings.Ports '%s') 0).HostPort}}" % innerPort, containerName))

def volumeGet(imageName):
    return saferun(("sudo", "docker", "inspect", "--format", "{{ .Config.Labels.VOLUME }}", imageName))

def dockerGateway():
    cmd = ("sudo", "docker", "network", "inspect", "bridge", "--format", '{{(index .IPAM.Config 0).Gateway}}')
    return saferun(cmd)

def dockerRun(imageName, MSBHost, MSBPort, backrun, append):
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

    for i in range(len(sys.argv)):
        if sys.argv[i].startswith("--runopt="):
            segs = [x for x in shlex.split(sys.argv[i][9:]) if x]
            if segs:
                cmd.extend(segs)

    if backrun:
        cmd.extend(["-d"])
    cmd.extend(["-e", "MSBHOST=%s" % MSBHost])
    cmd.extend(["-e", "MSBPORT=%s" % MSBPort])
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

    if len(sys.argv) == 1 or "--help" in sys.argv:
        print("Directly run msa services from the image.")
        print("It fetch the MSB's IPAddress and set to the container")
        print("Usage: mscontainer.py [-k:s:e=kill] [-b=backrun] [-a=append] [--msbName=MSBName] imageNames ...")
        return

    x = os.environ.get("MSB_NAME") # MSB的名字
    if x:
        MSB_NAME = x
    x = os.environ.get("MS_SUFFIX")
    if x:
        MS_SUFFIX = x

    #
    # Another MSB?
    #
    for name in sys.argv[1:]:
        if name.startswith("--msbName="):
            MSB_NAME = name[10:]
            continue

        if name.startswith("--suffix="):
            MS_SUFFIX = name[9:]
            continue

    MSBHost = getContainerIp(MSB_NAME)
    MSBPort = getContainerPort(MSB_NAME)

    backrun = "-b" in sys.argv
    append  = "-a" in sys.argv

    imageNames = []
    for name in sys.argv[1:]:
        if name[0] == "-":
            continue
        imageNames.append(name)

    if not imageNames:
        x = os.environ.get("MS_NAME")
        if x:
            imageNames.append(x)

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

            for xname in imageNames:
                if xname[0] == "-":
                    continue
                killContainer(xname, killFirst, killLast)
            break

    for imageName in imageNames:
        dockerRun(imageName, MSBHost, MSBPort, backrun, append)

if __name__ == "__main__":
    main()

