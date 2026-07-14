# agent-lab

A lightweight **AI Agent (智能体)** framework implemented in Python, demonstrating the **ReAct** (Reasoning + Acting) pattern.

## Architecture

```
agent-lab/
├── agent/
│   ├── __init__.py   # Public API: Agent, Tool
│   ├── agent.py      # Core ReAct loop (Step, Agent)
│   ├── llm.py        # LLM abstraction (MockLLM, OpenAILLM, create_llm)
│   └── tools.py      # Built-in tools (calculator, clock, file I/O, web search)
├── tests/
│   ├── test_agent.py
│   ├── test_llm.py
│   ├── test_main.py
│   └── test_tools.py
├── main.py           # Interactive CLI
├── pytest.ini
└── requirements.txt
```

## How it works

The agent follows the **ReAct loop**:

```
Question → Thought → Action → Observation → Thought → … → Final Answer
```

1. **Thought** — the agent reasons about the question.
2. **Action** — the agent calls a registered tool.
3. **Observation** — the tool's result is fed back to the agent.
4. Steps repeat until the agent produces a **Final Answer** or the step limit is reached.

## Quick start

```bash
# No dependencies required for the mock backend
python main.py
```

### Single-query mode

```bash
python main.py --query "What time is it?"
python main.py --query "calculate 2 ** 32"
python main.py --query "search for Python agent tutorials"
```

### Verbose output (shows each reasoning step)

```bash
python main.py --query "What is sqrt(144)?" --verbose
```

### OpenAI backend

```bash
pip install openai
OPENAI_API_KEY=sk-... python main.py --backend openai --model gpt-4o-mini
```

## Built-in tools

| Tool | Description |
|------|-------------|
| `calculator` | Evaluate math expressions (`sqrt(2)`, `2**32`, etc.) |
| `current_time` | Return the current local date and time |
| `read_file` | Read a local file by path |
| `write_file` | Write content to a local file (`<path>\|<content>`) |
| `web_search` | Search stub — replace with a real API for live results |

## Adding custom tools

```python
from agent import Agent, Tool

def reverse(text: str) -> str:
    return text[::-1]

agent = Agent()
agent.add_tool(Tool(name="reverse", description="Reverse a string.", func=reverse))
print(agent.run("reverse the word hello"))
```

## Running tests

```bash
python -m pytest
```
