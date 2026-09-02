# Systems-Forge: Engineering Excellence Study Guide

This curriculum is designed to produce high-caliber infrastructure and distributed systems engineers. It is a grueling, Piscine-style progression. If you are here to learn, you are here to build. Expect to fail, debug, and rewrite until your understanding is total.

## 1. The Methodological Doctrine

1.  **Active Recall & Anki (Spaced Repetition):**
    - You must maintain an **Anki deck per module**.
    - When a module is "authored" in `subjects/`, a corresponding Anki deck must be created.
    - Anki cards are not for "definitions"; they are for "Mechanics." 
    - *Example:* Instead of "What is a Mutex?", ask "How does `pthread_mutex_lock` interact with the CPU cache hierarchy and kernel scheduling?"

2.  **The "Cracked Engineer" Rhythm (Interleaving):**
    - **No Silos.** You will interleave different domains (C, Go, Python, Bash). 
    - **Context-Switching:** Your brain must become adept at switching from low-level memory layout (C/M3) to distributed consensus (Go/M12).

3.  **Feynman Technique (Conceptual Debugging):**
    - You are strictly prohibited from moving to the next exercise if you cannot explain the "failure variant." 
    - If you are stuck for 45 minutes, explain the logic to the "Rubber Duck" (me). If you cannot explain the state machine transitions or the protocol flow, the code is a liability, not an asset.

4.  **Error-Based Learning (The "Broken Variant" Protocol):**
    - You must *intentionally break* your own code. 
    - Once an exercise passes, identify the most dangerous line, delete it, and confirm the system fails exactly as expected. If the system doesn't fail, your test suite is inadequate.

5.  **Engineering Rigor:**
    - Every exercise is a real system, not a toy. If you ship an implementation with leaks, races, or deadlocks, it is a **Failure**. We adhere to `-Wall -Wextra -Werror` and use TSan/ASan as the final arbiters of truth.

---

## 2. The Portfolio Ladder

Your goal is not to "pass," but to "ship." Each module's `ex05` is a **Mini-Capstone Micro-App**.

| Module | Mini-Capstone Focus | Engineering Artifact |
| :--- | :--- | :--- |
| M2 | Process Isolation | Process-based Shell (supports `|`, `&`) |
| M3 | Memory Management | Mmap-backed custom allocator |
| M4 | Parallelism | Thread-safe bounded priority queue |
| M5 | Files/IO | Log-Parser + Streaming Log-Streamer |
| M6 | Networking | Stateful Protocol Gateway (state across TCP connections) |
| M15 | Transactions | Transactional KV Store (2PC, WAL, Chaos-Drilled) |

---

## 3. Operational Discipline

- **The "Artifact" Standard:** Documentation must be stable. If a link is dead, remove it and synthesize the documentation from source (man pages, official RFCs).
- **The "No-Panic" Rule:** If the grader fails, do not panic. The grader is the truth. If it fails, your mental model of the system is misaligned with the protocol implementation. **Fix the mental model first, then the code.**
