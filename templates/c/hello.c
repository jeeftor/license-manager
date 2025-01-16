#include <stdio.h>
#include "hello.h"

void greet(const char* name) {
    printf("Hello, %s!\n", name);
}

int main() {
    greet("World");
    return 0;
}
