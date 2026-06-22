import os
import subprocess
from multiprocessing.pool import ThreadPool
import pandas as pd
from termcolor import colored
#from report.report import Report
from exps.experiment import Experiment, run_benchmark_on_simulator
from exps.native import run_opencl_benchmark_on_gpu
from exps.data import create_output_folder, initialize_gold_data, save_gold_data


class PCExp(Experiment):

    def __init__(self):
        super(PCExp, self).__init__()
        self.path = os.path.dirname(os.path.realpath(__file__))
        self.__param1 = [2**i for i in range(5, 26)]
        self.__param2 = range(1, 2)
        # self.__param1 = [2**i for i in range(0, 2)]
        # self.__param2 = range(1, 5)

    def name(self):
        return "pointer_chasing"

    def desc(self):
        return ""

    def param1(self):
        return self.__param1

    def param2(self):
        return self.__param2

    def param1_name(self):
        return 'height of the matrix'

    def param2_name(self):
        return 'width of the matrix'

    def gold(self):
        d = pd.read_csv(self.path+'/gold.csv')
        return d

    def __compile_native(self):
        # path = self.path + '/native' + '/linear'
        path = self.path + '/native'
        p = subprocess.Popen('make', shell=True, cwd=path)
        p.wait()
        if p.returncode == 0:
            print(colored("Compiled " + path, 'green'))
            return False
        else:
            print(colored("Compile failed " + path, 'red'))
            return True

    def __compile_simulation(self):
        path = self.path + '/sim'
        p = subprocess.Popen('go build', shell=True,
                             cwd=path)
        p.wait()
        if p.returncode == 0:
            print(colored("Compiled " + path, 'green'))
            return False
        else:
            print(colored("Compile failed " + path, 'red'))
            return True

    def run_simulation(self, thread_pool, report, results, output_folder):
        self.__compile_simulation()
        for p1 in self.__param1:
            for p2 in self.__param2:
                args = (self.name(),
                        './sim -loop-count={0} -num-wf={1}',
                        self.path+'/sim/', p1, p2, output_folder, report)
                result = thread_pool.apply_async(
                    run_benchmark_on_simulator, args)
                results.append(result)

                # Uncomment the next line if you want to run the simulations
                # one by one.
                # run_benchmark_on_simulator(*args)

    def run_native(self, repeat, output_folder, data):
        self.__compile_native()
        warmup = int(repeat/2)
        for p1 in self.__param1:
            for p2 in self.__param2:
                for i in range(0, warmup):
                    time = self.run_native_once(p1, p2, output_folder)
                for i in range(0, repeat-warmup):
                    time = self.run_native_once(p1, p2, output_folder)
                    entry = [self.name(), 'native', p1, p2, time]
                    data.loc[len(data)] = entry

    def run_native_once(self, param1, param2, output_folder):
        # return run_opencl_benchmark_on_gpu(self.name(),
        #                                 './kernel {0}',
        #                                 self.path + '/native' + '/linear',
        #                                 param1, param2, output_folder)
        return run_opencl_benchmark_on_gpu(self.name(),
                                        './kernel {0}',
                                        self.path + '/native' + '/random',
                                        param1, param2, output_folder)


if __name__ == '__main__':
    # path = os.path.dirname(os.path.realpath(__file__))
    # output_folder = create_output_folder()

    # exp = FlatLoadExp()

    # data = initialize_gold_data()
    # exp.run_native(5, output_folder, data)
    # save_gold_data(data, output_folder)

    '''report = Report([exp], output_folder, exp.gold())
    report.initialize()

    tp = ThreadPool()
    results = []
    exp.run_simulation(tp, report, results, output_folder)
    tp.close()
    tp.join()'''
