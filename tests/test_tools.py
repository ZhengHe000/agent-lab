"""Tests for agent/tools.py"""

import os
import tempfile
import pytest

from agent.tools import (
    Tool,
    DEFAULT_TOOLS,
    _calculator,
    _current_time,
    _read_file,
    _write_file,
    _web_search,
)


class TestTool:
    def test_run_returns_string(self):
        t = Tool(name="echo", description="echo", func=lambda x: x.upper())
        assert t.run("hello") == "HELLO"

    def test_run_catches_exceptions(self):
        def boom(_):
            raise ValueError("oops")

        t = Tool(name="boom", description="", func=boom)
        result = t.run("x")
        assert "Error" in result
        assert "oops" in result


class TestCalculator:
    def test_addition(self):
        assert _calculator("2 + 2") == "4"

    def test_power(self):
        assert _calculator("2 ** 10") == "1024"

    def test_sqrt(self):
        assert _calculator("sqrt(144)") == "12.0"

    def test_abs(self):
        assert _calculator("abs(-5)") == "5"

    def test_invalid_expression(self):
        result = _calculator("import os")
        assert "Invalid expression" in result or "Error" in result or result == "None"

    def test_division(self):
        result = float(_calculator("10 / 4"))
        assert result == pytest.approx(2.5)


class TestCurrentTime:
    def test_returns_string(self):
        result = _current_time("")
        assert isinstance(result, str)
        # Should look like YYYY-MM-DD HH:MM:SS
        assert len(result) == 19
        assert result[4] == "-"


class TestReadFile:
    def test_reads_existing_file(self, tmp_path):
        p = tmp_path / "hello.txt"
        p.write_text("hello world", encoding="utf-8")
        assert _read_file(str(p)) == "hello world"

    def test_missing_file(self):
        result = _read_file("/nonexistent/path/file.txt")
        assert "not found" in result.lower() or "File not found" in result

    def test_strips_whitespace_from_path(self, tmp_path):
        p = tmp_path / "test.txt"
        p.write_text("data", encoding="utf-8")
        assert _read_file(f"  {p}  ") == "data"


class TestWriteFile:
    def test_writes_file(self, tmp_path):
        p = tmp_path / "out.txt"
        result = _write_file(f"{p}|hello agent")
        assert "Written" in result
        assert p.read_text(encoding="utf-8") == "hello agent"

    def test_missing_separator(self):
        result = _write_file("/some/path")
        assert "Usage" in result


class TestWebSearch:
    def test_returns_stub_message(self):
        result = _web_search("Python AI agents")
        assert "Python AI agents" in result


class TestDefaultTools:
    def test_all_tools_present(self):
        names = {t.name for t in DEFAULT_TOOLS}
        assert {"calculator", "current_time", "read_file", "write_file", "web_search"} <= names

    def test_all_tools_have_descriptions(self):
        for t in DEFAULT_TOOLS:
            assert t.description, f"Tool '{t.name}' has no description"
