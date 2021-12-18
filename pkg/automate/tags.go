package automate

func hasExpirationTag(tags map[string]*string) bool {
	return hasTag(tags, getExpirationTagName())
}

func hasKeeperTag(tags map[string]*string) bool {
	return hasTag(tags, getKeeperTagName())
}

func hasTag(tags map[string]*string, name string) bool {
	_, ok := tags[name]
	return ok
}
