package adapters


// MultiCommitter manages multiple commit points for atomic operations.
type MultiCommitter struct {
	committers []func() error
	rollbackers []func() error
}

// AddCommitter adds a commit function and its corresponding rollback function.
func (mc *MultiCommitter) AddCommitter(commit func() error, rollback func() error) {
	mc.committers = append(mc.committers, commit)
	mc.rollbackers = append(mc.rollbackers, rollback)
}

// Commit commits all registered transactions atomically.
func (mc *MultiCommitter) Commit() error {
	for _, commit := range mc.committers {
		if err := commit(); err != nil {
			mc.Rollback() // Roll back if any commit fails
			return err
		}
	}
	return nil
}

// Rollback rolls back all registered transactions.
func (mc *MultiCommitter) Rollback() {
	for _, rollback := range mc.rollbackers {
		_ = rollback() // Ignore rollback errors to ensure all are attempted
	}
}
