#remove all lines that start with a print statement even if they are commented out

import re

def starts_with_print(line):
    res = re.match("\s*#?print", line)
    if res != None:
        return True
    return False

def remove_print(filename):
    with open(filename, "r") as f:
        lines = f.readlines()
    with open(filename, "w") as f:
        for line in lines:
            if not starts_with_print(line):
                f.write(line)


import sys

if __name__ == "__main__":
    for file in sys.argv:
        remove_print(file)