#include <stdio.h>

void loop(){
    char * line;
    char * * args;
    int status = 1 ;
    do{
        printf("> ");
        line = read_line();
        flag = 0;
        args = split_lines(line);
        status = dash_launch(args);

        free(line);
        free(args);
    }while(status);
}

char * read_line(){
    int buffsize = 1024; // 缓冲区的大小
    int positon = 0;
    char * buffer = malloc(sizeof(char) * buffsize); // 申请缓冲区内存
    int c;
    
    if(!buffer){
        fprintf(stderr,"%sdash: Allocation error%s\n",RED, RESET);
        exit(EXIT_FAILTURE);
    }

    while(1){
        c = getchar();
        if(c == EOF || c == '\n'){
            buffer[positon] = '\0'; // 表示说到输出的最后了
            return buffer;
        }else{
            buffer[positon] = c;
        }
        positon++;

        if(positon >= buffsize){
            buffsize += 1024;
            buffer = realloc(buffer,buffsize);

            if(!buffer){
                fprintf(stderr,"dash: Allocation error\n");
                exit(EXIT_FAILTURE);
            }
        }
    }
}

int main(){

    return 0;
}