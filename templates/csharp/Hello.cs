using System;

class Hello {
    static string Greet(string name) {
        return $"Hello, {name}!";
    }

    static void Main() {
        Console.WriteLine(Greet("World"));
    }
}
