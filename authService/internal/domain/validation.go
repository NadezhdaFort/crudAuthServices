package domain

func (lp LoginPassword) IsValid() bool {
	if lp.Login == "" || lp.Password == "" {
		return false
	}
	return true
}

func (up UserPassword) IsValid() bool {
	if up.ID.IsZero() || up.Password == "" {
		return false
	}
	return true
}

func (ui UserInfo) IsValid() bool {
	if ui.ID.IsZero() || ui.Name == "" {
		return false
	}
	return true
}

func (ur UserRole) IsValid() bool {
	if ur.ID.IsZero() || (ur.Role != "admin" && ur.Role != "user") {
		return false
	}
	return true
}

func (ub SetBlockUser) IsValid() bool {
	return !ub.ID.IsZero()
}
