ALTER TABLE codebases ADD COLUMN stack_id TEXT;

-- Should we set the default to current codebases entries, hardcoded value is java-default
UPDATE codebases set stack_id ='java-default';
