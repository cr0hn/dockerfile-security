import os
import re
import sys
import json
import select
import argparse
import os.path as op
import urllib.request

from typing import List

import yaml

from terminaltables import AsciiTable

HERE = os.path.dirname(__file__)


def download(remote_file: str) -> str:
    with urllib.request.urlopen(remote_file) as f:
        v = f.read()

        try:
            return v.decode('utf-8')
        except:
            return v


def _process_results(args: argparse.Namespace, found_issues: List[dict]):
    # -------------------------------------------------------------------------
    # Building results
    # -------------------------------------------------------------------------
    table_data = [
        ['Rule Id', 'Description', 'Severity', 'Reference']
    ]

    if not found_issues:
        table_data.append(["No issues found"])
    else:

        for res in found_issues:
            table_data.append((
                res["id"],
                res["description"],
                res["severity"],
                res["reference"]
            ))

        # Export results
        if args.output_file:
            with open(args.output_file, "w") as f:
                json.dump(found_issues, f)

    if not args.quiet:
        if sys.stdout.isatty():
            # We're in terminal
            print(AsciiTable(table_data).table)
        else:
            try:
                sys.stdout.write(json.dumps(found_issues))
                sys.stdout.flush()
            except (BrokenPipeError, IOError) as e:
                # Piped command doesn't support data input as pipe
                sys.stderr.write(e)
            except Exception as e:
                pass


def _load_rules(args: argparse.Namespace) -> List[dict]:
    def __load_all_rules__() -> list:
        _all = ("core", "credentials")
        return [
            op.join(HERE, "rules", f"{a}.yaml")
            for a in _all
        ]

    built_in_rules = args.internal_rules

    rules_files = []

    if built_in_rules == "none":
        pass

    elif built_in_rules == "all":
        rules_files.extend(__load_all_rules__())

    elif built_in_rules == "credentials":
        rules_files.append(op.join(HERE, "rules", "credentials.yaml"))

    elif built_in_rules == "core":
        rules_files.append(op.join(HERE, "rules", "core.yaml"))

    else:
        rules_files.extend(__load_all_rules__())

    rules = []

    for r in rules_files:
        with open(r, "r") as f:
            rules.extend(yaml.safe_load(f.read()))

    if args.rules_file:
        for rule_file in args.rules_file:
            if rule_file.startswith("http"):
                # Load from remote URL
                rules.extend(
                    yaml.safe_load(download(rule_file))
                )
            else:
                # Load from local file
                real_file_path = os.path.join(os.getcwd(), rule_file)
                with open(real_file_path, "r") as f:
                    rules.extend(yaml.safe_load(f.read()))

    return rules


def _load_ignore_ids(args: argparse.Namespace) -> List[str]:
    ignores = []

    if args.ignore_rule:
        for x in args.ignore_rule:
            ignores.extend(x.split(","))

    if args.ignore_file:
        for rule_file in args.ignore_file:
            if rule_file.startswith("http"):
                # Load from remote URL
                ignores.extend(download(rule_file).splitlines())
            else:
                # Load from local file
                real_file_path = os.path.join(os.getcwd(), rule_file)
                with open(real_file_path, "r") as f:
                    ignores.extend([x.replace("\n", "") for x in f.readlines()])

    return ignores


def analyze(args: argparse.Namespace):
    # -------------------------------------------------------------------------
    # Load Dockerfile by stdin or parameter
    # -------------------------------------------------------------------------
    dockerfile_content = None
    if not args.DOCKERFILE:
        if sys.stdin in select.select([sys.stdin], [], [], 0)[0]:
            dockerfile_content = sys.stdin.read()
    else:
        with open(args.DOCKERFILE[0], "r") as f:
            dockerfile_content = f.read()

    if not dockerfile_content:
        raise FileNotFoundError("Dockerfile is needed")

    rules = _load_rules(args)
    ignores = set(_load_ignore_ids(args))

    found_issues = []

    # Matching
    for rule in rules:

        if rule["id"] in ignores:
            continue

        regex = re.search(rule["regex"], dockerfile_content)

        if regex:
            res = rule.copy()
            del res["regex"]

            found_issues.append(res)

    _process_results(args, found_issues)


def main():
    parser = argparse.ArgumentParser(
        description='Analyze a Dockerfile for security'
    )
    parser.add_argument('DOCKERFILE', help="dockerfile path", nargs="*")
    parser.add_argument('-F', '--ignore-file',
                        action="append",
                        help="ignore files")
    parser.add_argument('-i', '--ignore-rule',
                        action="append",
                        help="ignore rule")
    parser.add_argument('-r', '--rules-file',
                        action="append",
                        help="rules file. One rule ID per line")
    parser.add_argument('-R', '--internal-rules',
                        choices=("core", "credentials", "all", "none"),
                        help="use built-in rules. Default: all")
    parser.add_argument('-o', '--output-file',
                        help="output file path")
    parser.add_argument('-q', '--quiet',
                        action="store_true",
                        default=False,
                        help="quiet mode")
    parsed_cli = parser.parse_args()

    try:
        analyze(parsed_cli)
    except Exception as e:
        print("[!] ", e)


if __name__ == '__main__':
    main()
