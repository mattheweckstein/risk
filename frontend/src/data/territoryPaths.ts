export interface TerritoryMapData {
  id: string;
  path: string;
  labelX: number;
  labelY: number;
}

export const territoryPaths: TerritoryMapData[] = [
  // === NORTH AMERICA ===
  {
    id: 'alaska',
    path: 'M 45 65 L 80 55 L 110 60 L 115 80 L 100 100 L 70 105 L 40 95 L 30 80 Z',
    labelX: 72,
    labelY: 80,
  },
  {
    id: 'northwest_territory',
    path: 'M 115 60 L 170 45 L 230 42 L 260 55 L 255 80 L 220 95 L 170 100 L 115 80 Z',
    labelX: 185,
    labelY: 70,
  },
  {
    id: 'greenland',
    path: 'M 270 20 L 320 10 L 370 15 L 390 35 L 380 60 L 350 72 L 310 70 L 275 55 L 265 35 Z',
    labelX: 330,
    labelY: 42,
  },
  {
    id: 'alberta',
    path: 'M 100 100 L 170 100 L 175 135 L 170 165 L 120 165 L 100 140 Z',
    labelX: 138,
    labelY: 132,
  },
  {
    id: 'ontario',
    path: 'M 175 95 L 255 80 L 270 55 L 290 70 L 290 100 L 275 130 L 245 145 L 200 145 L 175 135 Z',
    labelX: 230,
    labelY: 115,
  },
  {
    id: 'quebec',
    path: 'M 290 70 L 310 70 L 340 80 L 345 110 L 325 135 L 290 140 L 275 130 L 290 100 Z',
    labelX: 312,
    labelY: 105,
  },
  {
    id: 'western_us',
    path: 'M 100 165 L 170 165 L 190 180 L 185 215 L 160 240 L 115 240 L 90 215 L 85 185 Z',
    labelX: 138,
    labelY: 200,
  },
  {
    id: 'eastern_us',
    path: 'M 190 145 L 245 145 L 275 160 L 280 195 L 265 225 L 230 245 L 190 240 L 170 220 L 185 215 L 190 180 Z',
    labelX: 228,
    labelY: 195,
  },
  {
    id: 'central_america',
    path: 'M 115 240 L 160 240 L 190 240 L 195 265 L 180 295 L 155 310 L 130 305 L 115 280 L 105 260 Z',
    labelX: 152,
    labelY: 275,
  },

  // === SOUTH AMERICA ===
  {
    id: 'venezuela',
    path: 'M 155 320 L 195 315 L 235 318 L 260 335 L 245 355 L 205 360 L 170 355 L 150 340 Z',
    labelX: 205,
    labelY: 338,
  },
  {
    id: 'peru',
    path: 'M 150 355 L 205 360 L 210 395 L 200 435 L 175 455 L 145 440 L 135 400 L 140 370 Z',
    labelX: 172,
    labelY: 405,
  },
  {
    id: 'brazil',
    path: 'M 210 355 L 245 355 L 275 365 L 305 380 L 310 420 L 290 455 L 255 470 L 220 460 L 200 435 L 210 395 Z',
    labelX: 260,
    labelY: 415,
  },
  {
    id: 'argentina',
    path: 'M 175 455 L 220 460 L 255 470 L 245 510 L 225 545 L 200 570 L 180 560 L 170 520 L 160 485 Z',
    labelX: 210,
    labelY: 515,
  },

  // === EUROPE ===
  {
    id: 'iceland',
    path: 'M 410 60 L 445 55 L 465 62 L 460 78 L 435 85 L 410 78 Z',
    labelX: 437,
    labelY: 70,
  },
  {
    id: 'scandinavia',
    path: 'M 480 42 L 520 35 L 545 45 L 555 75 L 545 105 L 520 115 L 495 105 L 480 80 Z',
    labelX: 515,
    labelY: 75,
  },
  {
    id: 'great_britain',
    path: 'M 420 105 L 450 98 L 465 108 L 462 135 L 445 150 L 425 145 L 415 125 Z',
    labelX: 440,
    labelY: 125,
  },
  {
    id: 'northern_europe',
    path: 'M 470 110 L 520 115 L 540 130 L 535 158 L 510 170 L 480 165 L 465 145 Z',
    labelX: 500,
    labelY: 140,
  },
  {
    id: 'western_europe',
    path: 'M 415 160 L 455 155 L 475 170 L 480 200 L 470 225 L 440 235 L 420 220 L 410 190 Z',
    labelX: 445,
    labelY: 195,
  },
  {
    id: 'southern_europe',
    path: 'M 480 170 L 520 170 L 545 185 L 550 215 L 535 240 L 505 248 L 480 235 L 475 205 Z',
    labelX: 512,
    labelY: 210,
  },
  {
    id: 'ukraine',
    path: 'M 545 55 L 590 45 L 640 50 L 660 75 L 665 110 L 655 145 L 635 170 L 600 180 L 565 175 L 545 155 L 540 130 L 545 105 Z',
    labelX: 600,
    labelY: 115,
  },

  // === AFRICA ===
  {
    id: 'north_africa',
    path: 'M 410 260 L 470 255 L 520 258 L 540 275 L 535 310 L 510 340 L 470 350 L 430 345 L 405 320 L 400 290 Z',
    labelX: 470,
    labelY: 300,
  },
  {
    id: 'egypt',
    path: 'M 540 255 L 585 250 L 610 265 L 615 295 L 600 315 L 570 320 L 540 310 L 535 280 Z',
    labelX: 575,
    labelY: 285,
  },
  {
    id: 'east_africa',
    path: 'M 570 325 L 610 315 L 635 335 L 640 370 L 625 400 L 595 415 L 565 405 L 545 375 L 545 345 Z',
    labelX: 590,
    labelY: 370,
  },
  {
    id: 'congo',
    path: 'M 470 355 L 540 350 L 545 375 L 555 405 L 540 435 L 505 445 L 475 430 L 460 400 L 460 370 Z',
    labelX: 505,
    labelY: 395,
  },
  {
    id: 'south_africa',
    path: 'M 475 445 L 540 440 L 570 455 L 580 490 L 565 525 L 535 540 L 500 535 L 475 510 L 465 480 Z',
    labelX: 525,
    labelY: 490,
  },
  {
    id: 'madagascar',
    path: 'M 625 435 L 645 430 L 660 445 L 660 480 L 648 500 L 630 495 L 620 470 L 618 450 Z',
    labelX: 640,
    labelY: 465,
  },

  // === ASIA ===
  {
    id: 'ural',
    path: 'M 665 50 L 710 42 L 740 50 L 745 85 L 740 120 L 720 140 L 690 135 L 665 115 L 660 80 Z',
    labelX: 705,
    labelY: 90,
  },
  {
    id: 'siberia',
    path: 'M 745 30 L 790 22 L 830 28 L 840 60 L 835 95 L 820 120 L 790 125 L 760 115 L 745 85 L 740 55 Z',
    labelX: 790,
    labelY: 75,
  },
  {
    id: 'yakutsk',
    path: 'M 840 25 L 890 18 L 930 25 L 940 55 L 930 85 L 900 95 L 870 88 L 845 65 Z',
    labelX: 890,
    labelY: 55,
  },
  {
    id: 'kamchatka',
    path: 'M 940 25 L 980 18 L 1020 25 L 1040 50 L 1035 85 L 1015 105 L 980 108 L 950 95 L 940 60 Z',
    labelX: 990,
    labelY: 65,
  },
  {
    id: 'irkutsk',
    path: 'M 835 95 L 875 90 L 920 95 L 935 115 L 925 140 L 895 150 L 860 145 L 840 125 Z',
    labelX: 885,
    labelY: 120,
  },
  {
    id: 'mongolia',
    path: 'M 860 150 L 920 145 L 960 155 L 975 180 L 960 205 L 920 210 L 885 200 L 865 178 Z',
    labelX: 918,
    labelY: 178,
  },
  {
    id: 'japan',
    path: 'M 1005 135 L 1030 130 L 1050 142 L 1055 170 L 1045 200 L 1025 210 L 1008 195 L 1000 165 Z',
    labelX: 1028,
    labelY: 170,
  },
  {
    id: 'afghanistan',
    path: 'M 665 145 L 720 140 L 755 152 L 765 180 L 755 210 L 720 220 L 690 215 L 670 195 L 660 170 Z',
    labelX: 715,
    labelY: 180,
  },
  {
    id: 'china',
    path: 'M 770 150 L 830 140 L 880 150 L 920 165 L 935 195 L 920 225 L 880 240 L 830 240 L 790 230 L 765 210 L 770 180 Z',
    labelX: 850,
    labelY: 195,
  },
  {
    id: 'india',
    path: 'M 720 225 L 770 220 L 795 235 L 805 270 L 790 305 L 760 320 L 730 310 L 710 280 L 710 250 Z',
    labelX: 755,
    labelY: 272,
  },
  {
    id: 'siam',
    path: 'M 840 245 L 880 240 L 905 255 L 910 290 L 895 320 L 865 330 L 840 315 L 830 285 L 835 260 Z',
    labelX: 870,
    labelY: 285,
  },
  {
    id: 'middle_east',
    path: 'M 580 195 L 630 185 L 665 195 L 680 225 L 675 260 L 650 280 L 615 285 L 585 270 L 570 240 L 575 215 Z',
    labelX: 625,
    labelY: 240,
  },

  // === AUSTRALIA ===
  {
    id: 'indonesia',
    path: 'M 870 370 L 920 365 L 955 375 L 960 400 L 945 420 L 910 425 L 880 415 L 865 395 Z',
    labelX: 915,
    labelY: 395,
  },
  {
    id: 'new_guinea',
    path: 'M 975 365 L 1020 358 L 1055 368 L 1060 392 L 1045 410 L 1010 415 L 980 405 L 972 385 Z',
    labelX: 1018,
    labelY: 388,
  },
  {
    id: 'western_australia',
    path: 'M 895 445 L 950 438 L 975 455 L 980 495 L 965 530 L 930 545 L 895 535 L 880 500 L 882 465 Z',
    labelX: 930,
    labelY: 490,
  },
  {
    id: 'eastern_australia',
    path: 'M 985 445 L 1030 438 L 1060 458 L 1070 495 L 1060 535 L 1030 550 L 1000 540 L 985 505 L 980 470 Z',
    labelX: 1025,
    labelY: 495,
  },
];

// Connection lines between territories for visual context (optional decorative lines)
export const connectionLines: { from: string; to: string; path: string }[] = [
  // Alaska - Kamchatka (wraps around)
  { from: 'alaska', to: 'kamchatka', path: 'M 45 70 Q 20 30 1040 50' },
  // Greenland - Iceland
  { from: 'greenland', to: 'iceland', path: 'M 390 50 L 410 62' },
  // Brazil - North Africa
  { from: 'brazil', to: 'north_africa', path: 'M 310 390 L 405 280' },
  // Central America - Venezuela
  { from: 'central_america', to: 'venezuela', path: 'M 165 310 L 165 320' },
  // Siam - Indonesia
  { from: 'siam', to: 'indonesia', path: 'M 870 330 L 880 370' },
];

// Continent colors for the background regions
export const continentColors: Record<string, string> = {
  north_america: 'rgba(255, 200, 50, 0.08)',
  south_america: 'rgba(233, 69, 96, 0.08)',
  europe: 'rgba(100, 149, 237, 0.08)',
  africa: 'rgba(255, 165, 0, 0.08)',
  asia: 'rgba(50, 205, 50, 0.08)',
  australia: 'rgba(148, 103, 189, 0.08)',
};
