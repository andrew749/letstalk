"""
To use, add to datadog.conf as follows:
    dogstreams: [path to ngnix log (e.g: "/var/log/nginx/access.log"]:[path to this python script (e.g "/usr/share/datadog/agent/dogstream/nginx.py")]:[name of parsing method of this file ("parse")]
so, an example line would be:
    dogstreams: /var/log/nginx/access.log:/usr/share/datadog/agent/dogstream/nginx.py:parse
Log of nginx should be defined like that:
    log_format time_log '$time_local "$request" S=$status $bytes_sent T=$request_time R=$http_x_forwarded_for';
when starting dd-agent, you can find the collector.log and check if the dogstream initialized successfully
"""

from argparse import ArgumentParser
from datetime import datetime
import time
import re

# mapping between datadog and supervisord log levels #tweaked by technovangelist to add 1xx-4xx status codes
METRIC_TYPES = {
    'AVERAGE_RESPONSE': 'nginx.net.avg_response',
    'FIVE_HUNDRED_STATUS': 'nginx.net.5xx_status',
    'FOUR_HUNDRED_STATUS': 'nginx.net.4xx_status',
    'THREE_HUNDRED_STATUS': 'nginx.net.3xx_status',
    'TWO_HUNDRED_STATUS': 'nginx.net.2xx_status',
    'ONE_HUNDRED_STATUS': 'nginx.net.1xx_status'
}

TIME_REGEX = "\sT=[-+]?[0-9]*\.?[0-9]+\s*"
TIME_REGEX_SPLIT = re.compile("T=")
STATUSFIVE_REGEX = "\sS=+5[0-9]{2}\s"
STATUSFOUR_REGEX = "\sS=+4[0-9]{2}\s"
STATUSTHREE_REGEX = "\sS=+3[0-9]{2}\s"
STATUSTWO_REGEX = "\sS=+2[0-9]{2}\s"
STATUSONE_REGEX = "\sS=+1[0-9]{2}\s"

def parse(log, line):
    if len(line) == 0:
        log.info("Skipping empty line")
        return None
    timestamp = getTimestamp(line)
    avgTime = parseAvgTime(line)
    objToReturn = []
    if isHttpResponse5XX(line):
        objToReturn.append((METRIC_TYPES['FIVE_HUNDRED_STATUS'], timestamp, 1, {'metric_type': 'counter'}))
    if isHttpResponse4XX(line):
        objToReturn.append((METRIC_TYPES['FOUR_HUNDRED_STATUS'], timestamp, 1, {'metric_type': 'counter'}))
    if isHttpResponse3XX(line):
        objToReturn.append((METRIC_TYPES['THREE_HUNDRED_STATUS'], timestamp, 1, {'metric_type': 'counter'}))
    if isHttpResponse2XX(line):
        objToReturn.append((METRIC_TYPES['TWO_HUNDRED_STATUS'], timestamp, 1, {'metric_type': 'counter'}))
    if isHttpResponse1XX(line):
        objToReturn.append((METRIC_TYPES['ONE_HUNDRED_STATUS'], timestamp, 1, {'metric_type': 'counter'}))
    if avgTime is not None:
        objToReturn.append((METRIC_TYPES['AVERAGE_RESPONSE'], timestamp, avgTime, {'metric_type': 'gauge'}))
    return objToReturn

def getTimestamp(line):
    line_parts = line.split()
    dt = line_parts[0]
    date = datetime.strptime(dt, "%d/%b/%Y:%H:%M:%S")
    date = time.mktime(date.timetuple())
    return date

def parseAvgTime(line):
    time = re.search(TIME_REGEX, line)
    if time is not None:
        time = time.group(0)
        time = TIME_REGEX_SPLIT.split(time)
        if len(time) == 2:
            return float(time[1])
        return None

def isHttpResponse5XX(line):
    response = re.search(STATUSFIVE_REGEX, line)
    return (response is not None)
def isHttpResponse4XX(line):
    response = re.search(STATUSFOUR_REGEX, line)
    return (response is not None)
def isHttpResponse3XX(line):
    response = re.search(STATUSTHREE_REGEX, line)
    return (response is not None)
def isHttpResponse2XX(line):
    response = re.search(STATUSTWO_REGEX, line)
    return (response is not None)
def isHttpResponse1XX(line):
    response = re.search(STATUSONE_REGEX, line)
    return (response is not None)

if __name__ == "__main__":
    import sys
    import pprint
    import logging
    logging.basicConfig()
    log = logging.getLogger()
    lines = open(sys.argv[1]).readlines()
    pprint.pprint([parse(log, line) for line in lines])
