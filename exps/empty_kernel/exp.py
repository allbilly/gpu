import subprocess
import re
import os
from exps.experiment import Experiment, run_benchmark_on_simulator
import pandas as pd
from termcolor import colored
from exps.native import run_opencl_benchmark_on_gpu
from exps.data import create_output_folder

class EmptyKernelExp(Experiment):
    def __init__(self):
        super(EmptyKernelExp,self).__init__()
        self.path = os.path.dirname(os.path.realpath(__file__))
        self.__param1 = [2**i for i in range(0,15)]
        self.__param2 = range(1,9)
    def name(self):
        return "empty_kernel"

    def desc(self):
        return "The Empty Kernel microbenchmark tests the GPU's wavefront " \
            "dispatching performance. The x-axis is the number of work-group to " \
            "dispatch while the line series represents the number of wavefronts " \
            "per work-group."

    def param1(self):
        return self.__param1

    def param2(self):
        return self.__param2

    def param1_name(self):
        return 'Number Work-Groups'

    def param2_name(self):
        return 'Number Wavefronts Per Work-Group'

    def gold(self):
        d = pd.read_csv(self.path+'/gold.csv')
        return d
    
    def __compile_native(self):
        path = self.path + '/native'
        p = subprocess.Popen('make', shell=True, cwd=path)
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
                        './sim -num-wg={0} -num-wf-per-wg={1}',
                        self.path+'/sim/', p1, p2, output_folder, report)
                result = thread_pool.apply_async(
                    run_benchmark_on_simulator, args)
                results.append(result)

                # Uncomment the next line if you want to run the simulations
                # one by one.
                # run_benchmark_on_simulator(*args)
    def run_native(self, repeat, output_folder, data):
        self.__compile_native()
        for p1 in self.__param1:
            for p2 in self.__param2:
                for i in range(0, repeat):
                    time = self.run_native_once(p1,p2,output_folder)
                    entry = [self.name(), "native", p1, p2, time]
                    data.loc[len(data)] =entry

    def run_native_once(self, param1, param2, output_folder):
        return run_opencl_benchmark_on_gpu(self.name(),'./kernel {0} {1}',self.path + '/native',param1,param2, output_folder)    

if __name__ == '__main__':
    path = os.path.dirname(os.path.realpath(__file__))
    output_folder = create_output_folder()

    exp = EmptyKernelExp()

    kt = exp.run_native_once(1,1,output_folder)
    print("kernel time = " + str(kt))
