#!/usr/bin/env python3

import argparse
from shutil import copyfile

from exps.data import create_output_folder, collect_gold_data
#from report.report import Report
from exps.experiments import Experiments
from multiprocessing.pool import ThreadPool


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-f", "--filter",  default="",
        help="regular expression to filter experiments")
    parser.add_argument(
        "-r", "--repeat", type=int, default=20,
        help="repeat times for each benchmark")
    args = parser.parse_args()

    output_folder = create_output_folder()

    exps = Experiments().get_experiments(args.filter)
    for exp in exps:
        gold_file = exp.measure_gold(output_folder, repeat=args.repeat)
        copyfile(gold_file, exp.path + "/gold.csv")


if __name__ == '__main__':
    main()
