package models

// TerritoryDefinition contains the static definition of a territory
type TerritoryDefinition struct {
	ID        string
	Name      string
	Continent string
	Neighbors []string
}

// AllTerritories defines all 42 Risk territories with their adjacencies
var AllTerritories = []TerritoryDefinition{
	// North America
	{ID: "alaska", Name: "Alaska", Continent: "north_america", Neighbors: []string{"northwest_territory", "alberta", "kamchatka"}},
	{ID: "northwest_territory", Name: "Northwest Territory", Continent: "north_america", Neighbors: []string{"alaska", "alberta", "ontario", "greenland"}},
	{ID: "greenland", Name: "Greenland", Continent: "north_america", Neighbors: []string{"northwest_territory", "ontario", "quebec", "iceland"}},
	{ID: "alberta", Name: "Alberta", Continent: "north_america", Neighbors: []string{"alaska", "northwest_territory", "ontario", "western_us"}},
	{ID: "ontario", Name: "Ontario", Continent: "north_america", Neighbors: []string{"northwest_territory", "alberta", "greenland", "quebec", "western_us", "eastern_us"}},
	{ID: "quebec", Name: "Quebec", Continent: "north_america", Neighbors: []string{"ontario", "greenland", "eastern_us"}},
	{ID: "western_us", Name: "Western United States", Continent: "north_america", Neighbors: []string{"alberta", "ontario", "eastern_us", "central_america"}},
	{ID: "eastern_us", Name: "Eastern United States", Continent: "north_america", Neighbors: []string{"ontario", "quebec", "western_us", "central_america"}},
	{ID: "central_america", Name: "Central America", Continent: "north_america", Neighbors: []string{"western_us", "eastern_us", "venezuela"}},

	// South America
	{ID: "venezuela", Name: "Venezuela", Continent: "south_america", Neighbors: []string{"central_america", "peru", "brazil"}},
	{ID: "peru", Name: "Peru", Continent: "south_america", Neighbors: []string{"venezuela", "brazil", "argentina"}},
	{ID: "brazil", Name: "Brazil", Continent: "south_america", Neighbors: []string{"venezuela", "peru", "argentina", "north_africa"}},
	{ID: "argentina", Name: "Argentina", Continent: "south_america", Neighbors: []string{"peru", "brazil"}},

	// Europe
	{ID: "iceland", Name: "Iceland", Continent: "europe", Neighbors: []string{"greenland", "scandinavia", "great_britain"}},
	{ID: "scandinavia", Name: "Scandinavia", Continent: "europe", Neighbors: []string{"iceland", "great_britain", "northern_europe", "ukraine"}},
	{ID: "great_britain", Name: "Great Britain", Continent: "europe", Neighbors: []string{"iceland", "scandinavia", "northern_europe", "western_europe"}},
	{ID: "northern_europe", Name: "Northern Europe", Continent: "europe", Neighbors: []string{"scandinavia", "great_britain", "western_europe", "southern_europe", "ukraine"}},
	{ID: "western_europe", Name: "Western Europe", Continent: "europe", Neighbors: []string{"great_britain", "northern_europe", "southern_europe", "north_africa"}},
	{ID: "southern_europe", Name: "Southern Europe", Continent: "europe", Neighbors: []string{"northern_europe", "western_europe", "ukraine", "north_africa", "egypt", "middle_east"}},
	{ID: "ukraine", Name: "Ukraine", Continent: "europe", Neighbors: []string{"scandinavia", "northern_europe", "southern_europe", "ural", "afghanistan", "middle_east"}},

	// Africa
	{ID: "north_africa", Name: "North Africa", Continent: "africa", Neighbors: []string{"brazil", "western_europe", "southern_europe", "egypt", "east_africa", "congo"}},
	{ID: "egypt", Name: "Egypt", Continent: "africa", Neighbors: []string{"southern_europe", "north_africa", "east_africa", "middle_east"}},
	{ID: "east_africa", Name: "East Africa", Continent: "africa", Neighbors: []string{"north_africa", "egypt", "congo", "south_africa", "madagascar", "middle_east"}},
	{ID: "congo", Name: "Congo", Continent: "africa", Neighbors: []string{"north_africa", "east_africa", "south_africa"}},
	{ID: "south_africa", Name: "South Africa", Continent: "africa", Neighbors: []string{"congo", "east_africa", "madagascar"}},
	{ID: "madagascar", Name: "Madagascar", Continent: "africa", Neighbors: []string{"east_africa", "south_africa"}},

	// Asia
	{ID: "ural", Name: "Ural", Continent: "asia", Neighbors: []string{"ukraine", "siberia", "afghanistan", "china"}},
	{ID: "siberia", Name: "Siberia", Continent: "asia", Neighbors: []string{"ural", "yakutsk", "irkutsk", "mongolia", "china"}},
	{ID: "yakutsk", Name: "Yakutsk", Continent: "asia", Neighbors: []string{"siberia", "kamchatka", "irkutsk"}},
	{ID: "kamchatka", Name: "Kamchatka", Continent: "asia", Neighbors: []string{"alaska", "yakutsk", "irkutsk", "mongolia", "japan"}},
	{ID: "irkutsk", Name: "Irkutsk", Continent: "asia", Neighbors: []string{"siberia", "yakutsk", "kamchatka", "mongolia"}},
	{ID: "mongolia", Name: "Mongolia", Continent: "asia", Neighbors: []string{"siberia", "irkutsk", "kamchatka", "japan", "china"}},
	{ID: "japan", Name: "Japan", Continent: "asia", Neighbors: []string{"kamchatka", "mongolia"}},
	{ID: "afghanistan", Name: "Afghanistan", Continent: "asia", Neighbors: []string{"ukraine", "ural", "china", "india", "middle_east"}},
	{ID: "china", Name: "China", Continent: "asia", Neighbors: []string{"ural", "siberia", "mongolia", "afghanistan", "india", "siam"}},
	{ID: "india", Name: "India", Continent: "asia", Neighbors: []string{"afghanistan", "china", "siam", "middle_east"}},
	{ID: "siam", Name: "Siam", Continent: "asia", Neighbors: []string{"china", "india", "indonesia"}},
	{ID: "middle_east", Name: "Middle East", Continent: "asia", Neighbors: []string{"ukraine", "southern_europe", "egypt", "east_africa", "afghanistan", "india"}},

	// Australia
	{ID: "indonesia", Name: "Indonesia", Continent: "australia", Neighbors: []string{"siam", "new_guinea", "western_australia"}},
	{ID: "new_guinea", Name: "New Guinea", Continent: "australia", Neighbors: []string{"indonesia", "western_australia", "eastern_australia"}},
	{ID: "western_australia", Name: "Western Australia", Continent: "australia", Neighbors: []string{"indonesia", "new_guinea", "eastern_australia"}},
	{ID: "eastern_australia", Name: "Eastern Australia", Continent: "australia", Neighbors: []string{"new_guinea", "western_australia"}},
}
