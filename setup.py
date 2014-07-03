from distutils.core import setup
from setuptools import find_packages

setup(
    name="cli-aoke",
    version="1.0.0",
    author="Jessica Frazelle",
    author_email="jfrazelle@me.com",
    url="https://github.com/jfrazelle/cli-aoke",
    packages=find_packages(),
    description="Command Line Karaoke",
    license="MIT",
    classifiers=(
        'Development Status :: 5 - Production/Stable',
        'Intended Audience :: Developers',
        'Programming Language :: Python',
        'Operating System :: OS Independent',
        'Topic :: Internet :: WWW/HTTP',
        'Topic :: Software Development :: Libraries :: Python Modules'
    ),
    install_requires=[
        "BeautifulSoup4==4.3.2"
    ],
    keywords=['karaoke', 'command', 'line', 'magic'],
    scripts=['bin/cli-aoke']
)
