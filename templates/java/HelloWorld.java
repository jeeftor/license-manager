public class HelloWorld {
    public static String greet(String name) {
        return String.format("Hello, %s!", name);
    }

    public static void main(String[] args) {
        System.out.println(greet("World"));
    }
}
