#include<stdio.h>
#include<stdlib.h>

int main(){
    int size = 5;
    char * s = malloc(sizeof(char)*size);
    int position = 0;

    while(1){
        char c;
        c = getchar();
        s[position] = c;
        position ++;
        if(position == size){
            printf("%c\n",*s);
            size+=5;
            s = realloc(s,size*sizeof(char));

        }
        if(position>=10){
            break;
        }
    }


    return 0;
}