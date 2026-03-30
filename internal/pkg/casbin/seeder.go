package casbin

import "manga-go/internal/pkg/logger"

const GlobalDomain = "global"

// SeedGlobalPolicies seeds the initial global Casbin policies if they don't exist.
// This should be called once on application startup.
func SeedGlobalPolicies(enforcer *Enforcer, log *logger.Logger) {
	policies := [][]string{
		// anonymous can read comics and chapters
		{"anonymous", GlobalDomain, "comic", "read"},
		{"anonymous", GlobalDomain, "chapter", "read"},
		// authenticated non-group users can create and join groups
		{"user", GlobalDomain, "group", "create"},
		{"user", GlobalDomain, "group", "join"},
	}

	for _, p := range policies {
		if ok, err := enforcer.AddPolicy(p[0], p[1], p[2], p[3]); err != nil {
			log.Errorf("Failed to add Casbin policy %v: %v", p, err)
		} else if ok {
			log.Infof("Casbin policy added: %v", p)
		}
	}

	// Role inheritance: user inherits anonymous in global domain
	if ok, err := enforcer.AddGroupingPolicy("user", "anonymous", GlobalDomain); err != nil {
		log.Errorf("Failed to add Casbin grouping policy user->anonymous: %v", err)
	} else if ok {
		log.Info("Casbin grouping policy added: user -> anonymous in global")
	}
}

// SeedGroupPolicies adds Casbin policies and role inheritance for a newly created group.
// groupID should be the string representation of the group's UUID.
func SeedGroupPolicies(enforcer *Enforcer, groupID string, log *logger.Logger) {
	// chapter_creator can create and update chapters in the group domain
	if ok, err := enforcer.AddPolicy("chapter_creator", groupID, "chapter", "create"); err != nil {
		log.Errorf("Failed to add chapter_creator create policy for group %s: %v", groupID, err)
	} else if ok {
		log.Infof("Casbin policy added: chapter_creator can create chapter in group %s", groupID)
	}

	if ok, err := enforcer.AddPolicy("chapter_creator", groupID, "chapter", "update"); err != nil {
		log.Errorf("Failed to add chapter_creator update policy for group %s: %v", groupID, err)
	} else if ok {
		log.Infof("Casbin policy added: chapter_creator can update chapter in group %s", groupID)
	}

	// group_owner can kick members
	if ok, err := enforcer.AddPolicy("group_owner", groupID, "group_member", "kick"); err != nil {
		log.Errorf("Failed to add group_owner kick policy for group %s: %v", groupID, err)
	} else if ok {
		log.Infof("Casbin policy added: group_owner can kick group_member in group %s", groupID)
	}

	// group_owner can grant permissions
	if ok, err := enforcer.AddPolicy("group_owner", groupID, "group_member", "grant"); err != nil {
		log.Errorf("Failed to add group_owner grant policy for group %s: %v", groupID, err)
	} else if ok {
		log.Infof("Casbin policy added: group_owner can grant group_member in group %s", groupID)
	}

	// Role inheritance within the group domain:
	// group_owner -> chapter_creator -> group_member
	if ok, err := enforcer.AddGroupingPolicy("chapter_creator", "group_member", groupID); err != nil {
		log.Errorf("Failed to add chapter_creator->group_member grouping for group %s: %v", groupID, err)
	} else if ok {
		log.Infof("Casbin grouping added: chapter_creator -> group_member in group %s", groupID)
	}

	if ok, err := enforcer.AddGroupingPolicy("group_owner", "chapter_creator", groupID); err != nil {
		log.Errorf("Failed to add group_owner->chapter_creator grouping for group %s: %v", groupID, err)
	} else if ok {
		log.Infof("Casbin grouping added: group_owner -> chapter_creator in group %s", groupID)
	}
}
