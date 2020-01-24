from os import path
from os.path import join
from setuptools import setup, find_packages

here = path.abspath(path.dirname(__file__))

with open(join(here, 'requirements.txt')) as f:
    required = f.read().splitlines()

setup(
    name='dockerfile-sec',
    version="1.0.1",
    packages=find_packages(),
    description='Find security issues in Dockerfiles',
    long_description="Find security issues in Dockerfiles",
    install_requires=required,
    include_package_data=True,
    zip_safe=True,
    url='https://github.com/cr0hn/dockerfile-sec',
    license='Apache 2.0',
    author='cr0hn',
    entry_points={'console_scripts': [
        'dockerfile-sec = dockerfile_sec.__main__:main'
    ]},
    classifiers=[
        'Environment :: Console',
        'Intended Audience :: System Administrators',
        'Intended Audience :: Other Audience',
        'License :: OSI Approved :: BSD License',
        'Operating System :: MacOS',
        'Operating System :: Microsoft :: Windows',
        'Operating System :: POSIX',
        'Programming Language :: Python :: 3',
    ],
)

