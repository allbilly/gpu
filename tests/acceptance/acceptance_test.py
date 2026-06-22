import subprocess
import os
import sys
import argparse
from termcolor import cprint


def compile(path):
    fp = open(os.devnull, 'w')
    p = subprocess.Popen('go build', shell=True,
                         cwd=path, stdout=fp, stderr=fp)
    p.wait()
    if p.returncode == 0:
        cprint("Compiled " + path, 'green')
        return False
    else:
        cprint("Compile failed " + path, 'red')
        return True


class Test(object):
    """ define a benchmark to test """

    def __init__(self, path, executable, size_args, benchmark_path):
        self.path = path
        self.executable = executable
        self.size_args = size_args
        self.benchmark_path = benchmark_path

    def test(self,
             test_disassemble=False,
             ):
        err = False

        if test_disassemble:
            err |= self.test_disassemble()
            return err

        err |= self.compile()

        err |= self.run_test(False, '1')
        err |= self.run_test(True, '1')
        

        return err

    def compile(self):
        fp = open(os.devnull, 'w')
        p = subprocess.Popen('go build', shell=True,
                             cwd=self.path, stdout=fp, stderr=fp)
        p.wait()
        if p.returncode == 0:
            cprint("Compiled " + self.path, 'green')
            return False
        else:
            cprint("Compile failed " + self.path, 'red')
            return True

    def run_test(self,
                 timing,
                 gpus):
        fp = open(os.devnull, 'w')
        cmd = ['./'+self.executable, '-verify']
        cmd.extend(self.size_args)


        cmd.append('-gpus='+gpus)

        if timing:
            cmd.append('-timing')

        
        cmd_string = 'cd ' + self.path + ' && ' + ' '.join(cmd)
        print('Running ' + cmd_string)

        p = subprocess.Popen(cmd, shell=False,
                             cwd=self.path,
                             stdout=fp, stderr=fp
                             )
        p.wait()

        if p.returncode == 0:
            cprint('Passed.', 'green')
            return False
        else:
            cprint('Failed.', 'red')
            return True

    def test_disassemble(self):
        output_filename = self.benchmark_path + '/disasm.disasm'
        fp = open(output_filename, 'w')
        cmd = ['../../insts/gcn3disassembler/gcn3disassembler',
               self.benchmark_path + '/kernels.hsaco']

        cmd_string = ' '.join(cmd)
        print('Running ' + cmd_string + ' > ' +
              output_filename + ' 2>&1 ' + ' ')

        p = subprocess.Popen(cmd, shell=False,
                             stdout=fp, stderr=fp
                             )
        p.wait()
        if p.returncode != 0:
            cprint(' Failed.', 'red')
            return True

        fp = open(self.benchmark_path + '/diff.debug', 'w')
        cmd = ['diff', 'kernels.disasm', 'disasm.disasm']
        p = subprocess.Popen(cmd, shell=False,
                             cwd=self.benchmark_path,
                             stdout=fp, stderr=fp)
        p.wait()

        if p.returncode == 0:
            cprint('Passed.', 'green')
            return False
        else:
            cprint('Failed.', 'red')
            return True





