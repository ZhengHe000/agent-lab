#!/usr/bin/env python3
"""Interactive command-line interface for agent-lab.

Usage
-----
    # Use the built-in mock LLM (no API key required):
    python main.py

    # Use the OpenAI backend:
    OPENAI_API_KEY=sk-... python main.py --backend openai --model gpt-4o-mini

    # Run a single non-interactive query:
    python main.py --query "What is 2 ** 32?"
"""

from __future__ import annotations

import argparse
import sys

from agent import Agent
from agent.llm import create_llm


BANNER = r"""
 ___                    _     _         _
/ _ \                  | |   | |       | |
/ /_\ \ __ _  ___ _ __ | |_  | |     __ _| |__
|  _  |/ _` |/ _ \ '_ \| __| | |    / _` | '_ \
| | | | (_| |  __/ | | | |_  | |___| (_| | |_) |
\_| |_/\__, |\___|_| |_|\__| \_____/\__,_|_.__/
        __/ |
       |___/    智能体 (AI Agent) — agent-lab
"""


def parse_args(argv: list[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="agent-lab: interactive AI agent (智能体)",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )
    parser.add_argument(
        "--backend",
        choices=["mock", "openai"],
        default="mock",
        help="LLM backend to use (default: mock)",
    )
    parser.add_argument(
        "--model",
        default="gpt-4o-mini",
        help="Model name, used only with --backend openai (default: gpt-4o-mini)",
    )
    parser.add_argument(
        "--max-steps",
        type=int,
        default=10,
        help="Maximum ReAct iterations per query (default: 10)",
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Print each reasoning step",
    )
    parser.add_argument(
        "--query",
        metavar="QUESTION",
        help="Run a single query and exit (non-interactive mode)",
    )
    return parser.parse_args(argv)


def build_agent(args: argparse.Namespace) -> Agent:
    llm_kwargs: dict[str, object] = {}
    if args.backend == "openai":
        llm_kwargs["model"] = args.model
    llm = create_llm(args.backend, **llm_kwargs)
    return Agent(llm=llm, max_steps=args.max_steps, verbose=args.verbose)


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)
    agent = build_agent(args)

    if args.query:
        answer = agent.run(args.query)
        print(answer)
        return 0

    # Interactive REPL
    print(BANNER)
    print("Type your question and press Enter. Type 'quit' or 'exit' to leave.\n")

    while True:
        try:
            question = input("You: ").strip()
        except (EOFError, KeyboardInterrupt):
            print("\nGoodbye!")
            break

        if not question:
            continue
        if question.lower() in ("quit", "exit", "q", "退出"):
            print("Goodbye! 再见！")
            break

        answer = agent.run(question)
        print(f"\nAgent: {answer}\n")

    return 0


if __name__ == "__main__":
    sys.exit(main())
