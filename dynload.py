import io
import sys
import requests
import importlib.util
from contextlib import contextmanager


# create file-like string to capture output

@contextmanager
def std_redirect():
    code_out = io.StringIO()
    code_err = io.StringIO()

    # capture output and errors
    sys.stdout = code_out
    sys.stderr = code_err

    # Get console err / out
    _stdout = code_out.getvalue()
    _stderr = code_err.getvalue()

    # Revert stdout/stderr
    sys.stdout = sys.__stdout__
    sys.stderr = sys.__stderr__

    yield _stdout, _stderr

    code_out.close()
    code_err.close()


def run_python_code(code: str, fn_name: str, *params, **kwargs):

    _locals = {}
    _globals = {}
    exec(code, _globals, _locals)

    return _locals[fn_name](*params, **kwargs)


if __name__ == '__main__':
    url = "http://localhost:8000/two.py"
    content = requests.get(url).content.decode("UTF-8")

    run_python_code(content, "hola", 11)

    files = ["demo.py", "two.py"]

    for f in files:
        with open(f, "r") as f:
            # ret, stdout = run_python_code(f.read(), "hola", 11)
            try:
                ret = run_python_code(f.read(), "hola", 11)
                print(ret)

            except ValueError as e:
                print(e)
                continue
