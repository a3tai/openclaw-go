#!/usr/bin/env python3
"""Validate openclaw-go coverage for a specific upstream release.

Checks:
- no missing upstream methods in go gateway wrappers
- release-delta methods are covered by gateway method tests
"""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path

from upstream_drift_report import collect_go_methods, collect_upstream_methods
from upstream_release_diff import build_report


TEST_METHOD_RE = re.compile(r'method:\s*"([a-z][a-z0-9_.-]+)"')


def collect_tested_methods(go_root: Path) -> set[str]:
    test_file = go_root / "gateway" / "methods_test.go"
    content = test_file.read_text(encoding="utf-8")
    return set(TEST_METHOD_RE.findall(content))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--upstream-root", required=True)
    parser.add_argument("--go-root", required=True)
    parser.add_argument("--version", required=True)
    parser.add_argument("--previous-version")
    parser.add_argument("--output")
    args = parser.parse_args()

    upstream_root = Path(args.upstream_root)
    go_root = Path(args.go_root)

    release_report = build_report(upstream_root, args.version, args.previous_version)

    upstream_methods, _, _, _ = collect_upstream_methods(upstream_root)
    go_methods, _ = collect_go_methods(go_root)

    missing = sorted(upstream_methods - go_methods)

    tested_methods = collect_tested_methods(go_root)
    delta_methods = set(release_report.get("added_methods", [])) | set(
        release_report.get("changed_methods", [])
    )
    delta_methods &= upstream_methods
    untested_delta_methods = sorted(delta_methods - tested_methods)

    result = {
        "version": args.version,
        "previous_version": release_report.get("previous_version"),
        "compare_url": release_report.get("compare_url"),
        "missing_methods": missing,
        "delta_methods": sorted(delta_methods),
        "untested_delta_methods": untested_delta_methods,
        "changed_schema_files": release_report.get("changed_schema_files", []),
        "ok": len(missing) == 0 and len(untested_delta_methods) == 0,
    }

    encoded = json.dumps(result, indent=2)
    if args.output:
        Path(args.output).write_text(encoded + "\n", encoding="utf-8")
    else:
        print(encoded)

    if missing:
        print("Validation failed: missing upstream methods in gateway wrappers.")
        return 2
    if untested_delta_methods:
        print("Validation failed: release-delta methods missing gateway/methods_test.go coverage.")
        return 3
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
