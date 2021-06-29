#include <mpi.h>
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>

#define DATA_SIZE (4*1024*1024)

int (*func_tune_alltoall)(const int);
void *handle;
void GetInfo(char* cname)
{
    MPI_Datatype datatype;
    MPI_T_enum enumtype;
    int err = MPI_SUCCESS;
    int cvar_index = 0;
    int name_len = 20;
    int verbosity = 1;
    int desc_len = 100;
    int binding;
    int scope;
    char desc[desc_len];
    char name[name_len];    
    int myrank=0;
    
    MPI_Comm_rank(MPI_COMM_WORLD, &myrank);
    if(myrank ==0) 
        printf("\n%s Rank=%d\n", cname, myrank);
    
    err = MPI_T_cvar_get_index( cname , &cvar_index);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr,"\nERROR IN get index (%d)\n",err);
    }

    err = MPI_T_cvar_get_info(cvar_index, name, &name_len, &verbosity,
                        &datatype, &enumtype, desc, &desc_len,
                        &binding, &scope);
    if(myrank ==0)
    {
        printf("name by index = %s\n",name);
        printf("verbosity=%d\n",verbosity);
        printf("datatype=%x\n",(int)datatype);
        printf("desc=%s\n",desc);
        printf("binding=%d\n",binding);
        printf("scope=%d\n",scope);
    }
}



int test_cvar(char* cname, int value, void* object){
    MPI_T_cvar_handle handle;
    int err = MPI_SUCCESS;
    int read_value, count, cvar_index = 0;
    
    err = MPI_T_cvar_get_index( cname , &cvar_index);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr,"\nERROR IN get index (%d)\n",err);
    }

    err = MPI_T_cvar_handle_alloc(cvar_index, object, &handle, &count);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr, "\nERROR cvar_handle_alloc fpr NO_OBJ (%d)\n", err);
    }
#if 0
    if(object != NULL)
        printf("MPI_T_cvar_handle_alloc : Bind object=%x, count=%d\n",*(int*)object, count);
#endif
    
    err = MPI_T_cvar_read(handle, &read_value);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr, "\nCVAR READ error %d\n", err);
    }
    fprintf(stderr, "\tPrevious value of   %s = %10d\n", cname,  read_value );

    err = MPI_T_cvar_write(handle, &value);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr,"\nERROR IN mpit write (%d)\n",err);
    }
    err = MPI_T_cvar_read(handle, &read_value);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr, "\nCVAR READ error %d\n", err);
    }
    //fprintf(stderr, "\tRequested value of  %s = %10d\n\tConfigured_Value of %s = %10d\n", cname, value, cname, read_value );

    MPI_T_cvar_handle_free(&handle);    
    return 0;
}


int read_cvar(char* cname, void* object) {
    MPI_T_cvar_handle handle;
    int err = MPI_SUCCESS;
    int read_value, count, cvar_index = 0;
    
    err = MPI_T_cvar_get_index( cname , &cvar_index);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr,"\nERROR IN get index (%d)\n",err);
    }

    err = MPI_T_cvar_handle_alloc(cvar_index, object, &handle, &count);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr, "\nERROR cvar_handle_alloc fpr NO_OBJ (%d)\n", err);
    }
    err = MPI_T_cvar_read(handle, &read_value);
    if(err!=MPI_SUCCESS) {
        fprintf(stderr, "\nCVAR READ error %d\n", err);
    }
    return read_value;
}


int MPI_Alltoall(const void *sendbuf, int sendcount, MPI_Datatype sendtype,
                 void *recvbuf, int recvcount, MPI_Datatype recvtype,
                                  MPI_Comm comm)
{
    int err, ithreadsupport, size;
    err=MPI_T_init_thread(MPI_THREAD_SINGLE,&ithreadsupport);    
    if(err!=MPI_SUCCESS)
        fprintf(stderr,"ERROR IN MPI_Alltoall (%d)",err);
    err = MPI_Type_size(sendtype, &size);
    if(err!=MPI_SUCCESS)
        fprintf(stderr,"ERROR in MPI_Alltoall (%d)",err);
    int tune_val = (*func_tune_alltoall)(size*sendcount);
    MPI_Comm world_comm = MPI_COMM_WORLD;
    test_cvar("MPIR_CVAR_ALLTOALL_TUNE", tune_val, (void*)&world_comm);
    return PMPI_Alltoall(sendbuf, sendcount, sendtype, 
            recvbuf, recvcount, recvtype, comm);
    err=MPI_T_finalize();
    if(err!=MPI_SUCCESS)
        fprintf(stderr,"ERROR in MPI_T finalize (%d)",err);
}
int MPI_Init(int *argc, char ***argv)
{
   
    handle = dlopen("./libmvapich_tuner.so", RTLD_LAZY);    
    if (!handle) {
        /* fail to load the library */
        fprintf(stderr, "Error: %s\n", dlerror());
        return -1;
    }

    func_tune_alltoall = (int (*)(int))dlsym(handle, "tune_alltoall");
    if (!func_tune_alltoall) {
        /* no such symbol */
        fprintf(stderr, "Error: %s\n", dlerror());
        dlclose(handle);
        return -1;
    }
    return  PMPI_Init(argc, argv);

    

}
int MPI_Finalize(void)
{
    dlclose(handle);
    return PMPI_Finalize();

}
int main(int argc, char **argv)
{
    int err, ithreadsupport, val = 0;
    err=MPI_T_init_thread(MPI_THREAD_SINGLE,&ithreadsupport);    
    if(err!=MPI_SUCCESS)
        fprintf(stderr,"ERROR IN MPI T INIT THREAD (%d)",err);

    //-----------------------------------------------------------------   
    
    val = 2;
    test_cvar("MPIR_CVAR_ALLTOALL_TUNE", val, NULL);

    
    int world_rank, world_size;
    MPI_Comm world_comm = MPI_COMM_WORLD;
    MPI_Init(&argc, &argv);

    MPI_Comm_rank(MPI_COMM_WORLD, &world_rank);
    MPI_Comm_size(MPI_COMM_WORLD, &world_size);

    if(world_rank == 0) {
        fprintf(stdout, "==================================================\n");
        fprintf(stdout, "Launched %d processes\n", world_size);
        fprintf(stdout, "==================================================\n");
        fflush(stdout);
    }
    val = 3;
    test_cvar("MPIR_CVAR_ALLTOALL_TUNE", val, (void*)&world_comm);
    size_t size = 4;
    char * sendbuf = malloc(size*world_size);
    char * recvbuf = malloc(size*world_size);
    MPI_Alltoall(sendbuf, size, MPI_CHAR, recvbuf, size, MPI_CHAR,
                         MPI_COMM_WORLD);
    val = 2;
    test_cvar("MPIR_CVAR_ALLTOALL_TUNE", val, (void*)&world_comm);
    MPI_Alltoall(sendbuf, size, MPI_CHAR, recvbuf, size, MPI_CHAR,
                         MPI_COMM_WORLD);
    err=MPI_T_finalize();
    if(err!=MPI_SUCCESS)
    fprintf(stderr,"ERROR in MPI_T finalize (%d)",err);
    MPI_Finalize();

    return 0;
}
