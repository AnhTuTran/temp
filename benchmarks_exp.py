from __future__ import print_function
from pymic._misc import _map_data_types as map_data_types

import pymic as mic
import numpy as np
import sys
import time


# load the library with the kernel function (on the target)
device = mic.devices[0]
library = device.load_library("liboffload_array.so")

# use the default stream
stream = device.get_default_stream()

# sizes of the matrices
data_size = [(128 + i * 128) for i in range(12)]

csv = open("benchmarks_exp.csv", "w")
try:
    print("size,numpy,pymic kernel,pymic,nrep", file=csv)

    for i in data_size:
        m, n = i, i
        nrep = 10

        print("Measuring data size ", i)

        # construct some matrices
        np.random.seed(10)
        a = np.random.random(m * n).reshape((m, n))
    
        # numpy calculation
        np_mxm_start = time.time()
        for j in xrange(nrep):
            c = np.exp(a)
        np_mxm_end = time.time()
        
        #timing kernel including memory transfering and computing
        np_mic_start = time.time()
        for j in xrange(nrep):
            offl_a = stream.bind(a)
            offl_c = mic.exp(offl_a)
            offl_c.update_host()
        stream.sync()
        np_mic_end = time.time()
    
        #timing kernel only
        np_mic_kernel_start = time.time()
        for j in xrange(nrep):
            offl_c = mic.exp(offl_a)    
        stream.sync()
        np_mic_kernel_end = time.time()

        #offl_c.update_host()
        #print(np.exp(a))
        #print(offl_c)

        # calculate execution time
        np_mic_time = (np_mic_end - np_mic_start) / nrep
        np_mic_kernel_time = (np_mic_kernel_end - np_mic_kernel_start) / nrep
        np_mxm_time = (np_mxm_end - np_mxm_start) / nrep

        print('{0},{1:.5},{2:.5},{3:.5},{4}'.format(i, np_mxm_time,
        np_mic_kernel_time, np_mic_time, nrep), file=csv)
finally:
    csv.close()
