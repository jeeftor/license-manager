#!/usr/bin/env python3






def greet(name: str) -> str:
    """Generate a personalized greeting message.

    Args:
        name (str): The name of the person to greet.

    Returns:
        str: A formatted greeting string.
    """
    return f"Hello, {name}!"


def taco_fn() -> None:
    """This is a function about tacos.

    Performs operations related to tacos (currently empty).

    Args:
        None

    Returns:
        None
    """


if __name__ == "__main__":
    print(greet("World"))
