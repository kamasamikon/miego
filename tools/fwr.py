#!/usr/bin/env python3

import subprocess
import os
import sys
import shlex
import time

# fuwu runner

f = open("/tmp/fwr.log", "a+")

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

def volumeGet(imageName):
    return saferun(("sudo", "-S", "docker", "inspect", "--format", "{{ .Config.Labels.VOLUME }}", imageName))

def dockerGateway():
    cmd = ("sudo", "-S", "docker", "network", "inspect", "bridge", "--format", '{{(index .IPAM.Config 0).Gateway}}')
    return saferun(cmd)

def findMSB(*args):
    cmd = ["sudo", "-S", "docker", "ps", "-aq", "--filter", "ancestor=msb"]
    for i in range(int(len(args)/2)):
        cmd.append("--filter")
        cmd.append("%s=%s" % (args[2*i], args[2*i+1]))
    x = saferun(cmd)
    return x

def dockerRun(imageName, msbName, msbPort, backrun, append):
    if not imageName:
        return
    if not msbName and not msbPort:
        return

    gw = dockerGateway()

    #
    # 检查MSB是否存在，如果不存在，就启动一个
    #
    msbRun = ["sudo", "-S", "docker", "run", "--restart=always", "-it", "-d"]

    if msbName and msbPort:
        # msbName = "msb.auv", msbPort = "7788"
        x = findMSB("name", msbName, "publish", msbPort)
        if not x:
            msbRun.extend(["--name", msbName, "-p", "%s:80" % msbPort])
            msbRun.append("msb")

    if not msbName and msbPort:
        # msbName = "", msbPort = "7788" => msbName = "msb.7788"
        x = findMSB("publish", msbPort)
        if not x:
            msbRun.extend(["-p", "%s:80" % msbPort])
            msbRun.append("msb")

    if msbName and not msbPort:
        # msbName = "msb.auv", msbPort = "" => msbName = "msb.auv"
        x = findMSB("name", msbName)
        if not x:
            msbRun.extend(["--name", msbName, "-p", ":80"])
            msbRun.append("msb")

    if len(msbRun) > 7:
        saferun(msbRun)

    #
    # 如果没有指定msbPort，就查一下，然后使用查到的第一个
    #
    if not msbPort:
        x = msbName or "msb"
        cmd = ["sudo", "-S", "docker", "port", x]
        try:
            msbPort = saferun(cmd).split("\n")[0].split(":")[-1]
        except:
            pass


    #
    # 开始启动服务
    #

    if msbName:
        postfix = (msbName + ".msb").split(".")[-2]
        container = imageName + "." + postfix
    else:
        container = imageName + ".msb." + msbPort

    if append:
        index = 0
        tmpName = container
        while True:
            tmpName = container if index == 0 else container + "_%d" % index
            index += 1
            cmd = ["sudo", "-S", "docker", "ps", "-aq", "--filter", r"""name=^/%s$""" % tmpName]
            print(">>> ", " ".join(cmd))
            if not subprocess.check_output(cmd):
                break
        container = tmpName
        print(container)

    cmd = [
            "sudo",
            "-S",
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
    cmd.extend(["-e", "MSBPORT=%s" % msbPort])
    cmd.extend(["-e", "MSBNAME=%s" % msbName])
    cmd.extend(["-e", "DOCKER_GATEWAY=%s" % gw])

    volumeMap = volumeGet(imageName)
    if volumeMap and volumeMap[0] != "<":
        cmd.extend(["-v", volumeMap])

    cmd.extend([imageName])
    return saferun(cmd)

def killContainer(imageName, msbName, killFirst, killLast):
    postfix = (msbName + ".msb").split(".")[-2]
    container = imageName + "." + postfix
    killFirst = killFirst or "0"
    killLast = killLast or "99999999999"

    cmd = ["sudo", "-S", "docker", "ps", "-aq", "--filter", r'''name=\b%s\b|\b%s_.*''' % (container, container)]
    idList = subprocess.check_output(cmd).strip().decode("utf-8").split()
    print(idList)
    if idList:
        a = int(killFirst)
        b = int(killLast)
        cmd = ["sudo", "-S", "docker", "rm", "-f"]
        cmd.extend(idList[a:b])
        saferun(cmd)

def main():
    if len(sys.argv) == 1 or "--help" in sys.argv:
        print("fwr.py [-k:s:e=kill] [-b=backrun] [-a=append] [--msbName=] [--msbPort=] imageNames ...")
        return

    #
    # Another MSB?
    #
    for name in sys.argv[1:]:
        if name.startswith("--msbName="):
            msbName = name[10:]
            continue
        if name.startswith("--msbPort="):
            msbPort = name[10:]
            continue

    backrun = "-b" in sys.argv
    append  = "-a" in sys.argv

    imageNames = []
    for name in sys.argv[1:]:
        if name[0] == "-":
            continue
        imageNames.append(name)

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
                killContainer(xname, msbName, killFirst, killLast)
            break

    for imageName in imageNames:
        dockerRun(imageName, msbName, msbPort, backrun, append)

if __name__ == "__main__":
    print("--- %s ---" % time.asctime(), file=f)
    print(sys.argv, file=f)
    main()

