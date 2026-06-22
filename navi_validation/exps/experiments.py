""" Generate experiments """

import re
import exps.empty_kernel.exp as empty_kernel
import exps.flat_load.exp as flat_load
import exps.pointer_chasing.exp as pc
# import exps.v_add.exp as v_add
# import exps.branch.exp as branch
# import exps.fir.exp as fir
# import exps.two_cache_line.exp as two_cacheline
# import exps.two_flat_load.exp as two_flat_load
# import exps.dram_sequential.exp as dram_sequential
# import exps.cache_size.exp as cache_size
# import exps.cache_size_zig_zag.exp as cache_size_zig_zag
# import exps.mem_read.exp as mem_read
# import exps.mem_write.exp as mem_write
# import exps.mem_copy.exp as mem_copy


class Experiments(object):
    """ Experiments generate experiments to run """

    def __init__(self):
        self.experiments = []

    def get_experiments(self, filter):
        """ returns a list of experiments to run """
        self.__list_experiments()
        self.__filter_experiments(filter)
        return self.experiments

    def __filter_experiments(self, filter):
        if filter == "":
            return

        new_exp = []
        filter_re = re.compile(filter)
        for exp in self.experiments:
            m = filter_re.match(exp.name())

            if m != None:
                new_exp.append(exp)

        self.experiments = new_exp

    def __list_experiments(self):
        exps = []
        exps.append(empty_kernel.EmptyKernelExp())
        exps.append(pc.PCExp())
        # exps.append(v_add.VAddExp())
        # exps.append(branch.BranchExp())
        exps.append(flat_load.FlatLoadExp())
        # exps.append(two_flat_load.FlatLoadTwoLoadExp())
        # exps.append(two_cacheline.FlatLoadTwoCachelineExp())
        # exps.append(cache_size.Exp())
        # exps.append(cache_size_zig_zag.Exp())
        # exps.append(dram_sequential.DramSequentialExp())
        # exps.append(mem_read.Exp())
        # exps.append(mem_write.Exp())
        # exps.append(mem_copy.Exp())
        # exps.append(fir.FIRExp())
        self.experiments = exps
