""" util functions for native execution """

import os
import re
import pandas as pd
import subprocess


def get_exp_folder(exp_name, param1, param2, output_folder):
    folder = output_folder
    create_folder_if_not_exist(folder)

    folder += '/gold_out/'
    create_folder_if_not_exist(folder)

    folder += exp_name + '/'
    create_folder_if_not_exist(folder)

    folder += str(param1) + '/'
    create_folder_if_not_exist(folder)

    folder += str(param2) + '/'
    create_folder_if_not_exist(folder)

    return folder


def run_hsa_benchmark_on_gpu(
    exp_name, exe, path, param1, param2, output_folder
):
    folder = get_exp_folder(exp_name, param1, param2, output_folder)

    stdout_filename = folder + '/stdout.out'

    fp = open(stdout_filename, 'w')
    trace_filename = folder + 'results.csv'
    cwd = path
    cmd = 'rocprof --stats -o ' + trace_filename + ' '
    cmd += exe.format(param1, param2)

    print("Executing: cd " + cwd + " && " + cmd)

    p = subprocess.Popen(cmd, shell=True, cwd=cwd, stdout=fp, stderr=fp)
    p.wait()

    kernel_time = parse_kernel_time_hsa(trace_filename)
    print("Kernel time: " + str(kernel_time))

    return kernel_time


def parse_kernel_time_hsa(filename):
    df = pd.read_csv(filename)

    kernel_time = df['DurationNs'][5]
    kernel_time /= 1e9

    return kernel_time


def run_opencl_benchmark_on_gpu(exp_name, exe, path, param1, param2, output_folder):
    create_folder_if_not_exist(output_folder)
    create_folder_if_not_exist(output_folder+'/gold_out/')

    output_prefix = output_folder + 'gold_out/' + \
        output_filename(exp_name, param1, param2)
    stdout_filename = output_prefix + '_stdout.out'
    fp = open(stdout_filename, 'w')

    cwd = path
    cmd = exe.format(param1, param2)

    print("Executing: cd " + cwd + " && " + cmd)

    p = subprocess.Popen(cmd, shell = True, cwd = cwd, stdout = fp, stderr = fp)
    p.wait()

    kernel_time = parse_kernel_time_opencl(stdout_filename)
    print("Kernel time: " + str(kernel_time))

    return kernel_time
    # m = re.search(r'kernel time: ([0-9\.]+)', str(stdout)) 
    # return float(m.group(1))
    
def parse_kernel_time_opencl(filename):
    pattern = re.compile(r'kernel time: ([0-9\.]+)')

    fp = open(filename, 'r')
    kernel_time = 0
    for line in fp:
        match = pattern.match(line)
        # print(line, match)
        if match != None:
            duration = float(match.group(1))
            kernel_time += duration
            # print(kernel_time, duration)
    return kernel_time


def output_filename(exp_name, param1, param2):
    return exp_name + '_' + str(param1) + '_' + str(param2)


def create_folder_if_not_exist(path):
    if not os.path.exists(path):
        os.makedirs(path)
