package tui

// skill describes the competency a module teaches and its "gate" style.
type skill struct {
	modID string
	name  string // human skill label shown in the profile
}

// moduleSkills maps each module to the skill a learner earns by completing it.
// Derived from each module's theme (see subjects/<mod>/module.md). The skill
// is considered *earned* once every exercise in the module passes.
var moduleSkills = []skill{
	{"M0-tools", "Tooling & Shell Scripting"},
	{"M1-clib", "C Foundations"},
	{"M2-procs", "Processes & Signals"},
	{"M3-memory", "Memory & Allocation"},
	{"M4-concurrency", "Concurrency"},
	{"M5-filesio", "Files & I/O"},
	{"M6-networking", "Networking & Sockets"},
	{"M7-resilience", "Resilience & Supervision"},
	{"M8-messaging", "Messaging & Fault Injection"},
	{"M9-time", "Time & Ordering"},
	{"M10-commit", "Commit & Consensus"},
	{"M11-log", "Replicated Logs"},
	{"M12-raft", "Raft Consensus"},
	{"M13-shard", "Sharding"},
	{"M14-membership", "Membership & Failure Detection"},
	{"M15-tx", "Transactions"},
	{"M16-storage", "Storage Engines"},
	{"M17-observability", "Observability"},
}

// skillFor returns the skill label for a module id, or the module id if unknown.
func skillFor(modID string) string {
	for _, s := range moduleSkills {
		if s.modID == modID {
			return s.name
		}
	}
	return modID
}
