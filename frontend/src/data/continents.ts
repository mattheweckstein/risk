export const continentTerritories: Record<string, string[]> = {
  north_america: [
    'alaska', 'northwest_territory', 'greenland', 'alberta',
    'ontario', 'quebec', 'western_us', 'eastern_us', 'central_america',
  ],
  south_america: [
    'venezuela', 'peru', 'brazil', 'argentina',
  ],
  europe: [
    'iceland', 'scandinavia', 'great_britain', 'northern_europe',
    'western_europe', 'southern_europe', 'ukraine',
  ],
  africa: [
    'north_africa', 'egypt', 'east_africa', 'congo',
    'south_africa', 'madagascar',
  ],
  asia: [
    'ural', 'siberia', 'yakutsk', 'kamchatka', 'irkutsk',
    'mongolia', 'japan', 'afghanistan', 'china', 'india',
    'siam', 'middle_east',
  ],
  australia: [
    'indonesia', 'new_guinea', 'western_australia', 'eastern_australia',
  ],
};

export const continentBonuses: Record<string, number> = {
  north_america: 5,
  south_america: 2,
  europe: 5,
  africa: 3,
  asia: 7,
  australia: 2,
};

export const continentDisplayNames: Record<string, string> = {
  north_america: 'North America',
  south_america: 'South America',
  europe: 'Europe',
  africa: 'Africa',
  asia: 'Asia',
  australia: 'Australia',
};