def main():

    atax = Test('../../samples/atax',
                'atax',
                ['-x=256', '-y=256'],
                '../../benchmarks/polybench/atax'
                )
    bicg = Test('../../samples/bicg',
                'bicg',
                ['-x=256', '-y=256'],
                '../../benchmarks/polybench/bicg'
                )
    fir = Test('../../samples/fir',
               'fir',
               ['-length=8192'],
               '../../benchmarks/heteromark/fir')
    aes = Test('../../samples/aes',
               'aes',
               ['-length=16384'],
               '../../benchmarks/heteromark/aes')
    km = Test('../../samples/kmeans', 'kmeans',
              [
                  '-points=1024',
                  '-features=32',
                  '-clusters=5',
                  '-max-iter=5'
              ],
              '../../benchmarks/heteromark/kmeans')
    pagerank = Test('../../samples/pagerank', 'pagerank',
                    [
                        '-node=64',
                        '-sparsity=0.5',
                        '-iterations=2',
                    ],
                    '../../benchmarks/heteromark/pagerank')
    mm = Test('../../samples/matrixmultiplication',
              'matrixmultiplication',
              ['-x=128', '-y=128', '-z=128'],
              '../../benchmarks/amdappsdk/matrixmultiplication')
    mt = Test('../../samples/matrixtranspose',
              'matrixtranspose',
              ['-width=256'],
              '../../benchmarks/amdappsdk/matrixtranspose')
    bs = Test('../../samples/bitonicsort',
              'bitonicsort',
              ['-length=4096'],
              '../../benchmarks/amdappsdk/bitonicsort')
    sc = Test('../../samples/simpleconvolution',
              'simpleconvolution',
              [],
              '../../benchmarks/amdappsdk/simpleconvolution')
    fw = Test('../../samples/floydwarshall',
              'floydwarshall',
              [],
              '../../benchmarks/amdappsdk/floydwarshall')
    re = Test('../../samples/relu',
              'relu',
              [],
              '../../benchmarks/dnn/relu')
    mp = Test('../../samples/maxpooling',
              'maxpooling',
              [
                  '-h=128',
                  '-w=128',
                  '-c=3',
                  '-n=4',
                  '-stride-w=2',
                  '-stride-h=2',
                  '-kernel-w=2',
                  '-kernel-h=2',
              ],
              '../../benchmarks/dnn/maxpooling')
    bfs = Test('../../samples/bfs',
               'bfs',
               ['-node=1024'],
               '../../benchmarks/shoc/bfs')
    st = Test('../../samples/stencil2d',
              'stencil2d',
              [],
              '../../benchmarks/shoc/stencil2d')
    cw = Test('../../samples/concurrentworkload',
              'concurrentworkload',
              [],
              '')
    ck = Test('../../samples/concurrentkernel',
              'concurrentkernel',
              [],
              '')
    sp = Test('../../samples/spmv',
              'spmv',
              [],
              '../../benchmarks/shoc/spmv')
    fft = Test('../../samples/fft',
               'fft',
               ['-MB=2'],
               '../../benchmarks/shoc/fft')
    nb = Test('../../samples/nbody',
              'nbody',
              [],
              '../../benchmarks/shoc/nbody')


    err = False

    err |= compile('../../rdnainsts/rdnadisassembler')
    err |= atax.test()
    err |= bicg.test()
    #err |= aes.test()
    err |= fir.test()
    err |= km.test()
    # err |= pagerank.test()
    err |= mm.test()
    err |= mt.test()
    # err |= bs.test()
    err |= sc.test()
    #err |= fw.test(test_multi_gpu=False)
    err |= re.test()
    err |= mp.test()
    # err |= bfs.test(test_multi_gpu=False)
    err |= st.test()
    err |= sp.test(test_multi_gpu=False)
    # err |= fft.test(test_multi_gpu=False)
    # err |= nb.test(test_multi_gpu=False)

    # err |= ck.test(test_disassemble=False,
    #                 test_unified_multi_gpu=False,
    #                 test_multi_gpu=False)
    # err |= cw.test(test_disassemble=False,
    #                 test_unified_multi_gpu=False,
    #                 test_multi_gpu=False)

    # error |= compile('acceptancetests/cupipelinedraining')
    # error |= run_test('CU Pipeline Draining',
    #                   './cupipelinedraining -timing',
    #                   'acceptancetests/cupipelinedraining')
    # error |= run_test('CU Pipeline Draining Parallel',
    #                   './cupipelinedraining -timing -parallel',
    #                   'acceptancetests/cupipelinedraining')

    # error |= compile('acceptancetests/tlbshootown')
    # error |= run_test('TLB Shootdown',
    #                   './tlbshootdown -timing',
    #                   'acceptancetestes/tlbshootdown')

    if err:
        sys.exit(1)


if __name__ == '__main__':
    main()
