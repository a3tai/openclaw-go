#!/usr/bin/env python3
"""Build an upstream release delta report for openclaw.

The report compares the target tag with the previous stable tag and surfaces:
- changed upstream files (focused gateway/protocol paths)
- added/removed methods
- methods touched by changed server-method files
"""

from __future__ import annotations

import argparse
import json
import re
import subprocess
from dataclasses import dataclass
from pathlib import Path


TAG_RE = re.compile(r"^v(\d+)\.(\d+)\.(\d+)(?:-(\d+))?$")
METHOD_RE = re.compile(r'"([a-z][a-z0-9_.-]+)"\s*:\s*(?:async\s*)?\(')


@dataclass(frozen=True)
class ReleaseTag:
    name: str
    key: tuple[int, int, int, int]


def run_git(repo: Path, args: list[str]) -> str:
    proc = subprocess.run(
        ["git", "-C", str(repo), *args],
        check=True,
        capture_output=True,
        text=True,
    )
    return proc.stdout


def parse_tag(tag: str) -> ReleaseTag | None:
    match = TAG_RE.match(tag.strip())
    if not match:
        return None
    major, minor, patch, suffix = match.groups()
    return ReleaseTag(
        name=tag.strip(),
        key=(int(major), int(minor), int(patch), int(suffix or 0)),
    )


def list_release_tags(repo: Path) -> list[ReleaseTag]:
    tags_raw = run_git(repo, ["tag", "--list", "v*"]).splitlines()
    parsed = [parsed for parsed in (parse_tag(tag) for tag in tags_raw) if parsed]
    return sorted(parsed, key=lambda item: item.key)


def resolve_previous_tag(repo: Path, version: str) -> str:
    tags = list_release_tags(repo)
    by_name = {tag.name: idx for idx, tag in enumerate(tags)}
    if version not in by_name:
        raise ValueError(f"version tag not found in upstream repo: {version}")
    idx = by_name[version]
    if idx == 0:
        raise ValueError(f"no previous stable tag available before {version}")
    return tags[idx - 1].name


def list_files_at_tag(repo: Path, tag: str, path_prefix: str) -> list[str]:
    out = run_git(repo, ["ls-tree", "-r", "--name-only", tag, "--", path_prefix])
    return [line.strip() for line in out.splitlines() if line.strip()]


def read_file_at_tag(repo: Path, tag: str, file_path: str) -> str:
    return run_git(repo, ["show", f"{tag}:{file_path}"])


def collect_methods_at_tag(repo: Path, tag: str) -> dict[str, str]:
    methods: dict[str, str] = {}
    files = list_files_at_tag(repo, tag, "src/gateway/server-methods")
    for file_path in files:
        name = Path(file_path).name
        if not name.endswith(".ts") or name.endswith(".test.ts") or ".helpers." in name:
            continue
        content = read_file_at_tag(repo, tag, file_path)
        for match in METHOD_RE.finditer(content):
            method = match.group(1)
            methods.setdefault(method, file_path)
    return methods


def changed_files_between(repo: Path, base: str, head: str) -> list[str]:
    out = run_git(repo, ["diff", "--name-only", f"{base}..{head}"])
    return [line.strip() for line in out.splitlines() if line.strip()]


def build_report(upstream_root: Path, version: str, previous_version: str | None) -> dict[str, object]:
    prev = previous_version or resolve_previous_tag(upstream_root, version)
    changed_files = changed_files_between(upstream_root, prev, version)

    changed_method_files = [
        path
        for path in changed_files
        if path.startswith("src/gateway/server-methods/") and path.endswith(".ts")
    ]
    changed_schema_files = [
        path
        for path in changed_files
        if path.startswith("src/gateway/protocol/schema/") and path.endswith(".ts")
    ]

    prev_methods = collect_methods_at_tag(upstream_root, prev)
    curr_methods = collect_methods_at_tag(upstream_root, version)

    prev_set = set(prev_methods)
    curr_set = set(curr_methods)
    added_methods = sorted(curr_set - prev_set)
    removed_methods = sorted(prev_set - curr_set)

    changed_methods = sorted(
        {
            *[m for m, src in prev_methods.items() if src in changed_method_files],
            *[m for m, src in curr_methods.items() if src in changed_method_files],
        }
    )

    return {
        "version": version,
        "previous_version": prev,
        "compare_range": f"{prev}..{version}",
        "compare_url": f"https://github.com/openclaw/openclaw/compare/{prev}...{version}",
        "changed_files_count": len(changed_files),
        "changed_method_files": changed_method_files,
        "changed_schema_files": changed_schema_files,
        "added_methods": added_methods,
        "removed_methods": removed_methods,
        "changed_methods": changed_methods,
    }


def render_markdown(report: dict[str, object]) -> str:
    changed_method_files = report.get("changed_method_files", [])
    changed_schema_files = report.get("changed_schema_files", [])
    added_methods = report.get("added_methods", [])
    removed_methods = report.get("removed_methods", [])
    changed_methods = report.get("changed_methods", [])

    lines = [
        "## Upstream Release Diff Evidence",
        "",
        f"- Compare range: `{report['compare_range']}`",
        f"- Upstream compare: {report['compare_url']}",
        f"- Changed files: **{report['changed_files_count']}**",
        "",
        "### RPC method deltas",
        f"- Added methods: **{len(added_methods)}**",
        *(f"  - `{m}`" for m in added_methods),
        f"- Removed methods: **{len(removed_methods)}**",
        *(f"  - `{m}`" for m in removed_methods),
        f"- Methods touched in changed server-method files: **{len(changed_methods)}**",
        *(f"  - `{m}`" for m in changed_methods),
        "",
        "### Changed upstream files (gateway/protocol focus)",
        f"- `src/gateway/server-methods/*.ts`: **{len(changed_method_files)}**",
        *(f"  - `{path}`" for path in changed_method_files),
        f"- `src/gateway/protocol/schema/*.ts`: **{len(changed_schema_files)}**",
        *(f"  - `{path}`" for path in changed_schema_files),
        "",
    ]
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--upstream-root", required=True)
    parser.add_argument("--version", required=True)
    parser.add_argument("--previous-version")
    parser.add_argument("--output")
    parser.add_argument("--markdown-output")
    args = parser.parse_args()

    report = build_report(Path(args.upstream_root), args.version, args.previous_version)

    encoded = json.dumps(report, indent=2)
    if args.output:
        Path(args.output).write_text(encoded + "\n", encoding="utf-8")
    else:
        print(encoded)

    if args.markdown_output:
        Path(args.markdown_output).write_text(render_markdown(report), encoding="utf-8")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
