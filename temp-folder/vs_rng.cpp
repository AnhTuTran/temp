#include <stdio.h> 
#include <mkl.h>
#include <iostream>
#include <time.h>
#include <limits.h>
#include <omp.h>

using namespace std;


void discrete_uniform_dis() {
    long N = 10;
    int a = 1;
    int b = 100;
    int cores = 240;        
    srand(time(0));
    
    VSLStreamStatePtr *stream = new VSLStreamStatePtr[cores];
    for (int i = 0; i < cores; i++)
        vslNewStream(&stream[i], VSL_BRNG_MT2203 + i, time(0));

    printf("%10s %13.7s %13.7s\n", "N", "mkl_T", "rand_T");    
    
    for (N = cores; N < INT_MAX*1L; N *= 10) {
        int *r = new int[N];
        long size = N / cores;

        double mkl_time = omp_get_wtime();
        #pragma omp parallel for
        for (int j = 0; j < cores; j++)
            viRngUniform(VSL_RNG_METHOD_UNIFORM_STD, stream[j], size, r + j*size, a, b);
        mkl_time = omp_get_wtime() - mkl_time;
            
        double rand_time = omp_get_wtime();
        for (long i = 0; i < N; i++)
            r[i] = rand();
        rand_time = omp_get_wtime() - rand_time;
        
        delete [] r;
        printf("%10ld %13.7f %13.7f\n", N, mkl_time, rand_time);
    }
    for (int j = 0; j < cores; j++)
        vslDeleteStream(&stream[j]);
}



int main() {
    discrete_uniform_dis();
    return 0; 
}


