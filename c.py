
#coding: utf-8

import csv
import re
import sys

from datetime import datetime, timedelta


def read(fname):
    with open(fname) as f:
        c = f.read()
    return c


def count_zadd(t):
    ptn = re.compile("zadd:\s(\d{1,10}.\d{1,10}).+?(\d{1,10}.\d{1,10})",
            re.DOTALL)
    return ptn.findall(t)


def count_set(t):
    ptn = re.compile("set:\s(\d{1,10}.\d{1,10}).+?(\d{1,10}.\d{1,10})",
            re.DOTALL)
    return ptn.findall(t)


def extra_cpumem(t):
    ptn = re.compile("(\d{1,3}.\d)\s(\d{1,3}.\d)", re.DOTALL)
    return ptn.findall(t)

def extra_io(t):
    ptn = re.compile("(\d{1,3}\.\d{2})", re.DOTALL)
    return ptn.findall(t)

def main():
    if len(sys.argv) != 2:
        print "usage: python %s filename" % sys.argv[0]
        sys.exit()

    fn = sys.argv[1]
    content = read(fn) 
    if fn.startswith("report"):
        z = count_zadd(content)
        time_z = []
        req_z = []
        for i in z:
            req_z.append(int(float(i[0])))
            time_z.append(int(float(i[1])))
        print "sum of zadd time is %d" % sum(time_z)
        s = count_set(content)
        time_s = []
        req_s = []
        for i in s:
            req_s.append(int(float(i[0])))
            time_s.append(int(float(i[1])))
        print "sum of set time is %d" % sum(time_s)
        print "------- time--------"
        print time_z
        print time_s
        print "*" * 80

        print "------request-------"
        print req_z
        print req_s

    elif fn.startswith("cpumem"):
        c = extra_cpumem(content)
        cpu = []
        mem = []
        for i in c:
            cpu.append(float(i[0]))
            mem.append(float(i[1]))

        print "---------cpu-----"
        print cpu
        print "-------mem------"
        print mem

    elif fn.startswith("io"):
        c = extra_io(content)
        now = datetime.now()
        delta = timedelta(days=1)

        with open("io_level.tsv", "w") as f:
            writer = csv.writer(f, delimiter='\t')
            writer.writerow(["date", "close"])
            for i in c:
                now += delta
                writer.writerow([now.strftime("%d-%b-%y"), float(i)]) 
            

    else:
        print "unsupported fname"
        



if __name__ == "__main__":
    main()

