from argparse import ArgumentParser
from datetime import datetime
import time
import re
from datadog import ThreadStats, initialize
from functools import partial
import tailer
import logging

initialize(statsd_use_default_route=True)
stats = ThreadStats()
stats.start()

# setup datadog to talk to local statsd instance and start thread
nginx_status_format = "nginx.net.status.{}"
nginx_average_response_format = "nginx.net.avg_response"

TIME_REGEX = "\sT=[-+]?[0-9]*\.?[0-9]+\s*"
TIME_REGEX_SPLIT = re.compile("T=")
STATUS_REGEX = "\sS=([12345]\d{2})\s"

logging.basicConfig(level=logging.DEBUG)
log = logging.getLogger()

def process(inc_function, gauge_function, line):
    if len(line) == 0:
        log.info("Skipping empty line")
        return None
    timestamp = getTimestamp(line)
    avgTime = parseAvgTime(line)
    status = parseStatus(line)
    if status is not None:
        inc_function(
            nginx_status_format.format(status),
            timestamp,
        )
    else:
        log.warn("Unable to parse status for line {}".format(line))

    if avgTime is not None:
        gauge_function(nginx_average_response_format, timestamp, avgTime)
    else:
        log.warn("Unable ot parse avgTime for line {}".format(line))

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

def parseStatus(line):
    status_res = re.search(STATUS_REGEX, line)
    return status_res.group(1)

def inc_metric(datadog, metric, timestamp, tags={}):
    if datadog:
        stats.increment(metric, tags=tags, timestamp=timestamp)
    log.info("INC {} tags=({})".format(metric, tags))

def gauge_metric(datadog, metric, timestamp, value=1):
    if datadog:
        stats.gauge(metric, tags=tags, value=value, timestamp=timestamp)
    log.info("GAUGE {} value=({})".format(metric, value))

def parse_args():
    parser = ArgumentParser(description="Tail nginx logs")
    parser.add_argument("--datadog", action="store_true", help="Whether to send to datadog.")
    parser.add_argument("--file", required=True, help="File to parse")
    return parser.parse_args()

def main():
    args = parse_args()
    inc_function = partial(inc_metric, args.datadog)
    gauge_function = partial(gauge_metric, args.datadog)
    process_func = partial(process, inc_function, gauge_function)
    log.info("Initializing tailer ...")
    for line in tailer.follow(open(args.file)):
        log.debug("Logging line {}".format(line))
        process_func(line)

if __name__ == "__main__":
    main()
