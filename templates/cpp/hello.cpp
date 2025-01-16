#include <iostream>
#include "hello.hpp"

namespace hello {
    std::string greet(const std::string& name) {
        return "Hello, " + name + "!";
    }
}

int main() {
    std::cout << hello::greet("World") << std::endl;
    return 0;
}
