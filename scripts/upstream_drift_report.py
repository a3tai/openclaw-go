#!/usr/bin/env python3
"""Compare upstream OpenClaw gateway methods to openclaw-go gateway methods.

Outputs JSON with missing/extra methods and exits non-zero when missing methods
are detected unless --allow-missing is set.
"""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path


UPSTREAM_METHOD_RE = re.compile(r'"([a-z][a-z0-9_.-]+)"\s*:\s*(?:async\s*)?\(')
GO_LITERAL_METHOD_RE = re.compile(r'sendRPC(?:Typed|Void)?\(ctx,\s*"([a-z][a-z0-9_.-]+)"')
GO_CONST_METHOD_RE = re.compile(r'sendRPC(?:Typed|Void)?\(ctx,\s*string\(protocol\.(Method[A-Za-z0-9]+)\)')
PROTOCOL_METHOD_CONST_RE = re.compile(r'\b(Method[A-Za-z0-9]+)\s+MethodName\s*=\s*"([a-z][a-z0-9_.-]+)"')


def collect_upstream_methods(upstream_root: Path) -> tuple[set[str], dict[str, str]]:
    methods: set[str] = set()
    sources: dict[str, str] = {}
    methods_dir = upstream_root / "src" / "gateway" / "server-methods"
    for path in methods_dir.glob("*.ts"):
        name = path.name
        if ".test." in name or name.endswith(".test.ts") or name.endswith(".helpers.ts"):
            continue
        content = path.read_text(encoding="utf-8")
        for match in UPSTREAM_METHOD_RE.finditer(content):
            method = match.group(1)
            methods.add(method)
            sources.setdefault(method, str(path.relative_to(upstream_root)))
    return methods, sources


def collect_protocol_method_constants(go_root: Path) -> dict[str, str]:
    protocol_file = go_root / "protocol" / "protocol.go"
    content = protocol_file.read_text(encoding="utf-8")
    return {name: value for name, value in PROTOCOL_METHOD_CONST_RE.findall(content)}


def collect_go_methods(go_root: Path) -> tuple[set[str], dict[str, str]]:
    methods: set[str] = set()
    sources: dict[str, str] = {}
    method_constants = collect_protocol_method_constants(go_root)
    gateway_dir = go_root / "gateway"
    for path in gateway_dir.glob("*.go"):
        if path.name.endswith("_test.go"):
            continue
        content = path.read_text(encoding="utf-8")
        for match in GO_LITERAL_METHOD_RE.finditer(content):
            method = match.group(1)
            methods.add(method)
            sources.setdefault(method, str(path.relative_to(go_root)))
        for match in GO_CONST_METHOD_RE.finditer(content):
            const_name = match.group(1)
            method = method_constants.get(const_name)
            if method:
                methods.add(method)
                sources.setdefault(method, str(path.relative_to(go_root)))
    return methods, sources


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--upstream-root", required=True)
    parser.add_argument("--go-root", required=True)
    parser.add_argument("--output", required=False)
    parser.add_argument("--allow-missing", action="store_true")
    args = parser.parse_args()

    upstream_root = Path(args.upstream_root)
    go_root = Path(args.go_root)

    upstream_methods, upstream_sources = collect_upstream_methods(upstream_root)
    go_methods, go_sources = collect_go_methods(go_root)

    missing = sorted(upstream_methods - go_methods)
    extra = sorted(go_methods - upstream_methods)

    report = {
        "upstream_method_count": len(upstream_methods),
        "go_method_count": len(go_methods),
        "missing_in_go": missing,
        "extra_in_go": extra,
        "missing_details": [
            {"method": m, "upstream_source": upstream_sources.get(m, "")}
            for m in missing
        ],
        "extra_details": [
            {"method": m, "go_source": go_sources.get(m, "")}
            for m in extra
        ],
    }

    encoded = json.dumps(report, indent=2)
    if args.output:
        Path(args.output).write_text(encoded + "\n", encoding="utf-8")
    else:
        print(encoded)

    if missing and not args.allow_missing:
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
