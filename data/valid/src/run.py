import json
import pathlib

P = pathlib.Path(__file__).parent.parent / "in" / "inputs.json"

with open(P, "r") as f:
    inputs = json.load(f)

print(inputs)

def fibonacci(n):
    if n <= 0:
        return 0
    elif n == 1:
        return 1
    else:
        return fibonacci(n - 1) + fibonacci(n - 2)

print(fibonacci(40))