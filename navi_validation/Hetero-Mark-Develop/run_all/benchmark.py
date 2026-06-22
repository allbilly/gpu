"""Benchmark provides definitions about all the benchmarks that can be run,
and how to run those benchmarks.
"""
from __future__ import print_function
from data import initialize_gold_data, save_gold_data

import os
import subprocess
import re
import sys
import numpy as np
from bcolors import bcolors


class Benchmark(object):
    """Benchmark defines how to run a benchmark
    """

    def __init__(self, options):
        self.options = options
        self.benchmark_name = ''
        self.benchmark_platforms = []
        self.verify_run = []
        self.benchmark_runs = []

        self.avg_times = []
        self.std_devs = []

        self.executable = ''
        self.cwd = ''
        self.executable_full_path = ''

    def run(self,output_folder):
        """ runs the benchmark, report benchmarking result
        """
        for platform in self.benchmark_platforms:
            #self._define_current_executing(platform)

            # if not self._is_executable_found(platform):
            #     continue

            # if not self._verify(self.verify_run):
            #     continue

            data = self._benchmark()
            save_gold_data(data, output_folder, self.benchmark_name)

    def _define_current_executing(self, platform):
        self.executable = self.benchmark_name + "_" + platform
        self.cwd = self.options.build_folder + 'src/' + \
            self.benchmark_name + '/' + platform + '/'
        self.executable_full_path = self.cwd + self.executable

    def _is_executable_found(self, platform):
        if not os.path.isfile(self.executable_full_path):
            print(self.executable_full_path, 'not found, skip.')
            return False
        return True

    def _verify(self, args):
        print("Verifying", self.executable, *args, sep=' ', end=' ')
        sys.stdout.flush()
        command = " ".join([self.executable_full_path, '-q', '-v'] + args)
        proc = subprocess.Popen(command,
                                cwd=self.cwd, shell=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        output = proc.communicate()
        if proc.returncode != 0:
            print(bcolors.FAIL, "error: ", self.executable, bcolors.ENDC,
                  sep='')
            print(output)
            sys.stdout.flush()
            return False

        print(bcolors.OKGREEN, "Passed", bcolors.ENDC, sep='')
        sys.stdout.flush()
        return True

    def _benchmark(self):
        data = initialize_gold_data()
        for run in self.benchmark_runs:
            if self.benchmark_name in ['fir','pr','kmeans']:
                avg = self._benchmark_specific_input(run)
                entry = [self.benchmark_name, "native", str(run[-1]), avg]
            else:
                avg = self._benchmark_rest(run)
                entry = [self.benchmark_name, "native", str(run), avg]
            data.loc[len(data)] = entry
        return data
        
    
    def _benchmark_rest(self, args):
        print("Benchmarking", self.benchmark_name, args)
        pattern = re.compile(r'GPU Runtime: ([0-9\.]+)')
        perf = []
        for i in range(0, self.options.repeat_time + 5):
            #command = " ".join([self.executable_full_path] + args)
            
            p = subprocess.Popen('./kernel '+str(args), shell = True, cwd = self.cwd, stdout=subprocess.PIPE)
            (stdout,_) = p.communicate()
            res = pattern.search(str(stdout))
            if res and i>4:
                perf.append(float(res.group(1))) 

        if self.options.repeat_time > 0:
            return np.mean(perf)

    def _benchmark_specific_input(self, args):

        if self.options.full_verification:
            self._verify(args)

        print("Benchmarking", self.executable, *args, sep=' ', end=' ')
        sys.stdout.flush()

        runtime_regex = re.compile(
            r'Run: (-?\ *[0-9]+\.?[0-9]*(?:[Ee]\ *-?\ *[0-9]+)?) second')

        perf = []
        for i in range(0, self.options.repeat_time):
            command = " ".join([self.executable_full_path] + args)
            # print(command)
            proc = subprocess.Popen(command,
                                    cwd=self.cwd, shell=True,
                                    stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            for line in proc.stderr:
                res = runtime_regex.search(str(line))
                if res:
                    perf.append(float(res.group(1)))
            print(".", end='')
            sys.stdout.flush()
        print("")

        if self.options.repeat_time > 0:
            return np.mean(perf)

    def _report(self):
        print("Benchmark ", self.benchmark_name, "results:",)
        print("\ttime: ", self.avg_times)
        print("\tstd_dev: ", self.std_devs)


class FirBenchmark(Benchmark):
    """FIR benchmark"""

    def __init__(self, options):
        super(FirBenchmark, self).__init__(options)
        self.benchmark_name = 'fir'
        self.benchmark_platforms = ['cl12']
        self.executable = 'fir_cl12'
        self.cwd = self.options.build_folder + 'src/' + \
            self.benchmark_name + '/cl12/'
        self.executable_full_path = self.cwd + self.executable +' -q -t'
        self.verify_run = []
        self.benchmark_runs = [
            ['-y', '1024', '-x', '1024'],
            ['-y', '1024', '-x', '2048'],
            ['-y', '1024', '-x', '3072'],
            ['-y', '1024', '-x', '4096'],
            ['-y', '1024', '-x', '5120'],
            ['-y', '1024', '-x', '6144'],
            ['-y', '1024', '-x', '7168'],
            ['-y', '1024', '-x', '8192'],
        ]


class AtaxBenchmark(Benchmark):
    """Atax benchmark"""

    def __init__(self, options):
        super(AtaxBenchmark, self).__init__(options)
        self.benchmark_name = 'atax'
        self.benchmark_platforms = ['cl12']
        self.verify_run = []
        self.executable = 'kernel'
        self.cwd = '../atax/native/' 
        self.executable_full_path = self.cwd + self.executable 
       
        self.benchmark_runs = [2**i for i in range(8,16)]


class PRBenchmark(Benchmark):
    """KMeans benchmark"""

    def __init__(self, options):
        super(PRBenchmark, self).__init__(options)
        self.benchmark_name = 'pr'
        self.benchmark_platforms = ['cl12']
        self.executable = 'pr_cl12'
        self.cwd = self.options.build_folder + 'src/' + \
            self.benchmark_name + '/cl12/'
        self.executable_full_path = self.cwd + self.executable +' -q -t'
        self.verify_run = ['-i', os.getcwd() + '/data/pr/1024.data']
        self.benchmark_runs = [
            ['-i', os.getcwd() + '/data/pr/1024.data'],
            ['-i', os.getcwd() + '/data/pr/2048.data'],
            ['-i', os.getcwd() + '/data/pr/4096.data'],
            ['-i', os.getcwd() + '/data/pr/8192.data'],
            ['-i', os.getcwd() + '/data/pr/16384.data'],
        ]


class KMeansBenchmark(Benchmark):
    """KMeans benchmark"""

    def __init__(self, options):
        super(KMeansBenchmark, self).__init__(options)
        self.benchmark_name = 'kmeans'
        self.benchmark_platforms = ['cl12']
        self.executable = 'kmeans_cl12'
        self.cwd = self.options.build_folder + 'src/' + \
            self.benchmark_name + '/cl12/'
        self.executable_full_path = self.cwd + self.executable +' -q -t'
        self.verify_run = ['-i', os.getcwd() + '/data/kmeans/1000_34.txt']
        self.benchmark_runs = [
            ['-i', os.getcwd() + '/data/kmeans/100_34.txt'],
            ['-i', os.getcwd() + '/data/kmeans/1000_34.txt'],
            ['-i', os.getcwd() + '/data/kmeans/10000_34.txt'],
            ['-i', os.getcwd() + '/data/kmeans/100000_34.txt'],
            ['-i', os.getcwd() + '/data/kmeans/1000000_34.txt'],
        ]


class BicgBenchmark(Benchmark):
    """Bicg benchmark"""

    def __init__(self, options):
        super(BicgBenchmark, self).__init__(options)
        self.benchmark_name = 'bicg'
        self.benchmark_platforms = ['cl12']
        self.verify_run = []
        self.executable = 'kernel'
        self.cwd = '../bicg/native/' 
        self.executable_full_path = self.cwd + self.executable
        self.benchmark_runs = [2**i for i in range(8,16)]


class BSBenchmark(Benchmark):
    """Bitonicsort benchmark"""

    def __init__(self, options):
        super(BSBenchmark, self).__init__(options)
        self.benchmark_name = 'bs'
        self.benchmark_platforms = ['cl12']
        self.verify_run = []
        self.executable = 'kernel'
        self.cwd = '../bitonicsort/native/' 
        self.executable_full_path = self.cwd + self.executable
        self.benchmark_runs = [2**i for i in range(8,20)]


class FTBenchmark(Benchmark):
    """fastwalshtransform benchmark"""

    def __init__(self, options):
        super(FTBenchmark, self).__init__(options)
        self.benchmark_name = 'ft'
        self.benchmark_platforms = ['cl12']
        self.verify_run = ['-i', os.getcwd() + '/data/be/320x180.mp4']
        self.executable = 'kernel'
        self.cwd = '../fastwalshtransform/native/' 
        self.executable_full_path = self.cwd + self.executable
        self.benchmark_runs = [2**i for i in range(8,20)] 


class FWBenchmark(Benchmark):
    """floydwarshall benchmark"""

    def __init__(self, options):
        super(FWBenchmark, self).__init__(options)
        self.benchmark_name = 'fw'
        self.benchmark_platforms = ['cl12']
        self.verify_run = ['-i', os.getcwd() + '/data/ga/65536_64.data']
        self.executable = 'kernel'
        self.cwd = '../floydwarshall/native/' 
        self.executable_full_path = self.cwd + self.executable
        self.benchmark_runs = [2**i for i in range(8,20)]

class MTBenchmark(Benchmark):
    """matrixtranspose benchmark"""

    def __init__(self, options):
        super(MTBenchmark, self).__init__(options)
        self.benchmark_name = 'mt'
        self.benchmark_platforms = ['cl12']
        self.verify_run = []
        self.executable = 'kernel'
        self.cwd = '../matrixtranspose/native/' 
        self.executable_full_path = self.cwd + self.executable
        self.benchmark_runs = [2**i for i in range(8,20)]

class ReluBenchmark(Benchmark):
    """relu benchmark"""

    def __init__(self, options):
        super(ReluBenchmark, self).__init__(options)
        self.benchmark_name = 'relu'
        self.benchmark_platforms = ['cl12']
        self.verify_run = []
        self.executable = 'kernel'
        self.cwd = '../relu/native/' 
        self.executable_full_path = self.cwd + self.executable
        self.benchmark_runs = [2**i for i in range(8,20)]