def hola(sum: int, **kwargs) -> int:

    sum2 = sum + 1

    print("XXX")

    raise ValueError("ooooo")

    return sum2


if __name__ == '__main__':
    hola(2)