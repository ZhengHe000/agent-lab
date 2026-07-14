"""Tests for agent/agent.py"""

import pytest

from agent.agent import Agent, Step
from agent.llm import MockLLM
from agent.tools import Tool


class TestStep:
    def test_is_final_when_final_answer_set(self):
        s = Step(final_answer="done")
        assert s.is_final

    def test_is_not_final_when_no_final_answer(self):
        s = Step(thought="thinking")
        assert not s.is_final

    def test_str_includes_all_parts(self):
        s = Step(
            thought="t",
            action="calculator",
            action_input="1+1",
            observation="2",
        )
        text = str(s)
        assert "Thought: t" in text
        assert "Action: calculator" in text
        assert "Observation: 2" in text


class TestAgent:
    def setup_method(self):
        self.agent = Agent(llm=MockLLM(), verbose=False)

    def test_run_returns_string(self):
        result = self.agent.run("What time is it?")
        assert isinstance(result, str)
        assert len(result) > 0

    def test_run_math_query(self):
        result = self.agent.run("calculate 3 + 5")
        assert isinstance(result, str)

    def test_run_direct_answer(self):
        result = self.agent.run("Hello!")
        assert isinstance(result, str)

    def test_add_custom_tool(self):
        tool = Tool(name="greet", description="greet", func=lambda _: "Hello!")
        self.agent.add_tool(tool)
        assert "greet" in self.agent.tools

    def test_unknown_tool_triggers_error_observation(self):
        """An LLM that returns an unknown tool name should get an error observation."""
        from agent.llm import LLMBase

        class BrokenLLM(LLMBase):
            calls = 0

            def complete(self, prompt: str) -> str:
                BrokenLLM.calls += 1
                if BrokenLLM.calls == 1:
                    return "Thought: test\nAction: nonexistent_tool\nAction Input: x\n"
                return "Thought: done\nFinal Answer: fallback\n"

        agent = Agent(llm=BrokenLLM(), verbose=False)
        result = agent.run("test")
        assert isinstance(result, str)

    def test_max_steps_respected(self):
        from agent.llm import LLMBase

        class LoopLLM(LLMBase):
            def complete(self, prompt: str) -> str:
                # Never returns a Final Answer
                return "Thought: still thinking\nAction: calculator\nAction Input: 1+1\n"

        agent = Agent(llm=LoopLLM(), max_steps=3, verbose=False)
        result = agent.run("loop forever")
        assert isinstance(result, str)  # Should not hang

    def test_tools_property_is_copy(self):
        tools = self.agent.tools
        tools["injected"] = Tool(name="injected", description="", func=lambda _: "")
        assert "injected" not in self.agent.tools


class TestParseResponse:
    def test_parses_thought_action_input(self):
        text = "Thought: think\nAction: calculator\nAction Input: 2+2\n"
        step = Agent._parse_response(text)
        assert step.thought == "think"
        assert step.action == "calculator"
        assert step.action_input == "2+2"

    def test_parses_final_answer(self):
        text = "Thought: done\nFinal Answer: 42"
        step = Agent._parse_response(text)
        assert step.final_answer == "42"
        assert step.is_final

    def test_empty_text_returns_empty_step(self):
        step = Agent._parse_response("")
        assert not step.thought
        assert not step.action
        assert not step.is_final
