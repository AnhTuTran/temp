# from __future__ import print_function
# from pymic._misc import _map_data_types as map_data_types

# import pymic as mic
# import numpy as np
# import sys
# import time


# def limiter(data_size):
#     if data_size < 128:
#         return 10000
#     if data_size < 1024:
#         return 1000
#     if data_size < 8192:
#         return 100
#     return 10

# # load the library with the kernel function (on the target)
# device = mic.devices[0]
# library = device.load_library("liboffload_array.so")

# # use the default stream
# stream = device.get_default_stream()

# # sizes of the matrices
# data_size = [(128 + i * 128) for i in range(12)]

# csv = open("benchmarks_exp.csv", "w")
# try:
#     print("size,numpy,pymic kernel,pymic,nrep", file=csv)

#     for i in data_size:
#         m, n = i, i
#         nrep = 10

#         print("Measuring data size ", i)

#         # construct some matrices
#         np.random.seed(10)
#         a = np.random.random(m * n).reshape((m, n))
    
#         # numpy calculation
#         np_mxm_start = time.time()
#         for j in xrange(nrep):
#             c = np.exp(a)
#         np_mxm_end = time.time()
        
#         #timing kernel including memory transfering and computing
#         np_mic_start = time.time()
#         for j in xrange(nrep):
#             offl_a = stream.bind(a)
#             offl_c = mic.exp(offl_a)
#             offl_c.update_host()
#         stream.sync()
#         np_mic_end = time.time()
    
#         #timing kernel only
#         np_mic_kernel_start = time.time()
#         for j in xrange(nrep):
#             offl_c = mic.exp(offl_a)    
#         stream.sync()
#         np_mic_kernel_end = time.time()

#         #offl_c.update_host()
#         #print(np.exp(a))
#         #print(offl_c)

#         # calculate execution time
#         np_mic_time = (np_mic_end - np_mic_start) / nrep
#         np_mic_kernel_time = (np_mic_kernel_end - np_mic_kernel_start) / nrep
#         np_mxm_time = (np_mxm_end - np_mxm_start) / nrep

#         print('{0},{1:.5},{2:.5},{3:.5},{4}'.format(i, np_mxm_time,
#         np_mic_kernel_time, np_mic_time, nrep), file=csv)
# finally:
#     csv.close()


from __future__ import print_function

import sys
import time

import pymic
import numpy as np


def limiter(data_size):
    if data_size < 128:
        return 10000
    if data_size < 1024:
        return 1000
    if data_size < 8192:
        return 100
    return 10

benchmark = "benchmarks_exp"

# number of elements to copyin (8B to 2 GB)
data_sizes = []
data_sizes = [(128 + i * 128) for i in range(12)]
repeats = map(limiter, data_sizes)

device = pymic.devices[0]
library = device.load_library("liboffload_array.so")
stream = device.get_default_stream()

timings = {}
timings_kernel = {}
np.random.seed(10)

try:
    csv = open(benchmark + ".csv", "w")
    print("benchmark;elements;numpy;pymic;pymic kernel;nrep", file=csv)
    for ds, nrep in zip(data_sizes, repeats):
        print("Measuring with data size {0}x{0} (repeating {2})".format(ds, ds, nrep))
        
        m, k, n = ds, ds, ds
        
        a = np.random.random(m * k).reshape((m, k))

        for i in range(nrep):
            np_ts = time.time()
            c = np.exp(a)
            np_te = time.time()

            ts = time.time()
            offl_a = stream.bind(a)
            offl_c = pymic.exp(offl_a)
            offl_c.update_host()
            stream.sync()
            te = time.time()


            ts_kernel = time.time()
            offl_c = pymic.exp(offl_a)
            stream.sync()
            te_kernel = time.time()

            print("{0};{1};{2};{3};{4};{5}".format(benchmark, ds, np_te - np_ts, 
            te - ts, te_kernel - ts_kernel, nrep), file=csv)
finally:
    csv.close()    
 
# try:
#     csv = open(benchmark + ".csv", "w")
#     print("benchmark;elements;avg time;avg time kernel;flops;gflops"
#           ";gflops kernel", 
#           file=csv)
                  
#     for ds in sorted(list(timings)):
#         t, nrep = timings[ds]
#         t = (float(t) / nrep) 
#         t_k, dummy = timings_kernel[ds]
#         t_k = (float(t_k) / nrep)
#         flops = 2 * ds * ds * ds
#         gflops = (float(flops) / (1000 * 1000 * 1000)) / t
#         gflops_k = (float(flops) / (1000 * 1000 * 1000)) / t_k
#         print("{0};{1};{2};{3};{4};{5};{6}".format(benchmark, ds, t, 
#                                                    t_k, flops, gflops, 
#                                                    gflops_k),
#               file=csv)
# finally:
#     csv.close()
