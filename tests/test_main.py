"""Tests for main.py CLI"""

import pytest

from main import parse_args, build_agent, main


class TestParseArgs:
    def test_defaults(self):
        args = parse_args([])
        assert args.backend == "mock"
        assert args.max_steps == 10
        assert not args.verbose
        assert args.query is None

    def test_backend_openai(self):
        args = parse_args(["--backend", "openai"])
        assert args.backend == "openai"

    def test_query_flag(self):
        args = parse_args(["--query", "hello"])
        assert args.query == "hello"

    def test_max_steps(self):
        args = parse_args(["--max-steps", "5"])
        assert args.max_steps == 5

    def test_verbose(self):
        args = parse_args(["--verbose"])
        assert args.verbose


class TestBuildAgent:
    def test_builds_mock_agent(self):
        args = parse_args([])
        agent = build_agent(args)
        from agent import Agent
        assert isinstance(agent, Agent)


class TestMain:
    def test_single_query_mode(self, capsys):
        exit_code = main(["--query", "What is 2 + 2?"])
        assert exit_code == 0
        captured = capsys.readouterr()
        assert len(captured.out) > 0

    def test_single_query_verbose(self, capsys):
        exit_code = main(["--query", "current time", "--verbose"])
        assert exit_code == 0
