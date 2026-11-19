"""Shared helpers for mutate.py (checks.py keeps its own copies from the first run, unchanged)."""
import re
from client import STYLE, text

HANGUL = re.compile(r"[ㄱ-ㆎ가-힣]")
USER_NAME_KO = re.compile(r"긴한글이름+")
LANG_LABEL = re.compile(r'<(span|a)[^>]*lang="ko"[^>]*>한국어</\1>')
results = []


def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": "" if ok else str(detail)[:400]})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:300])


def lang_of(body):
    m = re.search(r'<html lang="(\w*)"', body)
    return m.group(1) if m else None


def hangul_residue(body):
    b = LANG_LABEL.sub("", USER_NAME_KO.sub("", STYLE.sub("", body)))
    return sorted(set(re.findall(r"[ㄱ-ㆎ가-힣]+", b)))[:8]


def kst(body): return re.findall(r"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} KST", body)
def main_text(m): return text(m)
