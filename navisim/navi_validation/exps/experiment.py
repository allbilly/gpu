""" Defines base class for all the experiments """

from exps.data import gold_file_name, initialize_gold_data, save_gold_data
import re
import subprocess
import abc
from exps.native import create_folder_if_not_exist, output_filename


class Experiment(metaclass=abc.ABCMeta):
    """ Experiment defines the public APIs that are supported by all the 
        expriments.
    """

    def __init__(self):
        self.path = ""

    @abc.abstractmethod
    def name(self):
        pass

    @abc.abstractmethod
    def desc(self):
        pass

    @abc.abstractmethod
    def param1(self):
        pass

    @abc.abstractmethod
    def param2(self):
        pass

    @abc.abstractmethod
    def param1_name(self):
        pass

    @abc.abstractmethod
    def param2_name(self):
        pass

    @abc.abstractmethod
    def gold(self):
        pass

    def measure_gold(self, output_folder, repeat=10):
        data = initialize_gold_data()
        self.run_native(repeat, output_folder, data)
        save_gold_data(data, output_folder, self.name())
        return gold_file_name(output_folder, self.name())

    @abc.abstractmethod
    def run_simulation(self, thread_pool, report, results, output_folder):
        pass

    @abc.abstractmethod
    def run_native(self, repeat, output_folder, data):
        pass


def run_benchmark_on_simulator(name, exe, path, param1, param2, output_folder, report):
    create_folder_if_not_exist(output_folder)
    create_folder_if_not_exist(output_folder + '/sim_out/')

    cwd = path
    cmd = exe.format(param1, param2)
    filename = output_folder + '/sim_out/' + \
        output_filename(name, param1, param2) + '_sim_stdout.out'
    # print(filename)

    if report is not None:
        report.start_run(name, param1, param2)

    fp = open(filename, 'w')
    fp.write("Executing: " + cmd + "\n")
    fp.flush()

    process = subprocess.Popen(cmd, shell=True, cwd=cwd, stdout=fp, stderr=fp)
    process.wait()

    if process.returncode != 0:
        print("Error executing ", cmd)

    time = parse_kernel_time_sim(filename)

    entry = [name, 'sim', param1, param2, time]

    if report is not None:
        report.done_run(name, param1, param2, time)

    return entry


def parse_kernel_time_sim(filename):
    fp = open(filename, 'r')
    kernel_time = 0
    for line in fp:
        m = re.match(r'Kernel time: ([0-9.]+)', str(line))
        if m != None:
            kernel_time += float(m.group(1))
            return kernel_time
