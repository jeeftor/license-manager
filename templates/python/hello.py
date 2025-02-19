"""
A simple greeting module that provides a personalized greeting message.
"""

def greet(name: str) -> str:
    """
    Creates a personalized greeting message.

    Args:
        name (str): The name to include in the greeting

    Returns:
        str: A formatted greeting string including the provided name
    """
    return f"Hello, {name}!"

if __name__ == "__main__":
    """
    Main execution block that demonstrates the greet function.
    Prints a greeting for 'World' when the script is run directly.
    """
    print(greet("World"))
